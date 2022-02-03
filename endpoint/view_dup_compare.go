package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type DupCompareData struct {
	Hash   string
	Type   string
	Images []DupImage
}

type DupImage struct {
	Path    string
	IsVideo bool
}

func DupCompare(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	params := req.URL.Query()
	if params.Get("hash") == "" || params.Get("type") == "" {
		http.Error(w, "Bad parameters", http.StatusNotFound)
		return
	}

	var imageList []string
	switch params.Get("type") {
	case "MD5":
		imageList = COLLECTIONHANDLER.HashCache().GetDups()[params.Get("hash")]
	case "PerceptionHash":
		hashList := COLLECTIONHANDLER.HashDirectory().GetPHashDups()[params.Get("hash")]
		imageList = []string{}
		for _, hash := range hashList {
			imageList = append(imageList, COLLECTIONHANDLER.HashCache().GetImagePathByHash(hash))
		}
	default:
		http.Error(w, "Bad hash type", http.StatusNotFound)
		return
	}

	images := []DupImage{}
	for _, path := range imageList {
		image := COLLECTIONHANDLER.Image(path)
		images = append(images, DupImage{
			Path:    path,
			IsVideo: image.IsVideo(),
		})
	}

	var tagTemplate = template.New("DupCompare")
	if html, err := os.ReadFile(filepath.Join("templates", "dupcompare.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := tagTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			tagTemplate.Execute(w, DupCompareData{
				Hash:   params.Get("hash"),
				Type:   params.Get("type"),
				Images: images,
			})
		}
	}
}

// func MarkAsNotDup(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
// 	params := req.URL.Query()
// 	if params.Get("hash") == "" || params.Get("type") == "" {
// 		http.Error(w, "Bad parameters", http.StatusNotFound)
// 		return
// 	}

// 	tags, err := collection.New(ths.Root)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	images := []string{}
// 	for image, imageData := range tags.Mapping {
// 		if params.Get("type") == "MD5" && imageData.MD5 == params.Get("hash") {
// 			images = append(images, image)
// 		} else if params.Get("type") == "PerceptionHash" && imageData.PHash == params.Get("hash") {
// 			images = append(images, image)
// 		}
// 	}

// 	err = tags.SetDupList(ths.Root, params.Get("hash"), images)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	http.Redirect(w, req, "/dups", http.StatusFound)
// }
