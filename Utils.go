package main

import (
	"fmt"
	"image/color"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
	"gopkg.in/ini.v1"
)

var (
	fileMap   = make(map[string]int64)
	fileMapMu sync.RWMutex
)

type movedFile struct {
	oldImagePath string
	newImagePath string
	oldThumbPath string
	newThumbPath string
	thumbMoved   bool
	size         int64
	hadSize      bool
}

type ContentConfig struct {
	Count         int    `json:"count"`
	Sort          string `json:"sort"`
	MobileColumns int    `json:"mobileColumns"`
}

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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getDirPath(base string, dir *[]string) {
	entries, err := os.ReadDir(base)
	if err != nil {
		fmt.Println("getDirPath error:", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(base, entry.Name())
		rel, relErr := filepath.Rel(imgPath, fullPath)
		if relErr != nil {
			fmt.Println("getDirPath rel error:", relErr)
			continue
		}

		normalized := "/" + filepath.ToSlash(rel)
		*dir = append(*dir, normalized)
		getDirPath(fullPath, dir)
	}
}

func getFileSize(filename string) int64 {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0
	}
	return fi.Size()
}

func thumbnailDirectory(imageDir string) (string, error) {
	rel, err := filepath.Rel(imgPath, imageDir)
	if err != nil {
		return "", err
	}
	if rel == "." {
		rel = ""
	}
	return safeJoinUnderBase(thumPath, filepath.ToSlash(rel))
}

func thumbnailPath(imagePath string) (string, error) {
	rel, err := filepath.Rel(imgPath, imagePath)
	if err != nil {
		return "", err
	}
	thumbBase, err := safeJoinUnderBase(thumPath, filepath.ToSlash(rel))
	if err != nil {
		return "", err
	}
	return thumbBase + ".jpg", nil
}

func preExplorerDirectory(root string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	explorerDirectory(root)
	time.Sleep(time.Millisecond * 10)
}

func explorerDirectory(root string) {
	targetDir, err := thumbnailDirectory(root)
	if err != nil {
		fmt.Println("thumbnailDirectory error:", err)
		return
	}
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		fmt.Println("mkdir thumbnail error:", err)
		return
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		fmt.Println("explorerDirectory error:", err)
		return
	}

	for _, entry := range entries {
		fullPath := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			explorerDirectory(fullPath)
			continue
		}
		if !isImage(entry.Name()) {
			continue
		}

		size := getFileSize(fullPath)
		if shouldMakeThumbnail(fullPath, size) {
			makeThumbnail(fullPath, size)
		} else {
			time.Sleep(time.Millisecond * 5)
		}
	}
}

func shouldMakeThumbnail(filename string, size int64) bool {
	thumbname, err := thumbnailPath(filename)
	if err != nil {
		fmt.Println("thumbnailPath error:", err)
		return false
	}

	fileMapMu.RLock()
	previous, ok := fileMap[filename]
	fileMapMu.RUnlock()

	if !ok {
		return !fileExists(thumbname)
	}
	if previous != size {
		return true
	}
	return !fileExists(thumbname)
}

func makeThumbnail(filename string, size int64) {
	thumbname, err := thumbnailPath(filename)
	if err != nil {
		fmt.Println("thumbnailPath error:", err)
		return
	}

	fileMapMu.Lock()
	fileMap[filename] = size
	fileMapMu.Unlock()

	img, err := imaging.Open(filename)
	if err != nil {
		fmt.Println("makeThumbnail open error:", err)
		return
	}

	thumbnail := imaging.Thumbnail(img, 80, 80, imaging.Linear)
	thumbnail = imaging.AdjustFunc(thumbnail, func(c color.NRGBA) color.NRGBA {
		return color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255}
	})

	if err := os.MkdirAll(filepath.Dir(thumbname), os.ModePerm); err != nil {
		fmt.Println("makeThumbnail mkdir error:", err)
		return
	}
	if err := imaging.Save(thumbnail, thumbname); err != nil {
		fmt.Println("makeThumbnail save error:", err)
	}
}

func getUserData() (string, string) {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return "", ""
	}
	return cfg.Section("account").Key("username").String(), cfg.Section("account").Key("passwd").String()
}

func getEnvData() string {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return filepath.ToSlash(filepath.Join(assetPath, "html", "index.html.ahat"))
	}

	html := cfg.Section("envronment").Key("html").String()
	if html == "" {
		return filepath.ToSlash(filepath.Join(assetPath, "html", "index.html.ahat"))
	}
	return html
}

func getContentData() ContentConfig {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return ContentConfig{Count: 100, Sort: "name", MobileColumns: 4}
	}

	count, _ := cfg.Section("content").Key("count").Int()
	if count <= 0 {
		count = 100
	}

	sortValue := cfg.Section("content").Key("sort").String()
	if sortValue != "name" && sortValue != "date" && sortValue != "size" {
		sortValue = "name"
	}

	mobileColumns, _ := cfg.Section("content").Key("mobile_columns").Int()
	if mobileColumns < 2 || mobileColumns > 6 {
		mobileColumns = 4
	}

	return ContentConfig{Count: count, Sort: sortValue, MobileColumns: mobileColumns}
}

func setContentData(count string, sortValue string, mobileColumns string) error {
	cfg, err := ini.LooseLoad("ImageCloud.conf")
	if err != nil {
		return err
	}

	parsedCount, err := strconv.Atoi(count)
	if err != nil || parsedCount <= 0 {
		return fmt.Errorf("invalid image count")
	}
	if sortValue != "name" && sortValue != "date" && sortValue != "size" {
		return fmt.Errorf("invalid sort option")
	}
	parsedMobileColumns, err := strconv.Atoi(mobileColumns)
	if err != nil || parsedMobileColumns < 2 || parsedMobileColumns > 6 {
		return fmt.Errorf("invalid mobile column count")
	}

	section := cfg.Section("content")
	section.Key("count").SetValue(strconv.Itoa(parsedCount))
	section.Key("sort").SetValue(sortValue)
	section.Key("mobile_columns").SetValue(strconv.Itoa(parsedMobileColumns))

	return cfg.SaveTo("ImageCloud.conf")
}

func getPort() string {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return ":9090"
	}

	port := strings.TrimSpace(cfg.Section("network").Key("port").String())
	if port == "" {
		return ":9090"
	}
	return ":" + port
}

func loadPathConfig() {
	cfg, err := ini.Load("ImageCloud.conf")
	if err != nil {
		return
	}

	images := strings.TrimSpace(cfg.Section("path").Key("images").String())
	if images != "" {
		imgPath = images
	}

	thumbnail := strings.TrimSpace(cfg.Section("path").Key("thumbnail").String())
	if thumbnail != "" {
		thumPath = thumbnail
	}
}

func imgFilter(files []os.FileInfo, text string) []os.FileInfo {
	result := make([]os.FileInfo, 0, len(files))
	query := strings.ToLower(strings.TrimSpace(text))

	for _, f := range files {
		if !f.IsDir() && !isImage(f.Name()) {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(f.Name()), query) {
			continue
		}
		result = append(result, f)
	}

	return result
}

func isImage(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".png") || strings.HasSuffix(lower, ".gif") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".webp")
}

func isThumbnailFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".jpg")
}

func printLog(r *http.Request) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + "," + r.RemoteAddr + "," + r.Method + "," + r.URL.Path)
}

func fileRemove(files string, dir string) error {
	sourceDir, err := safeImageSubdir(dir)
	if err != nil {
		return err
	}

	names, err := parseFileList(files)
	if err != nil {
		return err
	}

	for _, name := range names {
		imagePath := filepath.Join(sourceDir, name)
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			return err
		}

		thumbPath, err := thumbnailPath(imagePath)
		if err != nil {
			return err
		}
		if err := os.Remove(thumbPath); err != nil && !os.IsNotExist(err) {
			return err
		}

		fileMapMu.Lock()
		delete(fileMap, imagePath)
		fileMapMu.Unlock()
	}

	return nil
}

func fileMove(files string, source string, dest string) error {
	sourceDir, err := safeImageSubdir(source)
	if err != nil {
		return err
	}
	destDir, err := safeImageSubdir(dest)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return err
	}
	destThumbDir, err := thumbnailDirectory(destDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destThumbDir, os.ModePerm); err != nil {
		return err
	}

	names, err := parseFileList(files)
	if err != nil {
		return err
	}

	moves := make([]movedFile, 0, len(names))

	for _, name := range names {
		oldImagePath := filepath.Join(sourceDir, name)
		newImagePath := filepath.Join(destDir, name)
		if err := os.Rename(oldImagePath, newImagePath); err != nil {
			if rollbackErr := rollbackMovedFiles(moves); rollbackErr != nil {
				return fmt.Errorf("%v (rollback failed: %w)", err, rollbackErr)
			}
			return err
		}

		oldThumbPath, err := thumbnailPath(oldImagePath)
		if err != nil {
			_ = os.Rename(newImagePath, oldImagePath)
			if rollbackErr := rollbackMovedFiles(moves); rollbackErr != nil {
				return fmt.Errorf("%v (rollback failed: %w)", err, rollbackErr)
			}
			return err
		}
		newThumbPath, err := thumbnailPath(newImagePath)
		if err != nil {
			_ = os.Rename(newImagePath, oldImagePath)
			if rollbackErr := rollbackMovedFiles(moves); rollbackErr != nil {
				return fmt.Errorf("%v (rollback failed: %w)", err, rollbackErr)
			}
			return err
		}
		if err := os.MkdirAll(filepath.Dir(newThumbPath), os.ModePerm); err != nil {
			_ = os.Rename(newImagePath, oldImagePath)
			if rollbackErr := rollbackMovedFiles(moves); rollbackErr != nil {
				return fmt.Errorf("%v (rollback failed: %w)", err, rollbackErr)
			}
			return err
		}
		thumbMoved := false
		if err := os.Rename(oldThumbPath, newThumbPath); err != nil && !os.IsNotExist(err) {
			_ = os.Rename(newImagePath, oldImagePath)
			if rollbackErr := rollbackMovedFiles(moves); rollbackErr != nil {
				return fmt.Errorf("%v (rollback failed: %w)", err, rollbackErr)
			}
			return err
		} else if err == nil {
			thumbMoved = true
		}

		move := movedFile{
			oldImagePath: oldImagePath,
			newImagePath: newImagePath,
			oldThumbPath: oldThumbPath,
			newThumbPath: newThumbPath,
			thumbMoved:   thumbMoved,
		}

		fileMapMu.Lock()
		size, ok := fileMap[oldImagePath]
		if ok {
			delete(fileMap, oldImagePath)
			fileMap[newImagePath] = size
			move.size = size
			move.hadSize = true
		}
		fileMapMu.Unlock()

		moves = append(moves, move)
	}

	return nil
}

func rollbackMovedFiles(moves []movedFile) error {
	for i := len(moves) - 1; i >= 0; i-- {
		move := moves[i]

		if err := os.Rename(move.newImagePath, move.oldImagePath); err != nil && !os.IsNotExist(err) {
			return err
		}
		if move.thumbMoved {
			if err := os.MkdirAll(filepath.Dir(move.oldThumbPath), os.ModePerm); err != nil {
				return err
			}
			if err := os.Rename(move.newThumbPath, move.oldThumbPath); err != nil && !os.IsNotExist(err) {
				return err
			}
		}

		fileMapMu.Lock()
		delete(fileMap, move.newImagePath)
		if move.hadSize {
			fileMap[move.oldImagePath] = move.size
		}
		fileMapMu.Unlock()
	}

	return nil
}

func parseFileList(files string) ([]string, error) {
	raw := strings.Split(files, ",")
	names := make([]string, 0, len(raw))

	for _, file := range raw {
		name := strings.TrimSpace(file)
		if name == "" {
			continue
		}
		if filepath.Base(name) != name || strings.Contains(name, "/") || strings.Contains(name, `\`) {
			return nil, fmt.Errorf("invalid file name")
		}
		names = append(names, name)
	}

	if len(names) == 0 {
		return nil, fmt.Errorf("no files selected")
	}
	return names, nil
}

func safeImageSubdir(dir string) (string, error) {
	cleaned := strings.TrimSpace(dir)
	cleaned = strings.TrimPrefix(cleaned, "/")
	cleaned = filepath.ToSlash(cleaned)
	cleaned = strings.TrimSuffix(cleaned, "/")
	return safeJoinUnderBase(imgPath, cleaned)
}

func safeJoinUnderBase(base string, rel string) (string, error) {
	baseClean := filepath.Clean(base)
	candidate := filepath.Clean(filepath.Join(baseClean, filepath.FromSlash(rel)))

	if candidate == baseClean {
		return candidate, nil
	}

	prefix := baseClean + string(os.PathSeparator)
	if !strings.HasPrefix(candidate, prefix) {
		return "", fmt.Errorf("invalid path")
	}
	return candidate, nil
}
