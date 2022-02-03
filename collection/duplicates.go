package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
)

type DuplicateManager struct {
	data map[string][]string

	path string
}

func (ths *DuplicateManager) Load() {
	ths.data = map[string][]string{}
	if rawJson, err := os.ReadFile(ths.path); err != nil {
		fmt.Println("Hashcache not found.")
	} else if err := json.Unmarshal(rawJson, &ths.data); err != nil {
		panic(err)
	}
}

func (ths *DuplicateManager) Save() {
	if jsonData, err := json.MarshalIndent(ths.data, "", "\t"); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(ths.path, jsonData, os.FileMode(int(0777))); err != nil {
			panic(err)
		}
	}
}

func (ths *DuplicateManager) IsDup(hash string, images []string) bool {
	if verified, ok := ths.data[hash]; ok {
		sort.Strings(images)
		sort.Strings(verified)
		return !reflect.DeepEqual(verified, images)
	}
	return true
}

func (ths *DuplicateManager) SetIsNotDup(hash string, images []string) {
	ths.data[hash] = images
	ths.Save()
}
