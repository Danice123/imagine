package main

import (
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/julienschmidt/httprouter"
)

var root string

func raw(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	http.ServeFile(w, req, filepath.Join(root, ps.ByName("path")))
}

func toggleRandom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if _, err := req.Cookie("random"); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "random",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:  "random",
			Value: fmt.Sprint(time.Now().UnixNano()),
			Path:  "/",
		})
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func image(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	data := &ImageData{
		Url:         "/raw/" + ps.ByName("path"),
		HasNext:     false,
		HasPrevious: false,
		RandomState: "Off",
	}

	isRandom := false
	if cookie, err := req.Cookie("random"); err == nil {
		isRandom = true
		data.RandomState = "On"
		seed, _ := strconv.ParseInt(cookie.Value, 10, 64)
		rand.Seed(seed)
	}

	specificFile := true
	baseDir := filepath.Dir(filepath.Join(root, ps.ByName("path")))
	if fileInfo, err := os.Stat(filepath.Join(root, ps.ByName("path"))); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	} else if fileInfo.IsDir() {
		baseDir = filepath.Join(root, ps.ByName("path"))
		specificFile = false
	}
	mime, _ := mimetype.DetectFile(filepath.Join(root, ps.ByName("path")))
	data.Mimetype = mime.String()
	data.IsVideo = strings.Split(mime.String(), "/")[0] == "video"

	if files, err := os.ReadDir(baseDir); err != nil {
		panic(err.Error())
	} else {
		if isRandom {
			rand.Shuffle(len(files), func(i, j int) {
				files[i], files[j] = files[j], files[i]
			})
		} else {
			sort.Slice(files, func(i, j int) bool {
				return strings.Compare(files[i].Name(), files[j].Name()) < 0
			})
		}

		if !specificFile {
			http.Redirect(w, req, strings.TrimSuffix(req.RequestURI, "/")+"/"+files[0].Name(), http.StatusFound)
			return
		}
		for n, file := range files {
			if filepath.Join(baseDir, file.Name()) == filepath.Join(root, ps.ByName("path")) {
				data.Next, data.HasNext = findNextFile(files, n)
				data.Previous, data.HasPrevious = findPreviousFile(files, n)
			}
		}
	}

	imageTemplate.Execute(w, data)
}

func findNextFile(files []fs.DirEntry, current int) (string, bool) {
	n := current + 1
	if n >= len(files) {
		return "", false
	}
	for files[n].IsDir() {
		n += 1
		if n >= len(files) {
			return "", false
		}
	}
	return files[n].Name(), true
}

func findPreviousFile(files []fs.DirEntry, current int) (string, bool) {
	n := current - 1
	if n < 0 {
		return "", false
	}
	for files[n].IsDir() {
		n -= 1
		if n < 0 {
			return "", false
		}
	}
	return files[n].Name(), true
}

func main() {
	root = os.Args[1]
	router := httprouter.New()
	router.GET("/raw/*path", raw)
	router.GET("/browse/*path", image)
	router.GET("/api/random", toggleRandom)
	http.ListenAndServe(":8080", router)
}
