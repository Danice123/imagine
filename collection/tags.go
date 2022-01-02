package collection

import (
	"regexp"
	"strings"
)

type Tag struct {
	Name  string
	Valid bool
}

type TagHandler struct {
	Tags []string

	hc *HashCache
	hd *HashDirectory
}

func (ths *TagHandler) Get(image *Image) []Tag {
	tags := []Tag{}
	for _, tag := range ths.Tags {
		isValid := false
		if data := ths.hd.Data(ths.hc.Hash(image.RelativePath)); data != nil {
			if _, ok := data.Tags[tag]; ok {
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

func (ths *TagHandler) HasTag(image *Image, tag string) (bool, error) {
	if data := ths.hd.Data(ths.hc.Hash(image.RelativePath)); data != nil {
		if _, ok := data.Tags[tag]; ok {
			return true, nil
		} else if expression, err := regexp.Compile("^(?i)" + strings.ReplaceAll(regexp.QuoteMeta(tag), "\\*", ".*") + "$"); err != nil {
			return false, err
		} else {
			for tagOnFile := range data.Tags {
				if expression.MatchString(tagOnFile) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (ths *TagHandler) WriteTag(image *Image, tag string) {
	data := ths.hd.Data(ths.hc.Hash(image.RelativePath))
	if data == nil {
		data = ths.hd.CreateData(ths.hc.Hash(image.RelativePath))
	}

	if _, ok := data.Tags[tag]; ok {
		delete(data.Tags, tag)
	} else {
		data.Tags[tag] = struct{}{}
	}

	ths.hd.Save()
}

func (ths *TagHandler) DeleteTag(image *Image) {
	ths.hd.DeleteData(image.RelativePath)
	ths.hd.Save()
}
