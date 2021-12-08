package main

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"

	"github.com/Danice123/imagine/endpoint"
)

func main() {
	router := httprouter.New()
	endpoints := endpoint.Endpoints{
		Root: os.Args[2],
	}

	router.GET("/", endpoints.Home)
	router.GET("/raw/*path", endpoints.RawImage)
	router.GET("/browse/*path", endpoints.ImageView)
	router.GET("/tags/*path", endpoints.TagView)
	router.GET("/dups", endpoints.DupsView)
	router.GET("/dupcompare", endpoints.DupCompare)

	router.GET("/api/random", endpoints.ToggleRandom)
	router.GET("/api/browsingmood", endpoints.SetBrowsingMood)
	router.GET("/api/tag/*path", endpoints.ToggleTag)
	router.GET("/api/mood/*path", endpoints.SetImageMood)
	router.GET("/api/clean", endpoints.CleanImages)
	router.GET("/api/trash/*path", endpoints.TrashImage)
	router.GET("/api/scan", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		websocket.Handler(endpoints.Scan).ServeHTTP(rw, r)
	})
	router.GET("/api/markasnotdup", endpoints.MarkAsNotDup)

	static := http.FileServer(http.Dir("./templates/static"))
	router.Handler("GET", "/static/*path", http.StripPrefix("/static/", static))

	if err := http.ListenAndServe(os.Args[1], router); err != nil {
		panic(err.Error())
	}
}
