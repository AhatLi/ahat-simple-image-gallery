package main

//필요 패키지 임포트
import (
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

func indexTandler(w http.ResponseWriter, r *http.Request) {
	// MAIN SECTION HTML CODE
	fmt.Fprintf(w, "<h1>Whoa, Go is neat!</h1>")
	fmt.Fprintf(w, "<title>Go</title>")

	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
	fmt.Fprintf(w, "<img src='assets/001.jpg' style='width:20%%;'>")
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
	err := http.ListenAndServe(":9091", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
