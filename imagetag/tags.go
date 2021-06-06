package imagetag

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Danice123/imagine/imageinstance"
)

type TagFileOld struct {
	Tags        []string
	TagMapping  map[string]map[string]struct{}
	ImageHashes map[string]string
}

type TagTable struct {
	Tags    []string
	Mapping map[string]*TagFile
}

type TagFile struct {
	Tags  map[string]struct{}
	MD5   string
	AHash string
}

func New(root string) (*TagTable, error) {
	if rawJson, err := os.ReadFile(filepath.Join(root, ".tags.json")); err != nil {
		return &TagTable{}, nil
	} else {
		tagTable := &TagTable{}
		if err := json.Unmarshal(rawJson, tagTable); err != nil {
			return nil, err
		} else {
			return tagTable, nil
		}
	}
}

func (ths *TagTable) writeFile(root string) error {
	if jsonData, err := json.MarshalIndent(ths, "", "\t"); err != nil {
		return err
	} else {
		if err := os.WriteFile(filepath.Join(root, ".tags.json"), jsonData, os.FileMode(int(0777))); err != nil {
			return err
		}
	}
	return nil
}

func (ths *TagTable) ReadTags(file string) []imageinstance.Tag {
	tags := []imageinstance.Tag{}
	for _, tag := range ths.Tags {
		isValid := false
		if _, ok := ths.Mapping[file]; ok {
			if _, ok := ths.Mapping[file].Tags[tag]; ok {
				isValid = true
			}
		}

		tags = append(tags, imageinstance.Tag{
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

func (ths *TagTable) WriteTag(root string, file string, tag string) error {
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

	return ths.writeFile(root)
}

func (ths *TagTable) DeleteFile(root string, file string) error {
	delete(ths.Mapping, file)
	return ths.writeFile(root)
}

func (ths *TagTable) scan(root string, f func(string, string)) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		file := strings.TrimPrefix(path, root)
		if !info.IsDir() && info.Name() != ".tags.json" {
			if _, ok := ths.Mapping[file]; !ok {
				ths.Mapping[file] = &TagFile{
					Tags: make(map[string]struct{}),
				}
			}
			f(file, path)
		}
		return nil
	})
}

func (ths *TagTable) ScanMD5(root string, all bool) error {
	err := ths.scan(root, func(file string, path string) {
		if all || ths.Mapping[file].MD5 == "" {
			if hash, err := md5Hash(path); err != nil {
				panic(err)
			} else {
				ths.Mapping[file].MD5 = hash
			}
		}
	})

	if err != nil {
		panic(err)
	}
	return ths.writeFile(root)
}

func (ths *TagTable) ScanAverage(root string, all bool) error {
	err := ths.scan(root, func(file string, path string) {
		if all || ths.Mapping[file].AHash == "" {
			if hash, err := averageHash(path); err != nil {
				panic(err)
			} else {
				ths.Mapping[file].AHash = hash
			}
		}
	})

	if err != nil {
		panic(err)
	}
	return ths.writeFile(root)
}
