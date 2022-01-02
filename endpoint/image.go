package endpoint

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RawImage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() || image.IsDir() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, req, image.FullPath)
}
