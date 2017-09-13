package controller

import (
	"net/http"
	"strings"

	"github.com/hyperscale/hyperpaas/server"
)

type UiController struct {
}

func NewUiController() (*UiController, error) {
	return &UiController{}, nil
}

func (c UiController) Mount(r *server.Router) {
	r.AddPrefixRouteFunc("/ui/", c.StaticFile)
}

func (c UiController) StaticFile(w http.ResponseWriter, r *http.Request) {
	filename := strings.Replace(r.URL.Path, "/ui/", "/", 1)

	extensions := []string{".js", ".css", ".map", ".ico"}
	for _, ext := range extensions {
		if strings.HasSuffix(r.URL.Path, ext) {
			http.ServeFile(w, r, "/opt/hyperpaas/ui/"+filename)

			return
		}
	}

	http.ServeFile(w, r, "/opt/hyperpaas/ui/index.html")
}
