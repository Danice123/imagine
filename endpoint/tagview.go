package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

type TagData struct {
	Path   string
	TagMap map[string]int
}

func (this *Endpoints) TagView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(this.Root)
	if err != nil {
		panic(err.Error())
	}

	if files, err := os.ReadDir(filepath.Join(this.Root, ps.ByName("path"))); err != nil {
		http.Error(w, "Incorrent directory specified", http.StatusBadRequest)
	} else {
		data := TagData{
			Path:   ps.ByName("path"),
			TagMap: make(map[string]int),
		}
		for _, file := range files {
			imageTags := tags.ReadTags(strings.ReplaceAll(filepath.Join(ps.ByName("path"), file.Name()), "\\", "/"))
			for _, tag := range imageTags {
				if tag.Valid {
					data.TagMap[tag.Name] = data.TagMap[tag.Name] + 1
				}
			}
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
}
