package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

type DupData struct {
	Duplicates [][]string
}

func (ths *Endpoints) DupsView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	hashset := make(map[string][]string)
	dupHashes := []string{}
	for image, hash := range tags.ImageHashes {
		if _, ok := hashset[hash]; !ok {
			hashset[hash] = []string{image}
		} else {
			dupHashes = append(dupHashes, hash)
			hashset[hash] = append(hashset[hash], image)
		}
	}

	data := DupData{
		Duplicates: [][]string{},
	}
	for _, hash := range dupHashes {
		data.Duplicates = append(data.Duplicates, hashset[hash])
	}

	var tagTemplate = template.New("Dups")
	if html, err := os.ReadFile(filepath.Join("templates", "dupview.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := tagTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			tagTemplate.Execute(w, data)
		}
	}
}
