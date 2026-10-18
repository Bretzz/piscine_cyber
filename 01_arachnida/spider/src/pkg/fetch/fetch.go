package fetch

import (
	"io"
	"net/http"
)

// var download_path string = "./data/"

// func init() {

// }

/* Fetches the URL and returns the body and the host */
func Raw(_url string) ([]byte, *http.Response, error) {

	// #debug
	// fmt.Printf("baseUrl = '%s'\n", baseUrl)

	// Fetch the URL
	resp, err := http.Get(_url)
	if err != nil {
		// handle error
		return nil, nil, err
	}
	defer resp.Body.Close()

	// read all body so we can close the connection (?)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil, nil, err
	}

	return body, resp, nil
}
