package endpoint

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

type DupData struct {
	Duplicates []*Dup
}

type Dup struct {
	Type   string
	Images []string
}

func (ths *Endpoints) DupsView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	md5set := make(map[string][]string)
	dhashset := make(map[string][]string)
	phashset := make(map[string][]string)
	for image, imageData := range tags.Mapping {
		if imageData.MD5 != "" {
			if _, ok := md5set[imageData.MD5]; !ok {
				md5set[imageData.MD5] = []string{image}
			} else {
				md5set[imageData.MD5] = append(md5set[imageData.MD5], image)
			}
		}
		if imageData.DHash != "" {
			if _, ok := dhashset[imageData.DHash]; !ok {
				dhashset[imageData.DHash] = []string{image}
			} else {
				dhashset[imageData.DHash] = append(dhashset[imageData.DHash], image)
			}
		}
		if imageData.PHash != "" {
			if _, ok := phashset[imageData.PHash]; !ok {
				phashset[imageData.PHash] = []string{image}
			} else {
				phashset[imageData.PHash] = append(phashset[imageData.PHash], image)
			}
		}
	}

	data := DupData{
		Duplicates: []*Dup{},
	}
	for _, images := range md5set {
		if len(images) > 1 {
			sort.Strings(images)
			data.Duplicates = append(data.Duplicates, &Dup{
				Type:   "MD5",
				Images: images,
			})
		}
	}
	for _, images := range dhashset {
		if len(images) > 1 {
			sort.Strings(images)
			data.Duplicates = append(data.Duplicates, &Dup{
				Type:   "DifferenceHash",
				Images: images,
			})
		}
	}
	for _, images := range phashset {
		if len(images) > 1 {
			sort.Strings(images)
			data.Duplicates = append(data.Duplicates, &Dup{
				Type:   "PerceptionHash",
				Images: images,
			})
		}
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
