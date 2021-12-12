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

type Filter interface {
	IsValid(*imagetag.TagTable, string) bool
}

type TagFilter struct {
	Tag string
}

func (ths *TagFilter) IsValid(tags *imagetag.TagTable, name string) bool {
	if ok, err := tags.HasTag(name, ths.Tag); err != nil {
		return true
	} else {
		return ok
	}
}

type MoodFilter struct {
	Mood string
}

func (ths *MoodFilter) IsValid(tags *imagetag.TagTable, name string) bool {
	if image, ok := tags.Mapping[name]; ok {
		return image.Mood == ths.Mood
	}
	return false
}

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
	}

	var iterator *imagedir.ImageDirIterator
	if isRandom {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.Randomize)
	} else {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.SortByName)
	}

	filters := []Filter{}

	query := req.URL.Query()
	if query.Get("filter") != "" {
		for _, filter := range query["filter"] {
			filters = append(filters, &TagFilter{
				Tag: filter,
			})
		}
	}

	if cookie, err := req.Cookie("mood"); err == nil {
		if cookie.Value != "" && query.Get("tags") == "" {
			filters = append(filters, &MoodFilter{
				Mood: cookie.Value,
			})
		}
	}

	filter := func(name string) bool {
		imageName := strings.ReplaceAll(strings.TrimPrefix(filepath.Join(targetImage.BaseDir(), name), ths.Root), "\\", "/")
		shouldFilter := false
		for _, f := range filters {
			if !f.IsValid(tags, imageName) {
				shouldFilter = true
			}
		}
		return shouldFilter
	}

	if targetImage.IsDir {
		path := strings.TrimSuffix(strings.TrimPrefix(ps.ByName("path"), "/"), "/")
		http.Redirect(w, req, "/browse/"+path+"/"+iterator.FindNextFile(1, filter)+"?"+req.URL.Query().Encode(), http.StatusFound)
		return
	}

	data.Tags = tags.ReadTags(ps.ByName("path"))
	data.Mood = tags.ReadMood(ps.ByName("path"))
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
