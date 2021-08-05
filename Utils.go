package main

// 2020-12-29
// Ahat Simple Gallery ver 0.9

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"gopkg.in/ini.v1"
)

//FileNameSort : img file sort for name
type FileNameSort []os.FileInfo

func (a FileNameSort) Len() int      { return len(a) }
func (a FileNameSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FileNameSort) Less(i, j int) bool {
	if a[i].IsDir() && a[j].IsDir() {
		return a[i].Name() < a[j].Name()
	} else if a[i].IsDir() {
		return true
	} else if a[j].IsDir() {
		return false
	}
	return a[i].Name() < a[j].Name()
}

//FileDateSort : img file sort for file last mod date
type FileDateSort []os.FileInfo

func (a FileDateSort) Len() int      { return len(a) }
func (a FileDateSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FileDateSort) Less(i, j int) bool {
	if a[i].IsDir() && a[j].IsDir() {
		return a[i].Name() < a[j].Name()
	} else if a[i].IsDir() {
		return true
	} else if a[j].IsDir() {
		return false
	}

	return a[i].ModTime().Unix() < a[j].ModTime().Unix()
}

//FileSizeSort : img file sort for file size
type FileSizeSort []os.FileInfo

func (a FileSizeSort) Len() int      { return len(a) }
func (a FileSizeSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FileSizeSort) Less(i, j int) bool {
	if a[i].IsDir() && a[j].IsDir() {
		return a[i].Name() < a[j].Name()
	} else if a[i].IsDir() {
		return true
	} else if a[j].IsDir() {
		return false
	}
	return a[i].Size() < a[j].Size()
}

//이미지를 base64 텍스트로 변환한다.
func imgToBase64(file string) string {
	f, _ := os.Open(file)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	f.Close()

	return "data:image/png;base64," + encoded
}

//파일 존재여부 확인
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//폴더 구조를 확인하는 함수
func getDirPath(path string, dir *[]string) {
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for _, f := range files {
		if f.IsDir() {
			*dir = append(*dir, path[8:]+"/"+f.Name())
			getDirPath(path+"/"+f.Name(), dir)
		}
	}
}

//썸네일을 생성한다. 현재 이미지의 특정색이 완전히 검은색인 경우 썸네일에서 흰색으로 표시되는 문제점이 있는것으로 보임.
func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	if fileExists(thumbname + ".jpg") {
		return
	}

	img, err := imaging.Open(filename)

	if err != nil {
		fmt.Println("makeThumbnail1 err : ", err)
		return
	}

	thumbnail := imaging.Thumbnail(img, 80, 80, imaging.Linear)
	thumbnail = imaging.AdjustFunc(
		thumbnail,
		func(c color.NRGBA) color.NRGBA {
			return color.NRGBA{c.R + uint8(255), c.G + uint8(255), c.B + uint8(255), uint8(255)}
		},
	)

	err = imaging.Save(thumbnail, thumbname+".jpg")

	if err != nil {
		fmt.Println("makeThumbnail2 err : ", err)
		return
	}
}

func preExplorerDirectory(path string) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg := new(sync.WaitGroup)
	wg.Add(1)
	go explorerDirectory(path, wg)

	wg.Wait()
}

//폴더를 탐색하여 이미지가 썸네일이 존재하지 않을 경우 썸네일 파일을 생성한다.
func explorerDirectory(path string, wg *sync.WaitGroup) {
	os.MkdirAll(thumPath+path, os.ModePerm)
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	for _, f := range files {

		if f.IsDir() {
			wg.Add(1)
			go explorerDirectory(path+"/"+f.Name(), wg)
		} else if isImage(f.Name()) {
			go makeThumbnail(path + "/" + f.Name())
		}
	}

	fmt.Println("explorerDirectory Done")
	wg.Done()
}

//유저 로그인을 위한 설정을 받아온다.
func getUserData() (string, string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return "", ""
	}
	confUsername := cfg.Section("account").Key("username").String()
	confPasswd := cfg.Section("account").Key("passwd").String()

	return confUsername, confPasswd
}

//환경 관련 설정을 받아온다.
func getEnvData() string {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return "assets/html/index.html.ahat"
	}
	html := cfg.Section("envronment").Key("html").String()

	if html == "" {
		return "assets/html/index.html.ahat"
	}

	return html
}

//이미지 표시를 위한 설저응ㄹ 받아온다.
func getContentData() (int, string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return 100, ""
	}
	count, _ := cfg.Section("content").Key("count").Int()
	sort := cfg.Section("content").Key("sort").String()

	if count == 0 {
		count = 100
	}

	if sort == "" {
		sort = "name"
	}

	return count, sort
}

//클라이언트가 설정을 변경하였을 경우 해당 설정을 파일에 반영한다.
func setContentData(count string, sort string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return
	}

	cfg.Section("content").Key("count").SetValue(count)
	cfg.Section("content").Key("sort").SetValue(sort)
	cfg.SaveTo("ImageCloud.conf")
}

//설정파일에서 포트번호를 받아온다. 기본 포트주소는 9090
func getPort() string {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return ":9090"
	}
	port := cfg.Section("network").Key("port").String()
	if port == "" {
		return ":9090"
	}

	return ":" + cfg.Section("network").Key("port").String()
}

//파일을 검색했을 경우 해당 텍스트로 파일을 필터링한다. 추후 와일드카드 등 추가 예정
func fileSearch(files []os.FileInfo, text string) []os.FileInfo {
	result := make([]os.FileInfo, 0)

	for _, f := range files {
		if strings.Contains(f.Name(), text) {
			result = append(result, f)
		}
	}

	return result
}

//이미지 파일만을 표시하기 위해서 필터링을 하는 동작
func imgFilter(files []os.FileInfo, text string) []os.FileInfo {
	result := make([]os.FileInfo, 0)

	if text == "" {
		for _, f := range files {
			if f.IsDir() || isImage(f.Name()) {
				result = append(result, f)
			}
		}
	} else {
		for _, f := range files {
			if f.IsDir() || strings.Contains(f.Name(), text) && isImage(f.Name()) {
				result = append(result, f)
			}
		}
	}

	return result
}

//현재는 네가지 파일 확장자만을 지원한다.
func isImage(filename string) bool {
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".gif") || strings.HasSuffix(filename, ".jpeg")
}

//엑세스 로그를 출력하는 함수
func printLog(r *http.Request) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + "," + r.RemoteAddr + "," + r.URL.Path)
}

//파일 삭제 동작을 하는 함수
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

//파일 이동 동작을 하는 함수
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
