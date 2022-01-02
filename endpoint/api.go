package endpoint

import (
	"net/http"
	"os"
	"path/filepath"

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

	if err := os.Rename(image.FullPath, filepath.Join(COLLECTIONHANDLER.Trash(), filepath.Base(image.FullPath))); err != nil {
		panic(err)
	}

	tags := COLLECTIONHANDLER.Tags()
	tags.DeleteTag(image)
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

// func CleanImages(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
// 	tags, err := COLLECTIONHANDLER.Tags()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	for path := range tags.Mapping {
// 		if _, err := os.Stat(filepath.Join(ths.Root, path)); err != nil {
// 			if err := tags.DeleteFile(ths.Root, path); err != nil {
// 				panic(err.Error())
// 			}
// 		}
// 	}

// 	http.Redirect(w, req, req.Referer(), http.StatusFound)
// }
