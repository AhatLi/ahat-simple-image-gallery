package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

//FileSort  aaa
type FileSort []os.FileInfo

func (a FileSort) Len() int      { return len(a) }
func (a FileSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FileSort) Less(i, j int) bool {
	if a[i].IsDir() && a[j].IsDir() {
		return a[i].Name() < a[j].Name()
	} else if a[i].IsDir() {
		return true
	} else if a[j].IsDir() {
		return false
	}
	return a[i].Name() < a[j].Name()
}

func imgToBase64(file string) string {
	f, _ := os.Open(file)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)

	return "data:image/png;base64," + encoded
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getAllDirPath() {
	var dir []string
	getDirPath("./images", &dir, 0)

	for _, d := range dir {
		fmt.Println(d)
	}

}

func getDirPath(path string, dir *[]string, level int) {
	files, errf := ioutil.ReadDir(path)
	if errf != nil {
		log.Fatal(errf)
	}

	b := ""
	for i := 0; i < level; i++ {
		b += "  "
	}

	for _, f := range files {
		if f.IsDir() {
			*dir = append(*dir, b+f.Name())
			getDirPath(path+"/"+f.Name(), dir, level+1)
		}
	}
}

func makeThumbnail(filename string) {
	thumbname := thumPath + filename

	if fileExists(thumbname) {
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

	//	imaging.Encode(os.Stdout, thumbnail, imaging.JPEG)
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
		} else {
			makeThumbnail(path + "/" + f.Name())
		}
	}
}
