// Package fitchner provides utilities to make HTTP requests
// and to extract information from the responses.
package fitchner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html"
)

// Fetch receives an http.Client and an http.Request to make a request.
// An error is returned if the client fails to make the request or if there
// is a non-2xx response. When there is no error, returns a []byte with
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

// Nodes receives a []byte and extracts all the nodes found.
// If something goes wrong at parsing it returns an error.
// When there is no error, returns an []*html.Node with all the nodes found.
func Nodes(b []byte) ([]*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("while parsing: %v", err)
	}

	var parser func(n *html.Node)
	parsed := []*html.Node{}

	parser = func(n *html.Node) {
		if n.Type == html.ElementNode {
			parsed = append(parsed, n)
		}
	}

	forEachNode(doc, parser, nil)
	return parsed, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}
