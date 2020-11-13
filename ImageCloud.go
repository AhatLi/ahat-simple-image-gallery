package main

//필요 패키지 임포트
import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

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
	thumbname := "./thumbnail/" + filename
	ext := thumbname[len(thumbname)-4:]
	fmt.Println(ext)
	if fileExists(thumbname) {
		return
	}

	file, err1 := os.Open(filename)
	if err1 != nil {
		log.Fatal(err1)
	}

	switch ext {
	case ".jpg", "jpeg":
		{
			fmt.Println("jpg")

			img, err2 := jpeg.Decode(file)
			if err2 != nil {
				log.Fatal(err2)
			}
			file.Close()

			m := resize.Resize(80, 0, img, resize.Lanczos3)
			out, err3 := os.Create(thumbname)
			if err3 != nil {
				log.Fatal(err3)
			}
			defer out.Close()
			jpeg.Encode(out, m, nil)
		}
	case ".png":
		{
			fmt.Println("png")

			img, err2 := png.Decode(file)
			if err2 != nil {
				log.Fatal(err2)
			}
			file.Close()

			m := resize.Resize(80, 0, img, resize.Lanczos3)
			out, err3 := os.Create(thumbname)
			if err3 != nil {
				log.Fatal(err3)
			}
			defer out.Close()
			png.Encode(out, m)
		}
	}

}

func initServer() {
	path := "./images"
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for _, f := range files {
		makeThumbnail(path + "/" + f.Name())
	}
}

func main() {
	initServer()

	//기본 Url 핸들러 메소드 지정
	http.HandleFunc("/", indexTandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	//서버 시작
	err := http.ListenAndServe(":9091", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
