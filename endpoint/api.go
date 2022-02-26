package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Danice123/imagine/collection"
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

	tags := COLLECTIONHANDLER.Tags(image.RelativePath)
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

func DetectFaces(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if image.IsDir() {
		http.Error(w, "Api call on directory not valid", http.StatusBadRequest)
		return
	}

	faces, err := image.DetectFaces()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hd := COLLECTIONHANDLER.HashDirectory()
	for i, face := range faces {
		data := hd.Data(image.MD5())
		if data == nil {
			data = hd.CreateData(image.MD5())
		}

		if data.Faces == nil {
			data.Faces = map[string]collection.FaceBox{}
		}
		data.Faces["UNKNOWN"+strconv.Itoa(i)] = face
	}
	hd.Save()

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
