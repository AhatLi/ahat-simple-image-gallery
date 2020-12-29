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
)

const imgPath string = "images/"
const thumPath string = "thumbnail/"
const assetPath string = "assets/"

var htmlFile string = "assets/index.html.ahat"

func main() {
	initServer()

	fmt.Println("init complete. server start")

	http.HandleFunc("/main", indexPageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(getPort(), nil)

	if err != nil {
		fmt.Println("ListenAndServe:" + err.Error())
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
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

	page, _ := strconv.Atoi(r.URL.Query().Get("p"))
	search := r.URL.Query().Get("searchText")
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
				fmt.Fprintf(w, "<script>var contentCount="+strconv.Itoa(count)+";var contentSort='"+contentSort+"';</script>")
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
	if strings.HasSuffix(r.URL.Path, "move") {
		fileMove(r.PostFormValue("files"), r.PostFormValue("source"), r.PostFormValue("dest"))
	} else if strings.HasSuffix(r.URL.Path, "remove") {
		fileRemove(r.PostFormValue("files"), r.PostFormValue("path"))
	} else if strings.HasSuffix(r.URL.Path, "config") {
		setContentData(r.PostFormValue("imgCount"), r.PostFormValue("imgSort"))
	}
}

func initServer() {
	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	htmlFile = getEnvData()
	fmt.Println(htmlFile)

	explorerDirectory(imgPath)
}

func makeContent(w http.ResponseWriter, r *http.Request, count int, page int, contentSort string, files []os.FileInfo) {

	fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(assetPath+"directory.png")+`")' 
	onClick='location.href=".."'></div>..</div>`)

	if len(files) > ((page - 1) * count) {
		files = files[(page-1)*count:]
	}
	if len(files) > count {
		files = files[0:count]
		fmt.Fprintf(w, "<script>var lastPage=false;</script>")
	} else {
		fmt.Fprintf(w, "<script>var lastPage=true;</script>")
	}

	for i, f := range files {
		if f.IsDir() {
			fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(assetPath+"directory.png")+`")' 
			onClick='location.href="http://`+r.Host+"/"+r.URL.Path+"/"+f.Name()+`/"'></div>`+f.Name()+"</div>")
		} else {
			fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url("`+imgToBase64(thumPath+imgPath+r.URL.Path+"/"+f.Name()+".jpg")+`")'
			id='img`+strconv.Itoa(i)+`' onClick='thumbClick(this.id)' ontouchstart='longTouch(this.id)' ontouchend='revert(this.id)' 
			title='http://`+r.Host+"/"+imgPath+r.URL.Path+"/"+f.Name()+"'  ></div>"+f.Name()+"</div>")
		}
	}
}

func makePage(w http.ResponseWriter, r *http.Request, count int, page int, files []os.FileInfo) {

	pageno := (len(files) / count) + 1

	fmt.Fprintf(w, "<script>var page="+strconv.Itoa(page)+";var count="+strconv.Itoa(count)+";</script>")
	fmt.Fprintf(w, "<select id='pageSelect'>")
	for i := 0; i < pageno; i++ {
		if page == (i + 1) {
			fmt.Fprintf(w, "<option selected>"+strconv.Itoa(i+1)+"</option>")
		} else {
			fmt.Fprintf(w, "<option>"+strconv.Itoa(i+1)+"</option>")
		}
	}
	fmt.Fprintf(w, "</select>")
}

func makeSelect(w http.ResponseWriter) {

	var dir []string
	getDirPath("./images", &dir)

	fmt.Fprintf(w, "<select id='selectDir'>")

	fmt.Fprintf(w, "<option>-filemove-</option>")
	for _, d := range dir {
		fmt.Fprintf(w, "<option>"+d+"</option>")
	}
	fmt.Fprintf(w, "</select>")
}
