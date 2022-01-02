package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

type DupData struct {
	Duplicates []*Dup
}

type Dup struct {
	Type   string
	Hash   string
	Images []string
}

func DupsView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	hc := COLLECTIONHANDLER.HashCache()
	hd := COLLECTIONHANDLER.HashDirectory()

	md5set := hc.GetDups()
	phashset := hd.GetPHashDups()

	data := DupData{
		Duplicates: []*Dup{},
	}
	for hash, images := range md5set {
		sort.Strings(images)
		data.Duplicates = append(data.Duplicates, &Dup{
			Type:   "MD5",
			Hash:   hash,
			Images: images,
		})
	}
	for hash, images := range phashset {
		sort.Strings(images)
		// if checked, ok := tags.HashDups[hash]; ok {
		// 	sort.Strings(checked)
		// 	if reflect.DeepEqual(checked, images) {
		// 		continue
		// 	}
		// }
		data.Duplicates = append(data.Duplicates, &Dup{
			Type:   "PerceptionHash",
			Hash:   hash,
			Images: images,
		})
	}

	sort.Slice(data.Duplicates, func(i, j int) bool {
		return data.Duplicates[i].Images[0] < data.Duplicates[j].Images[0]
	})

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
