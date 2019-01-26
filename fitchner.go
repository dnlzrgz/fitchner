// Package fitchner provides utilities for work with HTTP requests.
package fitchner

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

// Fetch receives an *http.Request and returns an *http.Response
// or an error if the request gets a bad Status Code or the
// *http.Request received is not valid.
func Fetch(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resp.Body.Close()
		return nil, fmt.Errorf("bad request: %v", err)
	}

	if checkBadStatus(resp.StatusCode) {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status code %v at %v", resp.StatusCode, req.URL)
	}

	// It's responsability of the caller to close the body of the resp.
	return resp, nil
}

// HTMLParser parses an *http.Response and returns an slice of html.Token
// with all the StartTagToken (<a> by example) found
// and closes the body of the *http.Response.
func HTMLParser(r *http.Response) []html.Token {
	b := html.NewTokenizer(r.Body)
	defer r.Body.Close()

	parsed := []html.Token{}

	for {
		tt := b.Next()

		switch {
		case tt == html.ErrorToken:
			return parsed
		case tt == html.StartTagToken:
			t := b.Token()
			parsed = append(parsed, t)
		}
	}
}

func checkBadStatus(s int) bool {
	return s != http.StatusOK
}
