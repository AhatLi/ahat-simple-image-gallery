package main

//필요 패키지 임포트
import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

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

//FileSort  aaa
type FileSort []os.FileInfo

func (a FileSort) Len() int           { return len(a) }
func (a FileSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FileSort) Less(i, j int) bool { return a[i].IsDir() || a[i].Name() < a[j].Name() }

func indexTandler(w http.ResponseWriter, r *http.Request) {
	decodedValue, _ := url.QueryUnescape(r.URL.String())
	path := imgPath + decodedValue
	fmt.Println(imgPath + decodedValue)

	if decodedValue != "/" && !fileExists(path) {
		return
	}

	fmt.Fprintf(w, "<html><head><title>Ahat Simple Gallary</title><style>.modal {  display: none;  position: fixed;   z-index: 1;  padding-top: 100px;  left: 0;  top: 0;  width: 100%%;  height: 100%%;  overflow: auto;  background-color: rgb(0,0,0);  background-color: rgba(0,0,0,0.4);}.modal-content {  background-color: #fefefe;  margin: auto;  padding: 20px;  border: 1px solid #888;  width: 80%%; margin-bottom: 30%%;}</style></head><body><div id='myModal' class='modal'>  <div class='modal-content'>    <span class='close'>&times;</span>    <p><img src='https://blog.jinbo.net/attach/615/200937431.jpg' style='width: 100%%' id='myImg'></p>  </div></div>")

	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	sort.Sort(FileSort(files))

	for i, f := range files {
		if i%5 == 0 {
			fmt.Fprintf(w, "<br>")
		}

		fmt.Println("1: " + f.Name())
		fmt.Println("2: " + decodedValue)
		fmt.Println("3: " + imgPath)

		if f.IsDir() {
			fmt.Fprintf(w, "<a href=\""+decodedValue+"/"+f.Name()+"\"><img src='assets/directory.png' style='width:20%%;'></a>")
			fmt.Println("4: <a href=\"" + decodedValue + "/" + f.Name() + "\"><img src='assets/directory.png' style='width:20%%;'></a>")
		} else {
			fmt.Fprintf(w, "<img src='"+imgToBase64(thumPath+path+"/"+f.Name())+"' style='width:20%%;' onClick='thumbClick(\""+path+"/"+f.Name()+"\")' name='myBtn'>")
		}
	}
	fmt.Fprintf(w, "<script src='/assets/script.js'></script>")
	fmt.Fprintf(w, "</body></html>")
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(err)
		return false
	}
	return true
}

func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	if fileExists(thumbname) {
		return
	}

	fmt.Println(thumbname)

	// load original image
	img, err := imaging.Open(filename)

	if err != nil {
		fmt.Println(err)
		return
	}

	thumbnail := imaging.CropCenter(imaging.Resize(img, 80, 0, imaging.Lanczos), 80, 80)
	err = imaging.Save(thumbnail, thumbname)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func initServer() {
	t := time.Now()
	fmt.Println("A : ", t.Hour(), t.Minute(), t.Second(), t.Nanosecond())

	os.MkdirAll(imgPath, os.ModePerm)
	os.MkdirAll(thumPath, os.ModePerm)
	os.MkdirAll(assetPath, os.ModePerm)

	explorerDirectory(imgPath)

	t = time.Now()
	fmt.Println("B : ", t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

func explorerDirectory(path string) {
	os.MkdirAll(thumPath+path, os.ModePerm)
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for _, f := range files {

		if f.IsDir() {
			explorerDirectory(path + "/" + f.Name())
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
