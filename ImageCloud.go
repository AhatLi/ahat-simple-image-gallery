package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const imgPath string = "images/"
const thumPath string = "thumbnail/"
const assetPath string = "assets/"

var htmlFile string = "assets/html/index.html.ahat"

func main() {
	//서버 구성 초기화
	startTime := time.Now()
	initServer()
	elapsedTime := time.Since(startTime)

	fmt.Println("실행시간: ", elapsedTime.Seconds())

	fmt.Println("init complete. server start")

	//로그인 페이지
	http.HandleFunc("/main", indexPageHandler)
	//로그인 동작 API
	http.HandleFunc("/login", loginHandler)
	//로그아웃 동작 API
	http.HandleFunc("/logout", logoutHandler)

	//이미지 파일을 요청하는 URL 주소
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	//서버 구성요소를 요청하는 URL 주소
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))
	//파일 이동, 삭제 등의 API를 요청하는 URL 주소
	http.HandleFunc("/api/", apiHandler)
	//이미지 페이지를 요청하는 URL 주소
	http.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(getPort(), nil)

	if err != nil {
		fmt.Println("ListenAndServe:" + err.Error())
	} else {
		fmt.Println("ListenAndServe Started! -> Port(" + getPort() + ")")
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	printLog(r)
	if loginCheck(w, r) == false {
		return
	}

	fmt.Println(htmlFile)
	data, _ := ioutil.ReadFile(htmlFile)
	str := fmt.Sprintf("%s", data)

	//현재 페이지 번호를 파라미터로 받는다.
	page, _ := strconv.Atoi(r.URL.Query().Get("p"))
	//검색 텍스트가있을경우 받아온다.
	search := r.URL.Query().Get("searchText")
	//설정파일에서 현재 표시할 이미지의 수와 정렬방식을 받아온다.
	count, contentSort := getContentData()

	if page == 0 {
		page = 1
	}

	path := imgPath + r.URL.Path

	if r.URL.Path != "/" && !fileExists(path) {
		return
	}

	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		fmt.Println(errf)
	}

	if search != "" {
		files = fileSearch(files, search)
	}

	files = imgFilter(files, search)

	switch contentSort {
	case "name":
		sort.Sort(FileNameSort(files))
	case "date":
		sort.Sort(FileDateSort(files))
	case "size":
		sort.Sort(FileSizeSort(files))
	}

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "#") {
			switch scanner.Text() {
			case "#select":
				fmt.Fprintf(w, "<script>const contentCount="+strconv.Itoa(count)+";const contentSort='"+contentSort+"';</script>")
				makeSelect(w)
			case "#content":
				makeContent(w, r, count, page, contentSort, files)
			case "#page":
				makePage(w, r, count, page, files)
			}
		} else {
			fmt.Fprintf(w, scanner.Text())
		}
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	printLog(r)
	if strings.HasSuffix(r.URL.Path, "move") { //파일 이동 동작을 하는 API를 실행한다.
		fileMove(r.PostFormValue("files"), r.PostFormValue("source"), r.PostFormValue("dest"))
	} else if strings.HasSuffix(r.URL.Path, "delete") { //파일 삭제 동작을 하는 API를 실행한다.
		fileRemove(r.PostFormValue("files"), r.PostFormValue("path"))
	} else if strings.HasSuffix(r.URL.Path, "config") { //설정을 변경했을 경우 해당 설정을 적용하고 파일에 반영한다.
		setContentData(r.PostFormValue("imgCount"), r.PostFormValue("imgSort"))
	}
}

func initServer() {
	//서버의 구성요소 폴더가 존재하지 않을경우 생성한다.
	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	//html 파일 위치를 설정파일에서 불러온다.
	//기본은 영문 html 파일이며 설정으로 한국어 파일을 가져올 수 있다.
	htmlFile = getEnvData()

	//이미지파일을 확인하여 썸네일을 생성하는 동작을 한다.
	preExplorerDirectory(imgPath)
}

// 디렉토리의 파일을 읽어 페이지에 HTML 형식으로 표시한다.
func makeContent(w http.ResponseWriter, r *http.Request, count int, page int, contentSort string, files []os.FileInfo) {
	// 첫번째 DIV로 뒤로 이동하는 동작을 하는 아이콘을 표시해준다.
	fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(assetPath+"images/directory.png")+`")' 
	onClick='location.href=".."'></div><div class='imgtitle'>..</div></div>`)

	//현재 페이지번호와 이미지 개수 변수를 이용하여 페이징 기능 동작을 한다.
	if len(files) > ((page - 1) * count) {
		files = files[(page-1)*count:]
	}
	//마지막 페이지일 경우 다음버튼일 동작하지 않도록 하기 위한 동작.
	if len(files) > count {
		files = files[0:count]
		fmt.Fprintf(w, "<script>var lastPage=false;</script>")
	} else {
		fmt.Fprintf(w, "<script>var lastPage=true;</script>")
	}
	//files의 요소가 폴더일 경우 폴더로 표시하고 이미지파일일 경우 썸네일을 표시한다.
	jsonText := "[ "
	c := 0
	for i, f := range files {
		if f.IsDir() {
			fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(assetPath+"images/directory.png")+`")' 
			onClick='location.href="http://`+r.Host+"/"+r.URL.Path+"/"+f.Name()+`/"'></div><div class='imgtitle'>`+f.Name()+"</div></div>")
		} else {
			fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(thumPath+imgPath+r.URL.Path+"/"+f.Name()+".jpg")+`")'
			' onClick='showGallery(`+strconv.Itoa(c)+`, this.id)' id='img`+strconv.Itoa(i)+`' ontouchstart='longTouch(this.id)' ontouchend='revert(this.id)'>
			</div><div class='imgtitle'>`+f.Name()+"</div></div>")
			jsonText += `{title: "` + f.Name() + `", src: "http://` + r.Host + "/" + imgPath + r.URL.Path + "/" + f.Name() + `"},`
			c++
		}
	}
	jsonText = jsonText[0:len(jsonText)-1] + "];"
	fmt.Fprintf(w, "<script>const gallery = "+jsonText+"</script>")
}

//페이지 이동을 위한 select 엘레멘트를 만드는 함수
func makePage(w http.ResponseWriter, r *http.Request, count int, page int, files []os.FileInfo) {

	pageno := (len(files) / count) + 1

	fmt.Fprintf(w, "<script>var page="+strconv.Itoa(page)+";var count="+strconv.Itoa(count)+";</script>")
	fmt.Fprintf(w, "<select class='form-control' style='background-color: #6c757d; color: white;' id='pageSelect'>")
	for i := 0; i < pageno; i++ {
		if page == (i + 1) {
			fmt.Fprintf(w, "<option selected>"+strconv.Itoa(i+1)+"</option>")
		} else {
			fmt.Fprintf(w, "<option>"+strconv.Itoa(i+1)+"</option>")
		}
	}
	fmt.Fprintf(w, "</select>")
}

//파일 이동을 위해 폴더 구조가 들어있는 select 엘레멘트를 생성한다.
func makeSelect(w http.ResponseWriter) {

	var dir []string
	getDirPath("./images", &dir)

	fmt.Fprintf(w, "<select class='form-select' id='selectDir'>")

	fmt.Fprintf(w, "<option>-Select Directory-</option>")
	for _, d := range dir {
		fmt.Fprintf(w, "<option>"+d+"</option>")
	}
	fmt.Fprintf(w, "</select>")
}
