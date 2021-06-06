package imagetag

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

func md5Hash(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(file)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func averageHash(path string) (string, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
	case ".jpg":
	case ".jpeg":
	case ".gif":
		break
	default:
		return "", nil
	}

	if file, err := os.Open(path); err != nil {
		return "", err
	} else {
		if img, _, err := image.Decode(file); err != nil {
			return "", err
		} else {
			if hash, err := goimagehash.AverageHash(img); err != nil {
				return "", err
			} else {
				return fmt.Sprintf("%d", hash.GetHash()), nil
			}
		}
	}
}
