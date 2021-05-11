package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

type CleanData struct {
	RemovedImages []string
}

func (ths *Endpoints) Clean(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	data := CleanData{
		RemovedImages: []string{},
	}
	for path := range tags.TagMapping {
		if _, err := os.Stat(filepath.Join(ths.Root, path)); err != nil {
			data.RemovedImages = append(data.RemovedImages, path)
		}
	}

	var temp = template.New("Clean")
	if html, err := os.ReadFile(filepath.Join("templates", "clean.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := temp.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			temp.Execute(w, data)
		}
	}
}
