package imageinstance

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type ImageInstance struct {
	FullPath string
	IsDir    bool
	Mimetype string
	IsVideo  bool
	Tags     []Tag
}

type Tag struct {
	Name  string
	Valid bool
}

func New(relativePath string, rootPath string) (*ImageInstance, error) {
	new := &ImageInstance{
		FullPath: filepath.Join(rootPath, relativePath),
		IsDir:    false,
	}

	if fileInfo, err := os.Stat(new.FullPath); err != nil {
		return nil, err
	} else if fileInfo.IsDir() {
		new.IsDir = true
		return new, nil
	}

	mime, _ := mimetype.DetectFile(new.FullPath)
	new.Mimetype = mime.String()
	new.IsVideo = strings.Split(mime.String(), "/")[0] == "video"

	return new, nil
}

func (this *ImageInstance) BaseDir() string {
	if this.IsDir {
		return this.FullPath
	}
	return filepath.Dir(this.FullPath)
}
