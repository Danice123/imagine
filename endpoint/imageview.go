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
	Filter      string
	Next        string
	Previous    string
	RandomState string
	Image       *imageinstance.ImageInstance
	Tags        []imageinstance.Tag
}

func (this *Endpoints) ImageView(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	data := &ImageData{
		Url:         ps.ByName("path"),
		RandomState: "Off",
	}

	isRandom := false
	if cookie, err := req.Cookie("random"); err == nil {
		isRandom = true
		data.RandomState = "On"
		seed, _ := strconv.ParseInt(cookie.Value, 10, 64)
		rand.Seed(seed)
	}

	targetImage, err := imageinstance.New(ps.ByName("path"), this.Root)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	data.Image = targetImage

	tags, err := imagetag.New(this.Root)
	if err != nil {
		panic(err.Error())
	} else {
		data.Tags = tags.ReadTags(ps.ByName("path"))
	}

	var iterator *imagedir.ImageDirIterator
	if isRandom {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.Randomize)
	} else {
		iterator = imagedir.New(targetImage.BaseDir(), ps.ByName("path"), imagedir.SortByName)
	}

	filter := func(string) bool {
		return false
	}

	if targetImage.IsDir {
		http.Redirect(w, req, strings.TrimSuffix(req.RequestURI, "/")+"/"+iterator.FindNextFile(1, filter), http.StatusFound)
		return
	}

	data.Filter = req.URL.Query().Get("filter")
	if data.Filter != "" {
		filter = func(name string) bool {
			imageName := strings.ReplaceAll(strings.TrimPrefix(filepath.Join(targetImage.BaseDir(), name), this.Root), "\\", "/")
			return !tags.HasTag(imageName, data.Filter)
		}
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
