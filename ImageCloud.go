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
	if loginCheck(w, r) {
		fmt.Println("RemoteAddr : " + r.RemoteAddr)

		data, _ := ioutil.ReadFile("./assets/index.html.ahat")
		str := fmt.Sprintf("%s", data)

		fmt.Fprintf(w, str[:strings.Index(str, "#select")])
		//	fmt.Fprintf(w, str[strings.LastIndex(str, "#content")+9:strings.LastIndex(str, "#select")])

		makeSelect(w)

		fmt.Fprintf(w, str[strings.LastIndex(str, "#select")+8:strings.LastIndex(str, "#content")])
		//	fmt.Fprintf(w, str[:strings.Index(str, "#content")])

		displayImages(w, r)

		fmt.Fprintf(w, str[strings.LastIndex(str, "#content")+8:])
	}
}

func apiTandler(w http.ResponseWriter, r *http.Request) {

	files := strings.Split(r.PostFormValue("files"), ",")
	for _, file := range files {
		err := os.Rename(r.PostFormValue("source")+file, "images"+r.PostFormValue("dest")+"/"+file)
		if err != nil {
			log.Fatal("Rename error1 : " + err.Error())
		}
		fmt.Println("." + r.PostFormValue("source") + file)
		fmt.Println("./images" + r.PostFormValue("dest") + "/" + file)
		fmt.Println()

		err = os.Rename(thumPath+r.PostFormValue("source")+file+".jpg", thumPath+imgPath+r.PostFormValue("dest")+"/"+file+".jpg")
		if err != nil {
			log.Fatal("Rename error2 : " + err.Error())
		}
		fmt.Println()
	}

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

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		fmt.Println("ListenAndServe:" + err.Error())
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
