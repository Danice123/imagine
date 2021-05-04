package imagetag

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Danice123/imagine/imageinstance"
)

type TagFile struct {
	Tags       []string
	TagMapping map[string]map[string]struct{}
}

func New(root string) (*TagFile, error) {
	if rawJson, err := os.ReadFile(filepath.Join(root, ".tags.json")); err != nil {
		return &TagFile{}, nil
	} else {
		tagFile := &TagFile{}
		if err := json.Unmarshal(rawJson, tagFile); err != nil {
			return nil, err
		} else {
			return tagFile, nil
		}
	}
}

func (ths *TagFile) ReadTags(file string) []imageinstance.Tag {
	tags := []imageinstance.Tag{}
	for _, tag := range ths.Tags {
		isValid := false
		if _, ok := ths.TagMapping[file][tag]; ok {
			isValid = true
		}
		tags = append(tags, imageinstance.Tag{
			Name:  tag,
			Valid: isValid,
		})
	}
	return tags
}

func (ths *TagFile) HasTag(file string, tag string) (bool, error) {
	if ths.TagMapping[file] == nil {
		return tag == "None", nil
	}
	if tag == "None" && len(ths.TagMapping[file]) == 0 {
		return true, nil
	}

	if _, ok := ths.TagMapping[file][tag]; ok {
		return ok, nil
	} else if expression, err := regexp.Compile("^(?i)" + strings.ReplaceAll(regexp.QuoteMeta(tag), "\\*", ".*") + "$"); err != nil {
		return false, err
	} else {
		for tagOnFile, _ := range ths.TagMapping[file] {
			if expression.MatchString(tagOnFile) {
				return true, nil
			}
		}
		return false, nil
	}
}

func (ths *TagFile) WriteTag(root string, file string, tag string) error {
	if ths.TagMapping == nil {
		ths.TagMapping = make(map[string]map[string]struct{})
	}
	if ths.TagMapping[file] == nil {
		ths.TagMapping[file] = make(map[string]struct{})
	}

	if _, ok := ths.TagMapping[file][tag]; ok {
		delete(ths.TagMapping[file], tag)
	} else {
		ths.TagMapping[file][tag] = struct{}{}
	}
	if jsonData, err := json.MarshalIndent(ths, "", "\t"); err != nil {
		return err
	} else {
		if err := os.WriteFile(filepath.Join(root, ".tags.json"), jsonData, os.ModeAppend); err != nil {
			return err
		}
	}
	return nil
}
