package collection

import (
	"encoding/json"
	"fmt"
	"os"
)

type HashCache struct {
	data     map[string]string
	reversed map[string]string

	path     string
	hashFunc func(*Image) (string, error)
}

func (ths *HashCache) Load() {
	ths.data = map[string]string{}
	if rawJson, err := os.ReadFile(ths.path); err != nil {
		fmt.Println("Hashcache not found.")
	} else if err := json.Unmarshal(rawJson, &ths.data); err != nil {
		panic(err)
	}

	ths.reversed = map[string]string{}
	for p, h := range ths.data {
		ths.reversed[h] = p
	}
}

func (ths *HashCache) Save() {
	if jsonData, err := json.MarshalIndent(ths.data, "", "\t"); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(ths.path, jsonData, os.FileMode(int(0777))); err != nil {
			panic(err)
		}
	}
}

func (ths *HashCache) Hash(image *Image) string {
	if ths.data[image.RelativePath] == "" {
		newHash, err := ths.hashFunc(image)
		if err != nil {
			panic(err)
		}
		ths.PutHash(image, newHash)
	}

	return ths.data[image.RelativePath]
}

func (ths *HashCache) PutHash(image *Image, hash string) {
	ths.data[image.RelativePath] = hash
	ths.reversed[hash] = image.RelativePath
	ths.Save()
}

func (ths *HashCache) RemoveHash(image *Image) {
	hash := ths.data[image.RelativePath]
	delete(ths.data, image.RelativePath)
	delete(ths.reversed, hash)
	ths.Save()
}

func (ths *HashCache) GetDups() map[string][]string {
	check := map[string][]string{}
	for path, hash := range ths.data {
		if _, ok := check[hash]; !ok {
			check[hash] = []string{path}
		} else {
			check[hash] = append(check[hash], path)
		}
	}

	for key, list := range check {
		if len(list) == 1 {
			delete(check, key)
		}
	}

	return check
}

func (ths *HashCache) GetImagePathByHash(hash string) string {
	return ths.reversed[hash]
}

func (ths *HashCache) CachedImages() []string {
	list := []string{}
	for path := range ths.data {
		list = append(list, path)
	}
	return list
}
