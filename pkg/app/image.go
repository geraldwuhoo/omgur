package app

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type DirectImage struct {
	contents      []byte
	contentLength string
	contentType   string
}

func (a *App) DirectImageHandler(w http.ResponseWriter, uri string) {
	// Attempt to get image information from Cache or Remote
	image, err := a.getImageFromCacheOrRemote(uri)
	if err != nil {
		re, _ := err.(*RequestError)
		http.Error(w, re.Error(), re.StatusCode)
		return
	}
	// Successfully got image, so return the proper response as a direct image
	w.Header().Set("Content-Length", image.contentLength)
	w.Header().Set("Content-Type", image.contentType)

	_, err = w.Write(image.contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getImageFromCacheOrRemote(uri string) (*DirectImage, error) {
	rdb := a.Rdb
	_, err := rdb.Get(context.Background(), uri).Bytes()
	if err == redis.Nil {
		log.Printf("%v not found in redis cache. Pulling from remote.", uri)
		return a.getImageFromRemote(uri)
	} else if err != nil {
		log.Print(err.Error())
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(uri),
		}
	} else {
		// TODO: Fix Redis
		return nil, nil
	}

}

func (a *App) getImageFromRemote(uri string) (*DirectImage, error) {
	// Get the image directly from i.imgur.com
	resp, err := http.Get(fmt.Sprintf("https://i.imgur.com/%v", uri))

	if err != nil {
		log.Print(err)
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(uri),
		}
	}
	defer resp.Body.Close()

	// If we were unable to get this image for any reason, then respond as such
	if resp.StatusCode != 200 {
		log.Printf("Error %v looking up %v\n", resp.StatusCode, uri)
		return nil, &RequestError{
			StatusCode: resp.StatusCode,
			Err:        errors.New(uri),
		}
	}

	// Get the image contents, with error handling
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(uri),
		}
	}

	// Construct image
	image := &DirectImage{
		contents:      contents,
		contentLength: fmt.Sprint(resp.ContentLength),
		contentType:   resp.Header.Get("Content-Type"),
	}

	// TODO: place image into Redis cache

	// Construct and return this direct image
	return image, nil
}
