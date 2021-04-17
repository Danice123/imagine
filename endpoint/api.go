package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Danice123/imagine/imageinstance"
	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

func (this *Endpoints) ToggleRandom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

func (this *Endpoints) ToggleTag(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	targetImage, err := imageinstance.New(ps.ByName("path"), this.Root)
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

	if tags, err := imagetag.New(this.Root); err != nil {
		panic(err.Error())
	} else {
		if err := tags.WriteTag(this.Root, ps.ByName("path"), tagName); err != nil {
			panic(err.Error())
		}
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
