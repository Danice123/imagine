package endpoint

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	"github.com/Danice123/imagine/collection"
	"github.com/julienschmidt/httprouter"
)

type TagFilter struct {
	Tag        string
	TagHandler *collection.TagHandler
}

func (ths *TagFilter) IsValid(image *collection.Image) bool {
	if ok, err := ths.TagHandler.HasTag(image, ths.Tag); err != nil {
		return true
	} else {
		return ok
	}
}

type ImageData struct {
	Url         string
	Path        string
	QueryString string
	Next        string
	Previous    string
	RandomState bool
	Image       *collection.Image
	ImageSeries string
	Hash        string
	ShowTags    bool
	Tags        []collection.Tag
	Faces       []collection.FaceBox
	Series      []string
}

func handleEditingCookie(req *http.Request) bool {
	if cookie, err := req.Cookie("editing"); err == nil {
		return cookie.Value == "true"
	}
	return false
}

func handleRandomCookie(req *http.Request) bool {
	if cookie, err := req.Cookie("random"); err == nil {
		seed, _ := strconv.ParseInt(cookie.Value, 10, 64)
		rand.Seed(seed)
		return true
	}
	return false
}

func ImageView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	query := req.URL.Query()

	image := COLLECTIONHANDLER.Image(ps.ByName("path"))
	if !image.IsValid() {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	series := COLLECTIONHANDLER.Series()

	data := &ImageData{
		Url:         ps.ByName("path"),
		Path:        filepath.Dir(ps.ByName("path")),
		RandomState: handleRandomCookie(req),
		QueryString: query.Encode(),
		ShowTags:    handleEditingCookie(req),
		Image:       image,
		Series:      series.SeriesList(),
	}

	dir := image.Directory()
	var iterator *collection.CollectionIterator
	if data.RandomState {
		iterator = dir.Iterator(image, collection.Randomize)
	} else {
		iterator = dir.Iterator(image, collection.SortByName)
	}

	tagHandler := COLLECTIONHANDLER.Tags(image.RelativePath)
	if query.Get("filter") != "" {
		for _, filter := range query["filter"] {
			iterator.Filters = append(iterator.Filters, &TagFilter{
				Tag:        filter,
				TagHandler: tagHandler,
			})
		}
	}

	if image.IsDir() {
		next := iterator.FindNextFile(1)
		if next == nil {
			http.Error(w, "No image found", http.StatusNotFound)
			return
		}

		http.Redirect(w, req, fmt.Sprintf("/browse/%s?%s", next.RelativePath, req.URL.Query().Encode()), http.StatusFound)
		return
	}

	data.Tags = tagHandler.Get(image)
	data.Hash = image.MD5()
	imageSeries, _ := series.IsImageInSeries(image)
	data.ImageSeries = imageSeries

	if data.ShowTags {
		hdd := COLLECTIONHANDLER.HashDirectory().Data(data.Hash)
		if hdd != nil && hdd.Faces != nil {
			data.Faces = hdd.Faces
		}
	}

	next := iterator.FindNextFile(1)
	if next != nil {
		data.Next = next.RelativePath
	}
	previous := iterator.FindNextFile(-1)
	if previous != nil {
		data.Previous = previous.RelativePath
	}

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
