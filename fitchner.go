// Package fitchner provides utilities to make HTTP requests
// and to extract information from the responses.
package fitchner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Fetch makes an HTTP request and returns the response.
// Receives an http.Client and an http.Request for it.
// An error is returned if:
//		*	The client fails to make the request.
// 		*	There is some problem while reading the response body.
//		*	There is a non-2xxx response.
// When there is no error, returns a *bytes.Reader with the body of the response.
func Fetch(c *http.Client, req *http.Request) (*bytes.Reader, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("while making a request to %q: %v", req.URL, err)
	}
	defer resp.Body.Close()

	if checkBadStatus(resp.StatusCode) {
		return nil, fmt.Errorf("got %v at %s", resp.Status, req.URL)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading the response to %q: %v", req.URL, err)
	}

	return bytes.NewReader(b), nil
}

func checkBadStatus(s int) bool { return s != http.StatusOK }
