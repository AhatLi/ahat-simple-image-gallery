package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

const imgPath string = "images/"
const thumPath string = "thumbnail/"
const assetPath string = "assets/"

func indexTandler(w http.ResponseWriter, r *http.Request) {

	data, _ := ioutil.ReadFile("./assets/index.html.ahat")
	str := fmt.Sprintf("%s", data)

	fmt.Fprintf(w, str[:strings.Index(str, "#content")])

	displayImages(w, r)

	fmt.Fprintf(w, str[strings.LastIndex(str, "#content")+9:])
}

func displayImages(w http.ResponseWriter, r *http.Request) {

	decodedValue, _ := url.QueryUnescape(r.URL.String())
	path := imgPath + decodedValue
	fmt.Println("path : " + path)

	if decodedValue != "/" && !fileExists(path) {
		return
	}

	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	sort.Sort(FileSort(files))
	fmt.Fprintf(w, "<td></td><td></td><td></td><td></td></tr><tr>")
	fmt.Fprintf(w, "<td class='equalDivide'><a href='..'><img src='http://"+r.Host+"/assets/directory.png'></a><br>..</td>")
	fmt.Println("<a href='http://" + r.Host + "" + decodedValue + "..'>")

	for i, f := range files {
		if (i+1)%4 == 0 {
			fmt.Fprintf(w, "</tr><tr>")
		}
		fmt.Fprintf(w, "<td class='equalDivide'>")
		if f.IsDir() {
			fmt.Fprintf(w, "<a href=\"http://"+r.Host+"/"+decodedValue+"/"+f.Name()+"/\"><img src='http://"+r.Host+"/assets/directory.png'></a>")
			fmt.Fprintf(w, "<br>"+f.Name())
		} else {
			fmt.Fprintf(w, "<img src='"+imgToBase64(thumPath+path+"/"+f.Name()+".jpg")+"' id='img"+strconv.Itoa(i)+"' ontouchstart='func(this.id)' ontouchend='revert(this.id)' onClick='thumbClick(this.id)' name='http://"+r.Host+"/"+path+"/"+f.Name()+"'>")
			fmt.Fprintf(w, "<br>"+f.Name())
		}
		fmt.Fprintf(w, "</td>")
	}

	fmt.Fprintf(w, "</tr>")
}

func initServer() {
	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	explorerDirectory(imgPath)
}

func main() {
	initServer()

	getAllDirPath()

	fmt.Println("init complete. server start")

	http.HandleFunc("/", indexTandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	http.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(thumPath))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		fmt.Println("ListenAndServe:" + err.Error())
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
