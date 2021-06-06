package endpoint

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Danice123/imagine/imageinstance"
	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

func (ths *Endpoints) ToggleRandom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if _, err := req.Cookie("random"); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "random",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:  "random",
			Value: fmt.Sprint(time.Now().UnixNano()),
			Path:  "/",
		})
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (ths *Endpoints) ToggleTag(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	targetImage, err := imageinstance.New(ps.ByName("path"), ths.Root)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if targetImage.IsDir {
		http.Error(w, "Tagging folders not allowed", http.StatusBadRequest)
		return
	}

	tagName := req.URL.Query().Get("name")
	if tagName == "" {
		return
	}

	if tags, err := imagetag.New(ths.Root); err != nil {
		panic(err.Error())
	} else {
		if err := tags.WriteTag(ths.Root, ps.ByName("path"), tagName); err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (ths *Endpoints) CleanImages(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	for path := range tags.Mapping {
		if _, err := os.Stat(filepath.Join(ths.Root, path)); err != nil {
			delete(tags.Mapping, path)
		}
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (ths *Endpoints) ScanImages(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	shouldScanAll := false
	overwriteArg := req.URL.Query().Get("overwrite")
	if overwriteArg != "" {
		shouldScanAll = true
	}

	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	switch ps.ByName("hash") {
	case "md5":
		err = tags.ScanMD5(ths.Root, shouldScanAll)
		if err != nil {
			panic(err)
		}
	case "ahash":
		err = tags.ScanAverage(ths.Root, shouldScanAll)
		if err != nil {
			panic(err)
		}
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
