package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (a *App) DirectImageHandler(w http.ResponseWriter, uri string) {
	// Attempt to get image information from Cache or Remote
	contents, contentLength, contentType, err := getImageFromCacheOrRemote(uri)
	if err != nil {
		re, _ := err.(*RequestError)
		http.Error(w, re.Error(), re.StatusCode)
		return
	}
	// Successfully got image, so return the proper response as a direct image
	w.Header().Set("Content-Length", contentLength)
	w.Header().Set("Content-Type", contentType)

	_, err = w.Write(contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getImageFromCacheOrRemote(uri string) ([]byte, string, string, error) {
	// Get the image directly from i.imgur.com
	resp, err := http.Get(fmt.Sprintf("https://i.imgur.com/%v", uri))

	if err != nil {
		log.Print(err)
		return nil, "", "", &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("%v: Internal Server Error", http.StatusInternalServerError),
		}
	}
	defer resp.Body.Close()

	// If we were unable to get this image for any reason, then respond as such
	if resp.StatusCode != 200 {
		log.Printf("Error %v looking up %v\n", resp.StatusCode, uri)
		return nil, "", "", &RequestError{
			StatusCode: resp.StatusCode,
			Err:        errors.New(uri),
		}
	}

	contentLength := fmt.Sprint(resp.ContentLength)
	contentType := resp.Header.Get("Content-Type")
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil, "", "", &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("%v: Internal Server Error", http.StatusInternalServerError),
		}
	}

	return contents, contentLength, contentType, nil
}
