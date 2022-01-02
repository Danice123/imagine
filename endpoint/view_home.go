package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type HomeData struct {
	Folders []string
}

func Home(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	var homeTemplate = template.New("Home")
	if html, err := os.ReadFile(filepath.Join("templates", "home.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := homeTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			homeTemplate.Execute(w, HomeData{
				Folders: COLLECTIONHANDLER.Folders(),
			})
		}
	}
}
