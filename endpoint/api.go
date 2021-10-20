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

func (ths *Endpoints) SetBrowsingMood(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	http.SetCookie(w, &http.Cookie{
		Name:  "mood",
		Value: req.URL.Query().Get("mood"),
		Path:  "/",
	})
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

func (ths *Endpoints) SetImageMood(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	targetImage, err := imageinstance.New(ps.ByName("path"), ths.Root)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if targetImage.IsDir {
		http.Error(w, "Setting mood of folders not allowed", http.StatusBadRequest)
		return
	}

	mood := req.URL.Query().Get("mood")
	if tags, err := imagetag.New(ths.Root); err != nil {
		panic(err.Error())
	} else {
		if err := tags.SetMood(ths.Root, ps.ByName("path"), mood); err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}

func (ths *Endpoints) TrashImage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	targetImage, err := imageinstance.New(ps.ByName("path"), ths.Root)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	if targetImage.IsDir {
		http.Error(w, "Trashing folders not allowed", http.StatusBadRequest)
		return
	}

	trashDir := filepath.Join(ths.Root, "trash")
	if _, err := os.Stat(trashDir); err != nil {
		os.Mkdir(trashDir, os.FileMode(int(0777)))
	}

	if err := os.Rename(targetImage.FullPath, filepath.Join(trashDir, filepath.Base(targetImage.FullPath))); err != nil {
		panic(err)
	}

	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	}

	if err := tags.DeleteFile(ths.Root, ps.ByName("path")); err != nil {
		panic(err.Error())
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
			if err := tags.DeleteFile(ths.Root, path); err != nil {
				panic(err.Error())
			}
		}
	}

	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
