package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

var imgPath string = "images"
var thumPath string = "thumbnail"
const assetPath string = "assets"

var htmlFile string = "assets/html/index.html.ahat"

func main() {
	loadPathConfig()
	startTime := time.Now()
	go initServer()
	elapsedTime := time.Since(startTime)

	fmt.Println("startup seconds:", elapsedTime.Seconds())
	fmt.Println("init complete. server start")

	http.HandleFunc("/main", indexPageHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	http.HandleFunc("/images/", protectedFileHandler("/images/", imgPath, isImage))
	http.HandleFunc("/thumbnail/", protectedFileHandler("/thumbnail/", thumPath, isThumbnailFile))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(assetPath))))
	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/", indexHandler)

	server := &http.Server{
		Addr:              getPort(),
		Handler:           nil,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("ListenAndServe: " + err.Error())
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	printLog(r)
	if !loginCheck(w, r) {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := os.ReadFile(htmlFile)
	if err != nil {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	str := string(data)

	cfg := getContentData()
	page, _ := strconv.Atoi(r.URL.Query().Get("p"))
	if page < 1 {
		page = 1
	}
	search := strings.TrimSpace(r.URL.Query().Get("searchText"))

	dirPath, err := resolveImageDir(r.URL.Path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	infos := make([]os.FileInfo, 0, len(files))
	for _, entry := range files {
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		infos = append(infos, info)
	}

	infos = imgFilter(infos, search)

	switch cfg.Sort {
	case "date":
		sort.Sort(FileDateSort(infos))
	case "size":
		sort.Sort(FileSizeSort(infos))
	default:
		sort.Sort(FileNameSort(infos))
	}

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			switch line {
			case "#select":
				jsCfg, _ := json.Marshal(cfg)
				fmt.Fprintf(w, "<script>const contentConfig=%s;document.documentElement.style.setProperty('--mobile-columns', String(contentConfig.mobileColumns));</script>", jsCfg)
				makeSelect(w)
			case "#content":
				makeContent(w, r, cfg, page, infos)
			case "#page":
				makePage(w, cfg.Count, page, len(infos))
			}
			continue
		}
		io.WriteString(w, line)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	printLog(r)
	if !loginCheck(w, r) {
		return
	}

	switch {
	case strings.HasSuffix(r.URL.Path, "/config") && r.Method == http.MethodGet:
		writeConfigResponse(w, getContentData())
	case strings.HasSuffix(r.URL.Path, "/move"):
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := fileMove(r.PostFormValue("files"), r.PostFormValue("source"), r.PostFormValue("dest"))
		writeAPIResponse(w, err)
	case strings.HasSuffix(r.URL.Path, "/delete"):
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := fileRemove(r.PostFormValue("files"), r.PostFormValue("path"))
		writeAPIResponse(w, err)
	case strings.HasSuffix(r.URL.Path, "/config"):
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err := setContentData(r.PostFormValue("imgCount"), r.PostFormValue("imgSort"), r.PostFormValue("mobileColumns"))
		writeAPIResponse(w, err)
	default:
		http.NotFound(w, r)
	}
}

func writeAPIResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"ok":false,"error":%q}`, err.Error())
		return
	}
	io.WriteString(w, `{"ok":true}`)
}

func writeConfigResponse(w http.ResponseWriter, cfg ContentConfig) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(cfg)
}

func initServer() {
	_ = os.MkdirAll(imgPath, os.ModePerm)
	_ = os.MkdirAll(thumPath, os.ModePerm)
	_ = os.MkdirAll(assetPath, os.ModePerm)

	htmlFile = getEnvData()

	for {
		preExplorerDirectory(imgPath)
		time.Sleep(time.Second * 60)
	}
}

func makeContent(w http.ResponseWriter, r *http.Request, cfg ContentConfig, page int, files []os.FileInfo) {
	totalPages := pageCount(len(files), cfg.Count)
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * cfg.Count
	end := start + cfg.Count
	if start > len(files) {
		start = len(files)
	}
	if end > len(files) {
		end = len(files)
	}
	files = files[start:end]

	currentDir := safeURLDirPath(r.URL.Path)
	parentDir := parentURLDirPath(currentDir)

	fmt.Fprint(w, "<div class='gallery-grid'>")
	fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv directory-tile' style='background-image: url("/assets/images/directory.png")' onClick='location.href=%q'></div><div class='imgtitle'>..</div></div>`, parentDir)

	gallery := make([]map[string]string, 0, len(files))
	galleryIndex := 0

	for _, f := range files {
		name := f.Name()
		escapedName := html.EscapeString(name)
		if f.IsDir() {
			nextDir := joinURLDir(currentDir, name)
			fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv directory-tile' style='background-image: url("/assets/images/directory.png")' onClick='location.href=%q'></div><div class='imgtitle'>%s</div></div>`, nextDir, escapedName)
			continue
		}

		imageURL := fileURL("/images", currentDir, name)
		thumbURL := fileURL("/thumbnail", currentDir, name+".jpg")
		fmt.Fprintf(w, `<div class='imgBox'><div class='imgDiv' style='background-image: url(%q)' onClick='showGallery(%d, this)' data-file-name=%q data-file-url=%q data-dir-path=%q id='img%d' ontouchstart='longTouch(this.id)' ontouchend='revert(this.id)'></div><div class='imgtitle'>%s</div></div>`, thumbURL, galleryIndex, name, imageURL, currentDir, galleryIndex, escapedName)
		gallery = append(gallery, map[string]string{"title": name, "src": imageURL})
		galleryIndex++
	}

	fmt.Fprint(w, "</div>")
	galleryJSON, _ := json.Marshal(gallery)
	fmt.Fprintf(w, "<script>const gallery=%s;</script>", galleryJSON)
}

func makePage(w http.ResponseWriter, count int, page int, total int) {
	totalPages := pageCount(total, count)
	fmt.Fprintf(w, "<script>var page=%d;var count=%d;</script>", page, count)
	fmt.Fprint(w, "<select class='form-control' style='background-color: #6c757d; color: white;' id='pageSelect'>")
	for i := 1; i <= totalPages; i++ {
		if page == i {
			fmt.Fprintf(w, "<option selected>%d</option>", i)
		} else {
			fmt.Fprintf(w, "<option>%d</option>", i)
		}
	}
	fmt.Fprint(w, "</select>")
}

func makeSelect(w http.ResponseWriter) {
	dirs := []string{"/"}
	getDirPath(imgPath, &dirs)

	fmt.Fprint(w, "<select class='form-select' id='selectDir'>")
	fmt.Fprint(w, "<option value='/'>/</option>")
	for _, d := range dirs {
		if d == "/" {
			continue
		}
		fmt.Fprintf(w, "<option value=%q>%s</option>", d, html.EscapeString(d))
	}
	fmt.Fprint(w, "</select>")
}

func pageCount(total int, count int) int {
	if count <= 0 {
		count = 100
	}
	if total == 0 {
		return 1
	}
	return (total + count - 1) / count
}

func safeURLDirPath(raw string) string {
	cleaned := path.Clean("/" + strings.TrimSpace(raw))
	if cleaned == "." || cleaned == "" {
		return "/"
	}
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}
	if cleaned != "/" && !strings.HasSuffix(cleaned, "/") {
		cleaned += "/"
	}
	return cleaned
}

func parentURLDirPath(raw string) string {
	cleaned := safeURLDirPath(raw)
	if cleaned == "/" {
		return "/"
	}
	trimmed := strings.TrimSuffix(cleaned, "/")
	parent := path.Dir(trimmed)
	if parent == "." {
		return "/"
	}
	if !strings.HasSuffix(parent, "/") {
		parent += "/"
	}
	return parent
}

func joinURLDir(baseDir string, name string) string {
	cleanName := strings.Trim(name, "/")
	joined := path.Join(strings.TrimSuffix(safeURLDirPath(baseDir), "/"), cleanName)
	if !strings.HasPrefix(joined, "/") {
		joined = "/" + joined
	}
	return joined + "/"
}

func fileURL(prefix string, dir string, name string) string {
	base := strings.TrimSuffix(prefix, "/")
	currentDir := safeURLDirPath(dir)
	currentDir = strings.TrimPrefix(currentDir, "/")
	currentDir = strings.TrimSuffix(currentDir, "/")

	parts := []string{base}
	if currentDir != "" {
		for _, part := range strings.Split(currentDir, "/") {
			if part == "" {
				continue
			}
			parts = append(parts, url.PathEscape(part))
		}
	}
	parts = append(parts, url.PathEscape(name))
	return strings.Join(parts, "/")
}

func resolveImageDir(urlPath string) (string, error) {
	rel := strings.TrimPrefix(safeURLDirPath(urlPath), "/")
	rel = strings.TrimSuffix(rel, "/")
	return safeJoinUnderBase(imgPath, rel)
}

func protectedFileHandler(prefix string, baseDir string, allow func(string) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		printLog(r)
		if !loginCheck(w, r) {
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		rel, err := requestRelativePath(r.URL.Path, prefix)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if !allow(path.Base(rel)) {
			http.NotFound(w, r)
			return
		}

		fullPath, err := safeJoinUnderBase(baseDir, rel)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(fullPath)
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, fullPath)
	}
}

func requestRelativePath(requestPath string, prefix string) (string, error) {
	if !strings.HasPrefix(requestPath, prefix) {
		return "", errors.New("invalid path")
	}

	decoded, err := url.PathUnescape(strings.TrimPrefix(requestPath, prefix))
	if err != nil {
		return "", err
	}

	cleaned := path.Clean("/" + strings.TrimSpace(decoded))
	if cleaned == "/" || cleaned == "." {
		return "", errors.New("invalid path")
	}

	return strings.TrimPrefix(cleaned, "/"), nil
}
