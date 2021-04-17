package endpoint

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
)

func (this *Endpoints) RawImage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	http.ServeFile(w, req, filepath.Join(this.Root, ps.ByName("path")))
}
