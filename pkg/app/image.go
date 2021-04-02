package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (a App) DirectImageHandler(w http.ResponseWriter, uri string) {
	// Get the image directly from i.imgur.com
	resp, err := http.Get(fmt.Sprintf("https://i.imgur.com/%v", uri))
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// If we were unable to get this image for any reason, then respond as such
	if resp.StatusCode != 200 {
		output := fmt.Sprintf("Error %v looking up %v\n", resp.StatusCode, uri)
		log.Print(output)
		http.Error(w, output, resp.StatusCode)
		return
	}

	// Successfully got image, so return the proper response as a direct image
	w.Header().Set("Content-Length", fmt.Sprint(resp.ContentLength))
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	if _, err = io.Copy(w, resp.Body); err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
