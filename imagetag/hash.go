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
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corona10/goimagehash"
)

func MD5Hash(root string, path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write(file)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func PerceptionHash(root string, path string) (string, error) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".png":
	case ".jpg":
	case ".jpeg":
	case ".gif":
		fallthrough
	case ".mp4":
		fallthrough
	case ".webm":
		path = ExtractFrame(root, path, filepath.Base(path))
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

func ExtractFrame(root string, path string, name string) string {
	tempDir := filepath.Join(root, "temp")
	if _, err := os.Stat(tempDir); err != nil {
		os.Mkdir(tempDir, os.FileMode(int(0777)))
	}

	output := filepath.Join(root, "temp", name+".png")
	cmd := exec.Command("ffmpeg",
		"-i", path,
		"-frames", "1",
		output)

	if err := cmd.Start(); err != nil {
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		panic(err.Error())
	}

	return output
}
