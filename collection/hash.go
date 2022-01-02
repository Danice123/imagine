package collection

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

func (ths *CollectionHandler) MD5Hash(input *Image) (string, error) {
	file, err := os.ReadFile(input.FullPath)
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	hasher.Write(file)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (ths *CollectionHandler) PerceptionHash(input *Image) (string, error) {
	path := input.FullPath

	switch strings.ToLower(filepath.Ext(input.RelativePath)) {
	case ".png":
	case ".jpg":
	case ".jpeg":
	case ".gif":
		fallthrough
	case ".mp4":
		fallthrough
	case ".webm":
		path = ths.ExtractFrame(input)
	default:
		return "", nil
	}

	if file, err := os.Open(path); err != nil {
		return "", err
	} else {
		if img, _, err := image.Decode(file); err != nil {
			return "", err
		} else {
			if hash, err := goimagehash.PerceptionHash(img); err != nil {
				return "", err
			} else {
				return fmt.Sprintf("%d", hash.GetHash()), nil
			}
		}
	}
}

func (ths *CollectionHandler) ExtractFrame(image *Image) string {
	tempDir := filepath.Join(ths.rootDirectory, "temp")
	if _, err := os.Stat(tempDir); err != nil {
		os.Mkdir(tempDir, os.FileMode(int(0777)))
	}

	output := filepath.Join(ths.rootDirectory, "temp", filepath.Base(image.RelativePath))
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", image.FullPath,
		"-frames", "1",
		output)

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		println(cmd.String())
		panic(err)
	}

	return output
}
