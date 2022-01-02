package collection

import "strings"

type Filter interface {
	IsValid(*Image) bool
}

type CollectionIterator struct {
	Images      []*Image
	Filters     []Filter
	series      *SeriesManager
	currentFile int
}

func (ths *CollectionIterator) FindNextFile(direction int) *Image {
	if _, ok := ths.series.IsImageInSeries(ths.Images[ths.currentFile]); ok {
		hash := ths.series.NextImageHashInSeries(ths.Images[ths.currentFile], direction)
		for _, image := range ths.Images {
			if image.MD5() == hash {
				return image
			}
		}
	}

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
		if n >= len(ths.Images) {
			n = 0
			if hasLooped {
				return nil
			}
			hasLooped = true
		}
		if n < 0 {
			n = len(ths.Images) - 1
			if hasLooped {
				return nil
			}
			hasLooped = true
		}
		if ths.Images[n].IsDir() || strings.HasPrefix(ths.Images[n].Name(), ".") || filter(ths.Images[n]) {
			n += direction
			continue
		}
		if series, ok := ths.series.IsImageInSeries(ths.Images[n]); ok {
			if direction > 0 && ths.series.Series[series][0] != ths.Images[n].MD5() {
				n += direction
				continue
			}
			if direction < 0 && ths.series.Series[series][len(ths.series.Series[series])-1] != ths.Images[n].MD5() {
				n += direction
				continue
			}
		}
		break
	}
	return ths.Images[n]
}
