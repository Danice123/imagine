package endpoint

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
	data := hd.Data(image.MD5())
	if data == nil {
		data = hd.CreateData(image.MD5())
	}
	data.Faces = faces
	hd.Save()

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
