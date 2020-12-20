package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"strings"

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

func imgToBase64(file string) string {
	f, _ := os.Open(file)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	f.Close()

	return "data:image/png;base64," + encoded
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

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

func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	if fileExists(thumbname + ".jpg") {
		return
	}

	img, err := imaging.Open(filename)

	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}
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
			//		fmt.Println(path + "/" + f.Name())
		} else if isImage(f.Name()) {
			makeThumbnail(path + "/" + f.Name())
		}
	}
}

func getUserData() (string, string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return "", ""
	}
	confUsername := cfg.Section("account").Key("username").String()
	confPasswd := cfg.Section("account").Key("passwd").String()

	return confUsername, confPasswd
}

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

func setContentData(count string, sort string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return
	}

	cfg.Section("content").Key("count").SetValue(count)
	cfg.Section("content").Key("sort").SetValue(sort)
	cfg.SaveTo("ImageCloud.conf")
}

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

func fileSearch(files []os.FileInfo, text string) []os.FileInfo {
	result := make([]os.FileInfo, 0)

	for _, f := range files {
		if strings.Contains(f.Name(), text) {
			result = append(result, f)
		}
	}

	return result
}

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

func isImage(filename string) bool {
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".gif") || strings.HasSuffix(filename, ".jpeg")
}
