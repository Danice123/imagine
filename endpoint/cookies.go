package endpoint

import (
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func ToggleRandom(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

func ToggleEditing(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if _, err := req.Cookie("editing"); err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:   "editing",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:  "editing",
			Value: "true",
			Path:  "/",
		})
	}
	http.Redirect(w, req, req.Referer(), http.StatusFound)
}
