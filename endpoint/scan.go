package endpoint

import (
	"github.com/Danice123/imagine/collection"
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

func Scan(conn *websocket.Conn) {
	req := &Scanrequest{}
	if err := websocket.JSON.Receive(conn, req); err != nil {
		panic(err)
	}

	hc := COLLECTIONHANDLER.HashCache()

	var hashFunc func(*collection.Image) (string, error)
	var checkHash func(*collection.Image) bool
	var writeHash func(*collection.Image, string)
	var finish func()
	switch req.ScanType {
	case "md5":
		hashFunc = COLLECTIONHANDLER.MD5Hash
		checkHash = func(image *collection.Image) bool {
			return hc.Hash(image) == ""
		}
		writeHash = hc.PutHash
		finish = func() {
			hc.Save()
		}
	case "phash":
		hd := COLLECTIONHANDLER.HashDirectory()
		hashFunc = COLLECTIONHANDLER.PerceptionHash
		checkHash = func(image *collection.Image) bool {
			d := hd.Data(hc.Hash(image))
			if d == nil {
				return true
			}
			return d.PHash == ""
		}
		writeHash = func(image *collection.Image, hash string) {
			if hd.Data(hc.Hash(image)) == nil {
				hd.CreateData(hc.Hash(image))
			}
			d := hd.Data(hc.Hash(image))
			d.PHash = hash
		}
		finish = func() {
			hd.Save()
		}
	default:
		conn.Close()
		return
	}

	images := COLLECTIONHANDLER.Scan()
	for i := 0; i < len(images); i++ {
		if req.ScanAll || checkHash(images[i]) {
			if hash, err := hashFunc(images[i]); err != nil {
				panic(err)
			} else {
				writeHash(images[i], hash)
			}
		}
		websocket.JSON.Send(conn, &Scanprogress{
			Progress: i + 1,
			Total:    len(images),
		})
	}
	finish()
	websocket.JSON.Send(conn, &Scanprogress{
		Progress: len(images),
		Total:    len(images),
	})
	conn.Close()
}
