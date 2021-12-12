package endpoint

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"text/template"

	"github.com/Danice123/imagine/imagetag"
	"github.com/corona10/goimagehash"
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

func (ths *Endpoints) DupsView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	md5set := make(map[string][]string)
	phashlist := []*goimagehash.ImageHash{}
	phashset := make(map[uint64][]string)
	for image, imageData := range tags.Mapping {
		if imageData.MD5 != "" {
			if _, ok := md5set[imageData.MD5]; !ok {
				md5set[imageData.MD5] = []string{image}
			} else {
				md5set[imageData.MD5] = append(md5set[imageData.MD5], image)
			}
		}
		if imageData.PHash != "" {
			hashbits, err := strconv.ParseUint(imageData.PHash, 10, 64)
			if err != nil {
				panic(err)
			}
			hash := goimagehash.NewImageHash(hashbits, goimagehash.PHash)
			phashlist = append(phashlist, hash)

			foundMatch := false
			for _, hashsearch := range phashlist {
				distance, err := hash.Distance(hashsearch)
				if err != nil {
					panic(err)
				}
				if distance <= 1 {
					phashset[hashsearch.GetHash()] = append(phashset[hash.GetHash()], image)
					foundMatch = true
				}
			}

			if !foundMatch {
				phashset[hash.GetHash()] = []string{image}
			}
		}
	}

	data := DupData{
		Duplicates: []*Dup{},
	}
	for hash, images := range md5set {
		if len(images) > 1 {
			sort.Strings(images)
			if checked, ok := tags.HashDups[hash]; ok {
				sort.Strings(checked)
				if reflect.DeepEqual(checked, images) {
					continue
				}
			}

			data.Duplicates = append(data.Duplicates, &Dup{
				Type:   "MD5",
				Hash:   hash,
				Images: images,
			})
		}
	}
	for hash, images := range phashset {
		if len(images) > 1 {
			sort.Strings(images)
			if checked, ok := tags.HashDups[fmt.Sprintf("%d", hash)]; ok {
				sort.Strings(checked)
				if reflect.DeepEqual(checked, images) {
					continue
				}
			}
			data.Duplicates = append(data.Duplicates, &Dup{
				Type:   "PerceptionHash",
				Hash:   fmt.Sprintf("%d", hash),
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
