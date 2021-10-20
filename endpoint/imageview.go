package endpoint

import (
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/Danice123/imagine/imagedir"
	"github.com/Danice123/imagine/imageinstance"
	"github.com/Danice123/imagine/imagetag"
	"github.com/julienschmidt/httprouter"
)

type ImageData struct {
	Url         string
	Path        string
	QueryString string
	Next        string
	Previous    string
	RandomState bool
	Image       *imageinstance.ImageInstance
	ShowTags    bool
	Tags        []imageinstance.Tag
	Mood        string
}

func (ths *Endpoints) ImageView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	data := &ImageData{
		Url:         ps.ByName("path"),
		Path:        filepath.Dir(ps.ByName("path")),
		RandomState: false,
		QueryString: req.URL.Query().Encode(),
		ShowTags:    req.URL.Query().Get("tags") != "",
	}

	isRandom := false
	if cookie, err := req.Cookie("random"); err == nil {
		isRandom = true
		data.RandomState = true
		seed, _ := strconv.ParseInt(cookie.Value, 10, 64)
		rand.Seed(seed)
	}

	targetImage, err := imageinstance.New(ps.ByName("path"), ths.Root)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	data.Image = targetImage

	tags, err := imagetag.New(ths.Root)
	if err != nil {
		panic(err.Error())
	} else {
		data.Tags = tags.ReadTags(ps.ByName("path"))
		data.Mood = tags.Mapping[ps.ByName("path")].Mood
	}

	var iterator *imagedir.ImageDirIterator
	if isRandom {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.Randomize)
	} else {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.SortByName)
	}

	filters := []func(string) bool{}

	if req.URL.Query().Get("filter") != "" {
		filters = append(filters, func(name string) bool {
			imageName := strings.ReplaceAll(strings.TrimPrefix(filepath.Join(targetImage.BaseDir(), name), ths.Root), "\\", "/")
			if ok, err := tags.HasTag(imageName, req.URL.Query().Get("filter")); err != nil {
				return false
			} else {
				return !ok
			}
		})
	}

	if cookie, err := req.Cookie("mood"); err == nil {
		if cookie.Value != "" && req.URL.Query().Get("tags") == "" {
			filters = append(filters, func(name string) bool {
				imageName := strings.ReplaceAll(strings.TrimPrefix(filepath.Join(targetImage.BaseDir(), name), ths.Root), "\\", "/")
				if image, ok := tags.Mapping[imageName]; ok {
					return image.Mood != cookie.Value
				}
				return true
			})
		}
	}

	filter := func(name string) bool {
		for _, f := range filters {
			if f(name) {
				return true
			}
		}
		return false
	}

	if targetImage.IsDir {
		path := strings.TrimSuffix(strings.TrimPrefix(ps.ByName("path"), "/"), "/")
		http.Redirect(w, req, "/browse/"+path+"/"+iterator.FindNextFile(1, filter)+"?"+req.URL.Query().Encode(), http.StatusFound)
		return
	}

	data.Next = iterator.FindNextFile(1, filter)
	data.Previous = iterator.FindNextFile(-1, filter)

	var imageTemplate = template.New("ImageView")
	if html, err := os.ReadFile(filepath.Join("templates", "imageview.html")); err != nil {
		panic(err.Error())
	} else {
		if _, err := imageTemplate.Parse(string(html)); err != nil {
			panic(err.Error())
		} else {
			imageTemplate.Execute(w, data)
		}
	}
}
