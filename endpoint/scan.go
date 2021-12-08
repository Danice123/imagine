package endpoint

import (
	"path/filepath"

	"github.com/Danice123/imagine/imagetag"
	"golang.org/x/net/websocket"
)

type Scanrequest struct {
	ScanType string
	ScanAll  bool
}

type Scanprogress struct {
	Progress int
	Total    int
}

func (ths *Endpoints) Scan(conn *websocket.Conn) {
	req := &Scanrequest{}
	if err := websocket.JSON.Receive(conn, req); err != nil {
		panic(err)
	}

	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	var hashFunc func(string, string) (string, error)
	var checkHash func(string) bool
	var writeHash func(string, string)
	switch req.ScanType {
	case "md5":
		hashFunc = imagetag.MD5Hash
		checkHash = func(file string) bool {
			return tags.Mapping[file].MD5 == ""
		}
		writeHash = func(file string, hash string) {
			tags.Mapping[file].MD5 = hash
		}
	case "dhash":
		hashFunc = imagetag.DifferenceHash
		checkHash = func(file string) bool {
			return tags.Mapping[file].DHash == ""
		}
		writeHash = func(file string, hash string) {
			tags.Mapping[file].DHash = hash
		}
	case "phash":
		hashFunc = imagetag.PerceptionHash
		checkHash = func(file string) bool {
			return tags.Mapping[file].PHash == ""
		}
		writeHash = func(file string, hash string) {
			tags.Mapping[file].PHash = hash
		}
	default:
		conn.Close()
		return
	}

	files := tags.Scan(ths.Root)
	for i := 0; i < len(files); i++ {
		if filepath.Dir(files[i].Path) == filepath.Join(ths.Root, "temp") || filepath.Dir(files[i].Path) == filepath.Join(ths.Root, "trash") {
			continue
		}
		if req.ScanAll || checkHash(files[i].File) {
			if hash, err := hashFunc(ths.Root, files[i].Path); err != nil {
				panic(err)
			} else {
				writeHash(files[i].File, hash)

			}
			if i%10 == 0 {
				tags.WriteFile(ths.Root)
			}
		}
		websocket.JSON.Send(conn, &Scanprogress{
			Progress: i + 1,
			Total:    len(files),
		})
	}
	websocket.JSON.Send(conn, &Scanprogress{
		Progress: len(files),
		Total:    len(files),
	})
	tags.WriteFile(ths.Root)
	conn.Close()
}
