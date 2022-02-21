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

func getImageList(hashType string, hash string) []string {
	switch hashType {
	case "MD5":
		return COLLECTIONHANDLER.HashCache().GetDups()[hash]
	case "PerceptionHash":
		hashList := COLLECTIONHANDLER.HashDirectory().GetPHashDups()[hash]
		imageList := []string{}
		for _, hash := range hashList {
			imageList = append(imageList, COLLECTIONHANDLER.HashCache().GetImagePathByHash(hash))
		}
		return imageList
	default:
		return nil
	}
}

func DupCompare(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	params := req.URL.Query()
	if params.Get("hash") == "" || params.Get("type") == "" {
		http.Error(w, "Bad parameters", http.StatusNotFound)
		return
	}

	imageList := getImageList(params.Get("type"), params.Get("hash"))
	if imageList == nil {
		http.Error(w, "Bad hash type", http.StatusNotFound)
		return
	}

	images := []DupImage{}
	for _, path := range imageList {
		image := COLLECTIONHANDLER.Image(path)
		if image.IsValid() {
			images = append(images, DupImage{
				Path:    path,
				IsVideo: image.IsVideo(),
			})
		}
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

func MarkAsNotDup(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	params := req.URL.Query()
	if params.Get("hash") == "" || params.Get("type") == "" {
		http.Error(w, "Bad parameters", http.StatusNotFound)
		return
	}

	imageList := getImageList(params.Get("type"), params.Get("hash"))
	if imageList == nil {
		http.Error(w, "Bad hash type", http.StatusNotFound)
		return
	}

	COLLECTIONHANDLER.Duplicate().SetIsNotDup(params.Get("hash"), imageList)
	http.Redirect(w, req, "/dups", http.StatusFound)
}
