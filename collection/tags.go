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
	sh *SeriesManager
}

func (ths *TagHandler) concat(add *TagHandler) {
	ths.Tags = append(ths.Tags, add.Tags...)
}

func (ths *TagHandler) Get(image *Image) []Tag {
	tags := []Tag{}
	imageTags := map[string]struct{}{}
	if series, ok := ths.sh.IsImageInSeries(image); ok {
		if ths.sh.Series[series].Tags != nil {
			imageTags = ths.sh.Series[series].Tags
		}
	} else {
		if data := ths.hd.Data(ths.hc.Hash(image)); data != nil {
			imageTags = data.Tags
		}
	}

	for _, tag := range ths.Tags {
		_, isValid := imageTags[tag]
		tags = append(tags, Tag{
			Name:  tag,
			Valid: isValid,
		})
	}
	return tags
}

func (ths *TagHandler) HasTag(image *Image, tag string) (bool, error) {
	var imageTags map[string]struct{}

	if series, ok := ths.sh.IsImageInSeries(image); ok {
		if ths.sh.Series[series].Tags == nil || len(ths.sh.Series[series].Tags) == 0 {
			return tag == "None", nil
		}
		imageTags = ths.sh.Series[series].Tags
	} else if data := ths.hd.Data(ths.hc.Hash(image)); data != nil {
		if data.Tags == nil || len(data.Tags) == 0 {
			return tag == "None", nil
		}
		imageTags = data.Tags
	} else {
		return tag == "None", nil
	}

	if _, ok := imageTags[tag]; ok {
		return true, nil
	} else if expression, err := regexp.Compile("^(?i)" + strings.ReplaceAll(regexp.QuoteMeta(tag), "\\*", ".*") + "$"); err != nil {
		return false, err
	} else {
		for tagOnFile := range imageTags {
			if expression.MatchString(tagOnFile) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (ths *TagHandler) WriteTag(image *Image, tag string) {
	if series, ok := ths.sh.IsImageInSeries(image); ok {
		if ths.sh.Series[series].Tags == nil {
			ths.sh.Series[series].Tags = map[string]struct{}{}
		}

		if _, ok := ths.sh.Series[series].Tags[tag]; ok {
			delete(ths.sh.Series[series].Tags, tag)
		} else {
			ths.sh.Series[series].Tags[tag] = struct{}{}
		}
		ths.sh.Write()
		return
	}

	data := ths.hd.Data(ths.hc.Hash(image))
	if data == nil {
		data = ths.hd.CreateData(ths.hc.Hash(image))
	}

	if data.Tags == nil {
		data.Tags = map[string]struct{}{}
	}

	if _, ok := data.Tags[tag]; ok {
		delete(data.Tags, tag)
	} else {
		data.Tags[tag] = struct{}{}
	}

	ths.hd.Save()
}

func (ths *TagHandler) DeleteTag(image *Image) {
	ths.hd.DeleteData(image.MD5())
	ths.hd.Save()
}
