package endpoint

import (
	"bytes"
	"image/png"
	"net/http"
	"strconv"
	"time"

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

func RawFace(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() || image.IsDir() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	faceId, err := strconv.Atoi(req.URL.Query().Get("face"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	face, err := image.GetFaceImage(faceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buffer bytes.Buffer
	err = png.Encode(&buffer, face)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, req, "test", time.Now(), bytes.NewReader(buffer.Bytes()))
}
