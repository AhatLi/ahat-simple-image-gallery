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

func indexTandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RemoteAddr : " + r.RemoteAddr)
	if loginCheck(w, r) == false {
		return
	}

	data, _ := ioutil.ReadFile("assets/index.html.ahat")
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

func apiTandler(w http.ResponseWriter, r *http.Request) {

	if strings.HasSuffix(r.URL.Path, "move") {
		fileMove(r.PostFormValue("files"), r.PostFormValue("source"), r.PostFormValue("dest"))
	} else if strings.HasSuffix(r.URL.Path, "remove") {
		fileRemove(r.PostFormValue("files"), r.PostFormValue("path"))
	} else if strings.HasSuffix(r.URL.Path, "config") {
		setContentData(r.PostFormValue("imgCount"), r.PostFormValue("imgSort"))
	}
}

func fileRemove(files string, path string) {
	filesSplit := strings.Split(files, ",")
	for _, file := range filesSplit {

		err := os.Remove(path + file)
		if err != nil {
			fmt.Println("remove error1 : " + err.Error())
		}
		fmt.Println("remove " + path + file)

		err = os.Remove(thumPath + path + file + ".jpg")
		if err != nil {
			fmt.Println("Rename error2 : " + err.Error())
		}
		fmt.Println()
	}
}

func fileMove(files string, source string, dest string) {
	filesSplit := strings.Split(files, ",")
	for _, file := range filesSplit {
		err := os.Rename(source+file, "images"+dest+"/"+file)
		if err != nil {
			fmt.Println("Rename error1 : " + err.Error())
		}
		fmt.Println("." + source + file)
		fmt.Println("./images" + dest + "/" + file)
		fmt.Println()

		err = os.Rename(thumPath+source+file+".jpg", thumPath+imgPath+dest+"/"+file+".jpg")
		if err != nil {
			fmt.Println("Rename error2 : " + err.Error())
		}
		fmt.Println()
	}
}

func makeContent(w http.ResponseWriter, r *http.Request, count int, page int, contentSort string, files []os.FileInfo) {

	fmt.Fprintf(w, "<td></td><td></td><td></td><td></td></tr><tr>")
	fmt.Fprintf(w, "<td class='equalDivide'><a href='..'><img src='http://"+r.Host+"/assets/directory.png'></a><br>..</td>")
	fmt.Println("<a href='http://" + r.Host + "" + r.URL.Path + "..'>")

	fmt.Println(len(files))
	if len(files) > ((page - 1) * count) {
		files = files[(page-1)*count:]
	}
	if len(files) > count {
		files = files[0:count]
		fmt.Fprintf(w, "<script>var lastPage=false;</script>")
	} else {
		fmt.Fprintf(w, "<script>var lastPage=true;</script>")
	}
	fmt.Println(len(files))

	for i, f := range files {
		if (i+1)%4 == 0 {
			fmt.Fprintf(w, "</tr><tr>")
		}
		fmt.Fprintf(w, "<td class='equalDivide'>")
		if f.IsDir() {
			fmt.Fprintf(w, "<a href=\"http://"+r.Host+"/"+r.URL.Path+"/"+f.Name()+"/\"><img src='http://"+r.Host+"/assets/directory.png'></a>")
			fmt.Fprintf(w, "<br>"+f.Name())
		} else {
			fmt.Fprintf(w, "<img src='"+imgToBase64(thumPath+imgPath+r.URL.Path+"/"+f.Name()+".jpg")+"' id='img"+strconv.Itoa(i)+"' ontouchstart='longTouch(this.id)' ontouchend='revert(this.id)' name='http://"+r.Host+"/"+imgPath+r.URL.Path+"/"+f.Name()+"'>")
			fmt.Fprintf(w, "<br>"+f.Name())
		}
		fmt.Fprintf(w, "</td>")
	}

	fmt.Fprintf(w, "</tr>")
}

func makePage(w http.ResponseWriter, r *http.Request, count int, page int, files []os.FileInfo) {

	pageno := (len(files) / count) + 1

	fmt.Fprintf(w, "<select>")
	for i := 0; i < pageno; i++ {
		fmt.Fprintf(w, "<option>"+strconv.Itoa(i+1)+"</option>")
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

func initServer() {
	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	explorerDirectory(imgPath)
}

func main() {
	initServer()

	fmt.Println("init complete. server start")

	http.HandleFunc("/main", indexPageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))
	http.HandleFunc("/api/", apiTandler)
	http.HandleFunc("/", indexTandler)

	err := http.ListenAndServe(getPort(), nil)

	if err != nil {
		fmt.Println("ListenAndServe:" + err.Error())
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
