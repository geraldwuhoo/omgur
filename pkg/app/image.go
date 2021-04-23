package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func (a *App) ImageHandler(w http.ResponseWriter, r *http.Request, uri string) {
	endpoint := fmt.Sprintf("https://api.imgur.com/3/image/%v", uri)
	// Build GET request to Imgur API
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", a.Authorization)

	// Execute GET request to get image details
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// If we were unable to get this image for any reason, then respond as such
	if resp.StatusCode != 200 {
		output := fmt.Sprintf("Error %v looking up %v\n", resp.StatusCode, uri)
		log.Print(output)
		http.Error(w, output, resp.StatusCode)
		return
	}

	// Get contents of the API request
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	image, err := a.ParseImage(contents)
	if err != nil {
		log.Printf("error looking up %v\n", uri)
	}
	log.Printf("Redirecting %v to %v\n", uri, image)
	http.Redirect(w, r, image, http.StatusSeeOther)
}

func (a *App) ParseImage(contents []byte) (string, error) {
	// Unpack to JSON response into an unstructured Go struct
	var result map[string]interface{}
	json.Unmarshal([]byte(contents), &result)

	// Get parent data in JSON
	data := result["data"].(map[string]interface{})
	// Get image link (safely)
	var link string
	if data["link"] == nil {
		return "", errors.New("error")
	}
	link = data["link"].(string)
	log.Printf("link: %v\n", link)

	// Parse the uri
	link = link[strings.LastIndex(link, "/"):]
	return link, nil
}
