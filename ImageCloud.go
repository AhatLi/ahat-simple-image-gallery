package main

//필요 패키지 임포트
import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

const imgPath string = "images/"
const thumPath string = "thumbnail/"
const assetPath string = "assets/"

func imgToBase64(file string) string {
	f, _ := os.Open(file)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	return "data:image/png;base64," + encoded
}

func indexTandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)

	if r.URL.String() != "/" && !fileExists(r.URL.String()) {
		return
	}

	fmt.Fprintf(w, "<html><head><style>.modal {  display: none;  position: fixed;   z-index: 1;  padding-top: 100px;  left: 0;  top: 0;  width: 100%%;  height: 100%%;  overflow: auto;  background-color: rgb(0,0,0);  background-color: rgba(0,0,0,0.4);}.modal-content {  background-color: #fefefe;  margin: auto;  padding: 20px;  border: 1px solid #888;  width: 80%%;}</style></head><body><div id='myModal' class='modal'>  <div class='modal-content'>    <span class='close'>&times;</span>    <p><img src='https://blog.jinbo.net/attach/615/200937431.jpg' style='width: 100%%;padding-bottom: 25%%' id='myImg'></p>  </div></div>")
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
	fmt.Fprintf(w, "<title>Go</title>")

	path := thumPath + imgPath + r.URL.String()
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for i, f := range files {
		if i%5 == 0 {
			fmt.Fprintf(w, "<br>")
		}

		if f.IsDir() {
			//	explorerDirectory(path + f.Name())'"+imgPath+r.URL.String()+f.Name()+"'
		} else {
			fmt.Fprintf(w, "<img src='"+imgToBase64(path+f.Name())+"' style='width:20%%;' onClick='thumbClick(\""+imgPath+r.URL.String()+f.Name()+"\")' name='myBtn'>")
		}
	}
	fmt.Fprintf(w, "<script src='/assets/script.js'></script>")
	fmt.Fprintf(w, "</body></html>")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	// load original image
	img, err := imaging.Open(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	thumbnail := imaging.CropCenter(imaging.Resize(img, 80, 0, imaging.Lanczos), 80, 80)
	err = imaging.Save(thumbnail, thumbname)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initServer() {
	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	explorerDirectory(imgPath)
}

func explorerDirectory(path string) {
	os.MkdirAll(thumPath+path, os.ModePerm)
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for _, f := range files {
		if f.IsDir() {
			explorerDirectory(path + f.Name())
		} else {
			makeThumbnail(path + "/" + f.Name())
		}
	}
}

func main() {
	initServer()

	//기본 Url 핸들러 메소드 지정
	http.HandleFunc("/", indexTandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	http.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(thumPath))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))
	//서버 시작
	err := http.ListenAndServe(":9090", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
