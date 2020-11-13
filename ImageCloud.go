package main

//필요 패키지 임포트
import (
	"fmt"
	"image/jpeg"
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

func indexThumbnail(w http.ResponseWriter, r *http.Request) {

	// open "test.jpg"
	file, err := os.Open("assets/gopher.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(80, 0, img, resize.Lanczos3)

	out, err := os.Create("test_resized.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

}

func main() {
	//기본 Url 핸들러 메소드 지정
	http.HandleFunc("/", indexTandler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/thumbnail", indexThumbnail)
	//서버 시작
	err := http.ListenAndServe(":9090", nil)
	//예외 처리
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
