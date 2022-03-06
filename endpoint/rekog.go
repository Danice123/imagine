package endpoint

import (
	"fmt"
	"net/http"
	"strconv"

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

func AddNewPerson(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if image.IsDir() {
		http.Error(w, "Api call on directory not valid", http.StatusBadRequest)
		return
	}

	faceId, err := strconv.Atoi(req.URL.Query().Get("face"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := req.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "must supply a name to tag with", http.StatusBadRequest)
		return
	}

	re := COLLECTIONHANDLER.RecognitionEngine()

	face, err := image.GetFaceImage(faceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = re.AddFace(face, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hd := COLLECTIONHANDLER.HashDirectory()
	data := hd.Data(image.MD5())
	data.Faces[faceId].Name = &name
	hd.Save()

	COLLECTIONHANDLER.Tags(image.RelativePath).WriteTag(image, name)

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func RecognizeFace(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if image.IsDir() {
		http.Error(w, "Api call on directory not valid", http.StatusBadRequest)
		return
	}

	faceId, err := strconv.Atoi(req.URL.Query().Get("face"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	re := COLLECTIONHANDLER.RecognitionEngine()

	face, err := image.GetFaceImage(faceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	matches, err := re.Search(face)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("%v\n", matches)
}
