package main

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

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
	router.GET("/api/random", endpoints.ToggleRandom)
	router.GET("/api/tag/*path", endpoints.ToggleTag)
	if err := http.ListenAndServe(os.Args[1], router); err != nil {
		panic(err.Error())
	}
}
