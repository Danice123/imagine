package collection

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/rekognition"
)

type CollectionHandler struct {
	rootDirectory string

	hc *HashCache
	hd *HashDirectory

	Rekog *rekognition.Client
}

func (ths *CollectionHandler) Initialize(root string) {
	ths.rootDirectory = root
}

func (ths *CollectionHandler) Folders(path string) []string {
	folders := []string{}
	dir, _ := ioutil.ReadDir(filepath.Join(ths.rootDirectory, path))
	for _, d := range dir {
		if d.IsDir() && d.Name() != "trash" && d.Name() != "temp" {
			folders = append(folders, d.Name())
		}
	}
	return folders
}

func (ths *CollectionHandler) Directory(path string) *Directory {
	if !ths.Image(path).IsDir() {
		return nil
	}
	return &Directory{
		FullPath:     filepath.Join(ths.rootDirectory, path),
		RelativePath: path,
		collection:   ths,
	}
}

func (ths *CollectionHandler) Image(path string) *Image {
	image := &Image{
		FullPath:     filepath.Join(ths.rootDirectory, path),
		RelativePath: path,
		collection:   ths,
	}
	return image
}

func (ths *CollectionHandler) HashCache() *HashCache {
	if ths.hc == nil {
		ths.hc = &HashCache{path: filepath.Join(ths.rootDirectory, ".hashcache.json"), hashFunc: ths.MD5Hash}
		ths.hc.Load()
	}
	return ths.hc
}

func (ths *CollectionHandler) HashDirectory() *HashDirectory {
	if ths.hd == nil {
		ths.hd = &HashDirectory{path: filepath.Join(ths.rootDirectory, ".hashdir.json")}
		ths.hd.Load()
	}
	return ths.hd
}

func (ths *CollectionHandler) Tags(contextPath string) *TagHandler {
	tagHandler := &TagHandler{}
	if rawJson, err := os.ReadFile(filepath.Join(ths.rootDirectory, ".tags.json")); err != nil {
		fmt.Println("Main Tag file not found")
	} else {
		if err := json.Unmarshal(rawJson, tagHandler); err != nil {
			panic(err)
		}
	}

	for folder := filepath.Dir(contextPath); folder != "/"; folder = filepath.Dir(folder) {
		folderTagFile := filepath.Join(ths.rootDirectory, folder, ".tags.json")
		if _, err := os.Stat(folderTagFile); err == nil {
			if rawJson, err := os.ReadFile(folderTagFile); err != nil {
				fmt.Printf("%s had read error\n", folderTagFile)
			} else {
				var additionalTags TagHandler
				if err := json.Unmarshal(rawJson, &additionalTags); err != nil {
					panic(err)
				}
				tagHandler.concat(&additionalTags)
			}

		}
	}

	tagHandler.hc = ths.HashCache()
	tagHandler.hd = ths.HashDirectory()
	tagHandler.sh = ths.Series()
	return tagHandler
}

func (ths *CollectionHandler) Series() *SeriesManager {
	sm := &SeriesManager{path: filepath.Join(ths.rootDirectory, ".series.json")}
	if rawJson, err := os.ReadFile(filepath.Join(ths.rootDirectory, ".series.json")); err != nil {
		sm.Initialize()
		return sm
	} else {
		if err := json.Unmarshal(rawJson, sm); err != nil {
			panic(err)
		}
	}
	sm.Initialize()
	return sm
}

func (ths *CollectionHandler) Duplicate() *DuplicateManager {
	dm := &DuplicateManager{path: filepath.Join(ths.rootDirectory, ".duplicates.json")}
	dm.Load()
	return dm
}

func (ths *CollectionHandler) RecognitionEngine() *RecognitionEngine {
	re := &RecognitionEngine{path: filepath.Join(ths.rootDirectory, ".recognition.json"), rekog: ths.Rekog}
	if rawJson, err := os.ReadFile(filepath.Join(ths.rootDirectory, ".recognition.json")); err != nil {
		re.Initialize()
		return re
	} else {
		if err := json.Unmarshal(rawJson, re); err != nil {
			panic(err)
		}
	}
	re.Initialize()
	return re
}

func (ths *CollectionHandler) Trash() string {
	trashDir := filepath.Join(ths.rootDirectory, "trash")
	if _, err := os.Stat(trashDir); err != nil {
		os.Mkdir(trashDir, os.FileMode(int(0777)))
	}
	return trashDir
}

func (ths *CollectionHandler) Scan() []*Image {
	images := []*Image{}
	err := filepath.Walk(ths.rootDirectory, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && !strings.HasPrefix(info.Name(), ".") {
			relativePath := strings.TrimPrefix(path, ths.rootDirectory)
			if filepath.Dir(relativePath) == "/temp" || filepath.Dir(relativePath) == "/trash" {
				return nil
			}
			images = append(images, ths.Image(relativePath))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return images
}
