package endpoint

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type HomeData struct {
	Name    string
	Folders []string
}

func (ths *Endpoints) Home(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	folders := []string{}
	filepath.WalkDir(ths.Root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			folders = append(folders, strings.ReplaceAll(strings.TrimPrefix(path+"/", ths.Root), "\\", "/"))
		}
		return nil
	})

	var homeTemplate = template.New("Home")
	if html, err := os.ReadFile(filepath.Join("templates", "home.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := homeTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			homeTemplate.Execute(w, HomeData{
				Name:    filepath.Base(ths.Root),
				Folders: folders,
			})
		}
	}
}
