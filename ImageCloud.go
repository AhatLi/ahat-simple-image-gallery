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
	"strings"

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

	if decodedValue != "/" && !fileExists(path) {
		return
	}

	data, _ := ioutil.ReadFile("./assets/index.html.ahat")
	str := fmt.Sprintf("%s", data)

	fmt.Fprintf(w, str[:strings.Index(str, "#content")])

	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	sort.Sort(FileSort(files))

	for i, f := range files {
		if (i+1)%5 == 0 {
			fmt.Fprintf(w, "</tr><tr>")
		}

		if i == 0 && decodedValue != "/" {
			fmt.Fprintf(w, "<td class='equalDivide'>")
			fmt.Fprintf(w, "<a href=\"..\"><img src='http://"+r.Host+"/assets/directory.png' style='width:100%%;'></a>")
			fmt.Fprintf(w, "<br>..")
			fmt.Fprintf(w, "</td>")
			i++
		}

		fmt.Fprintf(w, "<td class='equalDivide'>")
		if f.IsDir() {
			fmt.Fprintf(w, "<a href=\"."+decodedValue+"/"+f.Name()+"\"><img src='http://"+r.Host+"/assets/directory.png' style='width:100%%;'></a>")
			fmt.Fprintf(w, "<br>"+f.Name())
		} else {
			fmt.Fprintf(w, "<img src='"+imgToBase64(thumPath+path+"/"+f.Name())+"' style='width:100%%;' onClick='thumbClick(\"http://"+r.Host+"/"+path+"/"+f.Name()+"\")' name='myBtn'>")
			fmt.Fprintf(w, "<br>"+f.Name())
		}
		fmt.Fprintf(w, "</td>")
	}
	fmt.Fprintf(w, "</tr>")
	fmt.Fprintf(w, str[strings.LastIndex(str, "#content")+9:])
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	if fileExists(thumbname) {
		return
	}

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
			explorerDirectory(path + "/" + f.Name())
		} else {
			makeThumbnail(path + "/" + f.Name())
		}
	}
}

func main() {
	initServer()

	http.HandleFunc("/", indexTandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(imgPath))))
	http.Handle("/thumbnail/", http.StripPrefix("/thumbnail/", http.FileServer(http.Dir(thumPath))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))

	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("ListenAndServe Started! -> Port(9090)")
	}
}
