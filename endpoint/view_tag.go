package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type TagData struct {
	Path   string
	TagMap map[string]int
}

func TagView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	dir := COLLECTIONHANDLER.Directory(ps.ByName("path"))
	data := TagData{
		Path:   ps.ByName("path"),
		TagMap: dir.TagListing(),
	}

	var tagTemplate = template.New("Tags")
	if html, err := os.ReadFile(filepath.Join("templates", "tagview.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := tagTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			tagTemplate.Execute(w, data)
		}
	}
}
