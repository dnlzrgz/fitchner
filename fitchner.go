// Package fitchner provides utilities for work with HTTP requests.
package fitchner

import (
	"fmt"
	"net/http"
)

// Fetch receives an *http.Request and returns an *http.Response
// or an error if the request gets a bad Status Code or the
// *http.Request received is not valid.
func Fetch(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("bad request: %v", err)
	}

	defer resp.Body.Close()

	if checkBadStatus(resp.StatusCode) {
		return nil, fmt.Errorf("bad status code %v at %v", resp.StatusCode, req.URL)
	}

	return resp, nil
}

func checkBadStatus(s int) bool {
	return s != http.StatusOK
}
