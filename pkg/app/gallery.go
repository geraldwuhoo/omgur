package app

import (
	"fmt"
	"net/http"
)

func (a *App) GalleryHandler(w http.ResponseWriter, uri string) {
	a.GetAlbum(w, uri, fmt.Sprintf("https://api.imgur.com/3/gallery/%v", uri[8:]))
}
