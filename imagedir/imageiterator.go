package imagedir

import (
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ImageDirIterator struct {
	Files       []fs.DirEntry
	currentFile int
}

func New(imageDir string, currentImagePath string, sorter func([]fs.DirEntry)) *ImageDirIterator {
	new := &ImageDirIterator{}

	if files, err := os.ReadDir(imageDir); err != nil {
		panic(err.Error())
	} else {
		new.Files = files
	}

	sorter(new.Files)

	for i, file := range new.Files {
		if filepath.Join(imageDir, file.Name()) == filepath.Join(imageDir, filepath.Base(currentImagePath)) {
			new.currentFile = i
			break
		}
	}

	return new
}

func (this *ImageDirIterator) FindNextFile(direction int, filter func(string) bool) string {
	n := this.currentFile + direction
	hasLooped := false
	for {
		if n >= len(this.Files) {
			n = 0
			if hasLooped {
				return ""
			}
			hasLooped = true
		}
		if n < 0 {
			n = len(this.Files) - 1
			if hasLooped {
				return ""
			}
			hasLooped = true
		}
		if this.Files[n].IsDir() || this.Files[n].Name() == ".tags.json" || filter(this.Files[n].Name()) {
			n += direction
			continue
		}
		break
	}
	return this.Files[n].Name()
}

func SortByName(slice []fs.DirEntry) {
	sort.Slice(slice, func(i, j int) bool {
		return strings.Compare(slice[i].Name(), slice[j].Name()) < 0
	})
}

func Randomize(slice []fs.DirEntry) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
