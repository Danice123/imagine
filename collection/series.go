package collection

import (
	"encoding/json"
	"os"
)

type SeriesManager struct {
	Series map[string][]string

	hashToSeries map[string]string
	path         string
}

func (ths *SeriesManager) Initialize() {
	ths.hashToSeries = map[string]string{}
	for s, list := range ths.Series {
		for _, hash := range list {
			ths.hashToSeries[hash] = s
		}
	}
}

func (ths *SeriesManager) IsImageInSeries(image *Image) (string, bool) {
	s, ok := ths.hashToSeries[image.MD5()]
	return s, ok
}

func (ths *SeriesManager) NextImageHashInSeries(image *Image, direction int) string {
	if s, ok := ths.hashToSeries[image.MD5()]; ok {
		i := 0
		for ; i < len(ths.Series[s]); i++ {
			if ths.Series[s][i] == image.MD5() {
				break
			}
		}
		if direction > 0 && i == len(ths.Series[s])-1 {
			return ""
		}
		if direction < 0 && i == 0 {
			return ""
		}
		return ths.Series[s][i+direction]
	}
	return ""
}

func (ths *SeriesManager) RemoveImageFromSeries(hash string, change string) {
	delete(ths.hashToSeries, hash)
	newL := []string{}
	for _, h := range ths.Series[change] {
		if h != hash {
			newL = append(newL, h)
		}
	}
	ths.Series[change] = newL
	if len(ths.Series[change]) == 0 {
		delete(ths.Series, change)
	}
}

func (ths *SeriesManager) AddImageToSeries(hash string, change string) {
	if series, ok := ths.hashToSeries[hash]; ok {
		ths.RemoveImageFromSeries(hash, series)
		if series == change {
			return
		}
	}
	ths.hashToSeries[hash] = change
	if _, ok := ths.Series[change]; !ok {
		ths.Series[change] = []string{hash}
	} else {
		ths.Series[change] = append(ths.Series[change], hash)
	}
}

func (ths *SeriesManager) Write() {
	if data, err := json.MarshalIndent(ths, "", "\t"); err != nil {
		panic(err)
	} else {
		if err := os.WriteFile(ths.path, data, os.FileMode(int(0777))); err != nil {
			panic(err)
		}
	}
}

func (ths *SeriesManager) SeriesList() []string {
	l := []string{}
	for s := range ths.Series {
		l = append(l, s)
	}
	return l
}
