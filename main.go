package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/Danice123/imagine/collection"
	"github.com/Danice123/imagine/endpoint"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
)

func initilizeRekog() *rekognition.Client {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile("imagine"))
	if err != nil {
		panic(err)
	}
	return rekognition.NewFromConfig(config, func(o *rekognition.Options) { o.Region = "us-east-2" })
}

func main() {
	rekog := initilizeRekog()

	endpoint.COLLECTIONHANDLER = &collection.CollectionHandler{
		Rekog: rekog,
	}
	root := strings.TrimSuffix(os.Args[2], "/")
	_, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	endpoint.COLLECTIONHANDLER.Initialize(root)

	router := httprouter.New()

	router.GET("/", endpoint.Home)
	router.GET("/raw/*path", endpoint.RawImage)
	router.GET("/face/*path", endpoint.RawFace)
	router.GET("/browse/*path", endpoint.ImageView)
	router.GET("/tags/*path", endpoint.TagView)
	router.GET("/dups", endpoint.DupsView)
	router.GET("/dupcompare", endpoint.DupCompare)

	router.GET("/api/random", endpoint.ToggleRandom)
	router.GET("/api/editing", endpoint.ToggleEditing)
	router.GET("/api/tag/*path", endpoint.ToggleTag)
	router.GET("/api/changeseries", endpoint.ChangeSeries)
	router.GET("/api/clean", endpoint.CleanImages)
	router.GET("/api/trash/*path", endpoint.TrashImage)
	router.GET("/api/scan", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		websocket.Handler(endpoint.Scan).ServeHTTP(rw, r)
	})
	router.GET("/api/markasnotdup", endpoint.MarkAsNotDup)

	router.GET("/api/aws/detectface/*path", endpoint.DetectFaces)

	static := http.FileServer(http.Dir("./templates/static"))
	router.Handler("GET", "/static/*path", http.StripPrefix("/static/", static))

	if err := http.ListenAndServe(os.Args[1], router); err != nil {
		panic(err.Error())
	}
}
