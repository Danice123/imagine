package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func ToggleTag(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if image.IsDir() {
		http.Error(w, "Tagging folders not allowed", http.StatusBadRequest)
		return
	}

	tagName := req.URL.Query().Get("name")
	if tagName == "" {
		return
	}

	tags := COLLECTIONHANDLER.Tags()
	tags.WriteTag(image, tagName)
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func TrashImage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if image.IsDir() {
		http.Error(w, "Trashing folders not allowed", http.StatusBadRequest)
		return
	}

	COLLECTIONHANDLER.HashCache().RemoveHash(image)
	if err := os.Rename(image.FullPath, filepath.Join(COLLECTIONHANDLER.Trash(), filepath.Base(image.FullPath))); err != nil {
		panic(err)
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func ChangeSeries(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	hash := req.URL.Query().Get("hash")
	if hash == "" {
		return
	}

	series := req.URL.Query().Get("series")
	sm := COLLECTIONHANDLER.Series()
	sm.AddImageToSeries(hash, series)
	sm.Write()

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func CleanImages(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	hc := COLLECTIONHANDLER.HashCache()
	var count int
	for _, path := range hc.CachedImages() {
		image := COLLECTIONHANDLER.Image(path)
		if !image.IsValid() {
			hc.RemoveHash(image)
			count++
		}
	}
	w.Write([]byte(strconv.Itoa(count)))
}
