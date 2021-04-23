package app

import (
	"log"
	"net/http"
	"regexp"
)

func (a *App) HTTPServer(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving page for %s", r.URL.Path)
	uri := r.URL.Path[1:]

	// Security
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Determine if this is a direct i.imgur.com image
	directImage, err := regexp.MatchString(`.*\.(jpg|jpeg|png|gif|gifv|apng|tiff|mp4|mpeg|avi|webm)`, uri)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Determine if this is an imgur.com/a/ album
	album, err := regexp.MatchString(`a/.+`, uri)
	if err != nil {
		log.Fatal(err)
	}

	// Determine if this is an imgur.com/gallery/ gallery
	gallery, err := regexp.MatchString(`gallery/.+`, uri)
	if err != nil {
		log.Fatal(err)
	}

	if directImage {
		// This is a direct image, so use the direct image handler
		a.DirectImageHandler(w, uri)
	} else if album {
		// This is an album, so use the album handler
		a.AlbumHandler(w, uri)
	} else if gallery {
		// This is a gallery, so use the gallery handler
		a.GalleryHandler(w, uri)
	} else {
		// Future proxying features not yet implemented
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
		return
	}
}
