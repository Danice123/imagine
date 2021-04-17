package imagetag

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Danice123/imagine/imageinstance"
)

type TagFile struct {
	Tags       []string
	TagMapping map[string]map[string]struct{}
}

func New(root string) (*TagFile, error) {
	if rawJson, err := os.ReadFile(filepath.Join(root, ".tags.json")); err != nil {
		return nil, err
	} else {
		tagFile := &TagFile{}
		if err := json.Unmarshal(rawJson, tagFile); err != nil {
			return nil, err
		} else {
			return tagFile, nil
		}
	}
}

func (this *TagFile) ReadTags(file string) []imageinstance.Tag {
	tags := []imageinstance.Tag{}
	for _, tag := range this.Tags {
		isValid := false
		if _, ok := this.TagMapping[file][tag]; ok {
			isValid = true
		}
		tags = append(tags, imageinstance.Tag{
			Name:  tag,
			Valid: isValid,
		})
	}
	return tags
}

func (this *TagFile) HasTag(file string, tag string) bool {
	if this.TagMapping[file] == nil {
		return false
	}
	_, ok := this.TagMapping[file][tag]
	return ok
}

func (this *TagFile) WriteTag(root string, file string, tag string) error {
	if this.TagMapping == nil {
		this.TagMapping = make(map[string]map[string]struct{})
	}
	if this.TagMapping[file] == nil {
		this.TagMapping[file] = make(map[string]struct{})
	}

	if _, ok := this.TagMapping[file][tag]; ok {
		delete(this.TagMapping[file], tag)
	} else {
		this.TagMapping[file][tag] = struct{}{}
	}
	if jsonData, err := json.Marshal(this); err != nil {
		return err
	} else {
		if err := os.WriteFile(filepath.Join(root, ".tags.json"), jsonData, os.ModeAppend); err != nil {
			return err
		}
	}
	return nil
}
