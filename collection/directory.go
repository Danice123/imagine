package collection

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Directory struct {
	FullPath     string
	RelativePath string

	collection *CollectionHandler
}

func (ths *Directory) Contents() []*Image {
	contents := []*Image{}
	if files, err := os.ReadDir(ths.FullPath); err != nil {
		panic(err)
	} else {
		for _, file := range files {
			if !file.IsDir() {
				contents = append(contents, ths.collection.Image(filepath.Join(ths.RelativePath, file.Name())))
			}
		}
	}
	return contents
}

func (ths *Directory) TagListing() map[string]int {
	tagMap := map[string]int{}
	hashDir := ths.collection.HashDirectory()
	hashCache := ths.collection.HashCache()
	for _, image := range ths.Contents() {
		md5 := hashCache.Hash(image)
		if data := hashDir.Data(md5); data != nil {
			for tag := range data.Tags {
				tagMap[tag] = tagMap[tag] + 1
			}
		}
	}
	return tagMap
}

func (ths *Directory) Iterator(currentImage *Image, sorter func([]*Image)) *CollectionIterator {
	iterator := &CollectionIterator{
		Files:   ths.Contents(),
		Filters: []Filter{},
	}
	sorter(iterator.Files)
	for i, file := range iterator.Files {
		if file.RelativePath == currentImage.RelativePath {
			iterator.currentFile = i
			break
		}
	}
	return iterator
}

func SortByName(slice []*Image) {
	sort.Slice(slice, func(i, j int) bool {
		return strings.Compare(slice[i].Name(), slice[j].Name()) < 0
	})
}

func Randomize(slice []*Image) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
