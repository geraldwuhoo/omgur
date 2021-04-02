package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type DirectImage struct {
	Contents      []byte
	ContentLength string
	ContentType   string
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
	w.Header().Set("Content-Length", image.ContentLength)
	w.Header().Set("Content-Type", image.ContentType)

	_, err = w.Write(image.Contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getImageFromCacheOrRemote(uri string) (*DirectImage, error) {
	var directImage *DirectImage
	var err error

	// If we are using Redis, then go through the cache first
	if a.Rdb != nil {
		// Attempt to get from Redis cache
		image, err := a.Rdb.Get(context.Background(), uri).Result()

		if err == redis.Nil {
			// Get from remote if not in cache
			log.Printf("%v Redis cache miss. Pulling from remote.", uri)
			directImage, err = a.getImageFromRemote(uri)
		} else if err != nil {
			log.Print(err.Error())
			return nil, &RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(uri),
			}
		} else {
			// Unmarshal the Redis object to a DirectImage struct
			log.Printf("%v Redis cache hit. Serving from Redis.", uri)
			directImage = &DirectImage{}
			if err := json.Unmarshal([]byte(image), directImage); err != nil {
				log.Printf("Error unmarshalling from Redis: %v\n", err.Error())
				return nil, &RequestError{
					StatusCode: http.StatusInternalServerError,
					Err:        errors.New(uri),
				}
			}
		}
	} else {
		directImage, err = a.getImageFromRemote(uri)
	}

	return directImage, err
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
		Contents:      contents,
		ContentLength: fmt.Sprint(resp.ContentLength),
		ContentType:   resp.Header.Get("Content-Type"),
	}
	log.Printf("Got %v from remote\n", uri)

	// Marshall the DirectImage struct to byte array
	imageMarshall, err := json.Marshal(image)
	if err != nil {
		log.Print(err.Error())
		return nil, &RequestError{
			StatusCode: http.StatusInternalServerError,
			Err:        errors.New(uri),
		}
	}

	// Place the marshalled struct into Redis with the uri as the key, if applicable
	if a.Rdb != nil {
		if err := a.Rdb.Set(context.Background(), uri, imageMarshall, 0).Err(); err != nil {
			log.Printf("Error setting Redis cache: %v\n", err.Error())
			return nil, &RequestError{
				StatusCode: http.StatusInternalServerError,
				Err:        errors.New(uri),
			}
		}
		log.Printf("Set %v in Redis cache\n", uri)
	}

	// Construct and return this direct image
	return image, nil
}
