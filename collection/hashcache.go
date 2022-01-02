package collection

import (
	"encoding/json"
	"os"
)

type HashCache struct {
	data *map[string]string

	path string
}

func (ths *HashCache) Load() {
	ths.data = &map[string]string{}
	if rawJson, err := os.ReadFile(ths.path); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(rawJson, ths.data); err != nil {
			panic(err)
		}
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

func (ths *HashCache) Hash(image string) string {
	return (*ths.data)[image]
}

func (ths *HashCache) PutHash(image string, hash string) {
	(*ths.data)[image] = hash
}

func (ths *HashCache) GetDups() map[string][]string {
	check := map[string][]string{}
	for path, hash := range *ths.data {
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
