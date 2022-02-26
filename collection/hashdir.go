package collection

import (
	"encoding/json"
	"fmt"
	"os"
)

type HashDirectory struct {
	data map[string]*ImageHashData

	path string
}

type ImageHashData struct {
	Tags  map[string]struct{}
	Faces map[string]FaceBox
	PHash string
}

func (ths *HashDirectory) Load() {
	ths.data = map[string]*ImageHashData{}
	if rawJson, err := os.ReadFile(ths.path); err != nil {
		fmt.Println("Hashdir not found.")
	} else {
		if err := json.Unmarshal(rawJson, &ths.data); err != nil {
			panic(err)
		}
	}
}

func (ths *HashDirectory) Save() {
	if jsonData, err := json.MarshalIndent(ths.data, "", "\t"); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(ths.path, jsonData, os.FileMode(int(0777))); err != nil {
			panic(err)
		}
	}
}

func (ths *HashDirectory) Data(hash string) *ImageHashData {
	return ths.data[hash]
}

func (ths *HashDirectory) CreateData(hash string) *ImageHashData {
	newData := &ImageHashData{}
	ths.data[hash] = newData
	return newData
}

func (ths *HashDirectory) DeleteData(hash string) {
	delete(ths.data, hash)
}

func (ths *HashDirectory) GetPHashDups() map[string][]string {
	check := map[string][]string{}
	for hash, data := range ths.data {
		if data.PHash != "" {
			if _, ok := check[data.PHash]; !ok {
				check[data.PHash] = []string{hash}
			} else {
				check[data.PHash] = append(check[data.PHash], hash)
			}
		}
	}

	for key, list := range check {
		if len(list) == 1 {
			delete(check, key)
		}
	}

	return check
}
