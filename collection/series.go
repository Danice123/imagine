package collection

type SeriesManager struct {
	Series map[string][]string

	hashToSeries map[string]string
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
