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
	directImage, err := regexp.MatchString(`.*\.(jpg|jpeg|png|gif|gifv|apng|tiff|mp4|mpeg|avi|webm|ogg)`, uri)
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

	// Determine if this is an imgur.com/ image
	image, err := regexp.MatchString(`\w+`, uri)
	if err != nil {
		log.Fatal(err)
	}

	if directImage {
		// This is a direct image, so use the direct image handler
		log.Print("Handling direct image")
		a.DirectImageHandler(w, uri)
	} else if album {
		// This is an album, so use the album handler
		log.Print("Handling album")
		a.AlbumHandler(w, uri)
	} else if gallery {
		// This is a gallery, so use the gallery handler
		log.Print("Handling gallery")
		a.GalleryHandler(w, uri)
	} else if image {
		// This is an image, so use the image handler
		log.Print("Handling image")
		a.ImageHandler(w, r, uri)
	} else {
		// Future proxying features not yet implemented
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
		return
	}
}
