package collection

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type CollectionHandler struct {
	rootDirectory string

	hc *HashCache
	hd *HashDirectory
}

func (ths *CollectionHandler) Initialize(root string) {
	ths.rootDirectory = root
}

func (ths *CollectionHandler) Folders() []string {
	folders := []string{}
	filepath.WalkDir(ths.rootDirectory, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			folders = append(folders, strings.TrimPrefix(path+"/", ths.rootDirectory))
		}
		return nil
	})
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

func (ths *CollectionHandler) Tags() *TagHandler {
	tagHandler := &TagHandler{}
	if rawJson, err := os.ReadFile(filepath.Join(ths.rootDirectory, ".tags.json")); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(rawJson, tagHandler); err != nil {
			panic(err)
		}
	}
	tagHandler.hc = ths.HashCache()
	tagHandler.hd = ths.HashDirectory()
	return tagHandler
}

func (ths *CollectionHandler) Series() *SeriesManager {
	sm := &SeriesManager{}
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
