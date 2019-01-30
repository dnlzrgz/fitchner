// Package fitchner provides utilities to make HTTP requests
// and to extract information from the responses.
package fitchner

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Fetch receives an *http.Client and an *http.Request to make a request.
// An error is returned if the client's fails to make the request or if there
// is a non-2xx response. When there is no error, returns a byte slice with
// the body of the response.
func Fetch(c *http.Client, req *http.Request) ([]byte, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if checkBadStatus(resp.StatusCode) {
		return nil, fmt.Errorf("got %v at %s", resp.Status, req.URL)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading %s: %v", req.URL, err)
	}

	return b, nil
}

func checkBadStatus(s int) bool {
	return s != http.StatusOK
}
