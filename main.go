package main

import (
	"net/http"
	"os"

	"github.com/Danice123/imagine/collection"
	"github.com/Danice123/imagine/endpoint"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
)

func main() {
	endpoint.COLLECTIONHANDLER = &collection.CollectionHandler{}
	endpoint.COLLECTIONHANDLER.Initialize(os.Args[2])

	router := httprouter.New()

	router.GET("/", endpoint.Home)
	router.GET("/raw/*path", endpoint.RawImage)
	router.GET("/browse/*path", endpoint.ImageView)
	router.GET("/tags/*path", endpoint.TagView)
	router.GET("/dups", endpoint.DupsView)
	router.GET("/dupcompare", endpoint.DupCompare)

	router.GET("/api/random", endpoint.ToggleRandom)
	router.GET("/api/editing", endpoint.ToggleEditing)
	router.GET("/api/tag/*path", endpoint.ToggleTag)
	router.GET("/api/changeseries", endpoint.ChangeSeries)
	// router.GET("/api/clean", endpoint.CleanImages)
	router.GET("/api/trash/*path", endpoint.TrashImage)
	router.GET("/api/scan", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		websocket.Handler(endpoint.Scan).ServeHTTP(rw, r)
	})
	router.GET("/api/markasnotdup", endpoint.MarkAsNotDup)

	static := http.FileServer(http.Dir("./templates/static"))
	router.Handler("GET", "/static/*path", http.StripPrefix("/static/", static))

	if err := http.ListenAndServe(os.Args[1], router); err != nil {
		panic(err.Error())
	}
}
