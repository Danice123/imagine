package collection

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

type Image struct {
	FullPath     string
	RelativePath string

	fileInfo   fs.FileInfo
	mimetype   *mimetype.MIME
	collection *CollectionHandler
}

func (ths *Image) processMime() {
	if mime, err := mimetype.DetectFile(ths.FullPath); err != nil {
		fmt.Fprintf(os.Stderr, "File (%s) does not exist or is not readable", ths.FullPath)
	} else {
		ths.mimetype = mime
	}
}

func (ths *Image) processFile() {
	if fileInfo, err := os.Stat(ths.FullPath); err != nil {
		fmt.Fprintf(os.Stderr, "File (%s) does not exist", ths.FullPath)
	} else {
		ths.fileInfo = fileInfo
	}
}

func (ths *Image) Mimetype() string {
	if ths.mimetype == nil {
		ths.processMime()
	}
	return ths.mimetype.String()
}

func (ths *Image) IsVideo() bool {
	if ths.mimetype == nil {
		ths.processMime()
	}
	return strings.Split(ths.mimetype.String(), "/")[0] == "video"
}

func (ths *Image) IsValid() bool {
	if ths.fileInfo == nil {
		ths.processFile()
	}
	return ths.fileInfo != nil
}

func (ths *Image) IsDir() bool {
	if ths.fileInfo == nil {
		ths.processFile()
	}
	return ths.fileInfo.IsDir()
}

func (ths *Image) Name() string {
	if ths.fileInfo == nil {
		ths.processFile()
	}
	return ths.fileInfo.Name()
}

func (ths *Image) Directory() *Directory {
	if ths.IsDir() {
		return &Directory{
			FullPath:     ths.FullPath,
			RelativePath: ths.RelativePath,
			collection:   ths.collection,
		}
	}
	return &Directory{
		FullPath:     filepath.Dir(ths.FullPath),
		RelativePath: filepath.Dir(ths.RelativePath),
		collection:   ths.collection,
	}
}
