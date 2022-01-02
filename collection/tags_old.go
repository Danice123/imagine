package collection

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TagTable struct {
	Tags     []string
	Mapping  map[string]*TagFile
	HashDups map[string][]string

	path string
}

type TagFile struct {
	Tags  map[string]struct{}
	MD5   string
	PHash string
}

func (ths *TagTable) WriteFile() error {
	if jsonData, err := json.MarshalIndent(ths, "", "\t"); err != nil {
		return err
	} else {
		if err := os.WriteFile(ths.path, jsonData, os.FileMode(int(0777))); err != nil {
			return err
		}
	}
	return nil
}

func (ths *TagTable) ReadTags(image *Image) []Tag {
	tags := []Tag{}
	for _, tag := range ths.Tags {
		isValid := false
		if _, ok := ths.Mapping[image.RelativePath]; ok {
			if _, ok := ths.Mapping[image.RelativePath].Tags[tag]; ok {
				isValid = true
			}
		}

		tags = append(tags, Tag{
			Name:  tag,
			Valid: isValid,
		})
	}
	return tags
}

func (ths *TagTable) HasTag(file string, tag string) (bool, error) {
	// Doesn't have a entry for the file yet
	if _, ok := ths.Mapping[file]; !ok {
		return tag == "None", nil
	}
	// Doesn't have a Tag map or an empty Tag map
	if ths.Mapping[file].Tags == nil || len(ths.Mapping[file].Tags) == 0 {
		return tag == "None", nil
	}

	if _, ok := ths.Mapping[file].Tags[tag]; ok {
		return ok, nil
	} else if expression, err := regexp.Compile("^(?i)" + strings.ReplaceAll(regexp.QuoteMeta(tag), "\\*", ".*") + "$"); err != nil {
		return false, err
	} else {
		for tagOnFile := range ths.Mapping[file].Tags {
			if expression.MatchString(tagOnFile) {
				return true, nil
			}
		}
		return false, nil
	}
}

func (ths *TagTable) WriteTag(file string, tag string) error {
	if ths.Mapping == nil {
		ths.Mapping = make(map[string]*TagFile)
	}

	if _, ok := ths.Mapping[file]; !ok {
		ths.Mapping[file] = &TagFile{
			Tags: make(map[string]struct{}),
		}
	}

	if _, ok := ths.Mapping[file].Tags[tag]; ok {
		delete(ths.Mapping[file].Tags, tag)
	} else {
		ths.Mapping[file].Tags[tag] = struct{}{}
	}

	return ths.WriteFile()
}

func (ths *TagTable) DeleteFile(file string) error {
	delete(ths.Mapping, file)
	return ths.WriteFile()
}

func (ths *TagTable) SetDupList(hash string, images []string) error {
	if ths.HashDups == nil {
		ths.HashDups = make(map[string][]string)
	}

	ths.HashDups[hash] = images
	return ths.WriteFile()
}

type FilePath struct {
	File string
	Path string
}

func (ths *TagTable) Scan(root string) []FilePath {
	files := []FilePath{}
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		file := strings.TrimPrefix(path, root)
		if !info.IsDir() && info.Name() != ".tags.json" {
			if _, ok := ths.Mapping[file]; !ok {
				ths.Mapping[file] = &TagFile{
					Tags: make(map[string]struct{}),
				}
			}
			files = append(files, FilePath{
				File: file,
				Path: path,
			})
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
