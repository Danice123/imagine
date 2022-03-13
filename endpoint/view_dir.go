package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type DirData struct {
	Folders []DirDataDir
	CWD     string
}

type DirDataDir struct {
	Name string
	Path string
}

func Dir(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	folders := []DirDataDir{}
	for _, d := range COLLECTIONHANDLER.Folders(ps.ByName("path")) {
		folders = append(folders, DirDataDir{
			Name: d,
			Path: filepath.Join(ps.ByName("path"), d),
		})
	}

	var homeTemplate = template.New("Dir")
	if html, err := os.ReadFile(filepath.Join("templates", "dir.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := homeTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			homeTemplate.Execute(w, DirData{
				Folders: folders,
				CWD:     ps.ByName("path"),
			})
		}
	}
}
