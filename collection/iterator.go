package collection

import "strings"

type Filter interface {
	IsValid(*Image) bool
}

type CollectionIterator struct {
	Files       []*Image
	Filters     []Filter
	currentFile int
}

func (ths *CollectionIterator) FindNextFile(direction int) *Image {
	n := ths.currentFile + direction
	hasLooped := false

	filter := func(image *Image) bool {
		shouldFilter := false
		for _, f := range ths.Filters {
			if !f.IsValid(image) {
				shouldFilter = true
			}
		}
		return shouldFilter
	}

	for {
		if n >= len(ths.Files) {
			n = 0
			if hasLooped {
				return nil
			}
			hasLooped = true
		}
		if n < 0 {
			n = len(ths.Files) - 1
			if hasLooped {
				return nil
			}
			hasLooped = true
		}
		if ths.Files[n].IsDir() || strings.HasPrefix(ths.Files[n].Name(), ".") || filter(ths.Files[n]) {
			n += direction
			continue
		}
		break
	}
	return ths.Files[n]
}
