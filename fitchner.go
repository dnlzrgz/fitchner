// Package fitchner provides utilities to make HTTP requests
// and to extract information from the responses.
package fitchner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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

// SimpleFetch receives an URL and returns an error if the
// client fails to make the request or if there is a non-2xx response.
// When there is no error, returns a []byte with the body of the response.
func SimpleFetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	return Fetch(client, req)
}

// Filter receives a []byte and returns an []*html.Node of nodes
// which contain the attribute (with or without) a value. If
// while parsing is any error returns the error. If no attribute or value
// are provided simple returns all *html.Node of type html.ElementNode found.
// If value but no attribute are provided returns all the nodes with contains
// any attribute with the value specified.
func Filter(b []byte, attr, val string) ([]*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("while parsing: %v", err)
	}

	var pre func(n *html.Node)
	var post func(n *html.Node)
	parsed := []*html.Node{}

	pre = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if attr == "" {
				parsed = append(parsed, n)
				return
			}

			for _, a := range n.Attr {
				if a.Key != attr {
					continue
				}

				parsed = append(parsed, n)
			}
		}
	}

	post = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if val == "" {
				return
			}

			tmp := parsed[:0]
			for _, n := range parsed {
				for _, a := range n.Attr {
					if !strings.Contains(a.Val, val) {
						continue
					}

					tmp = append(tmp, n)
				}
			}

			parsed = tmp
		}
	}

	forEachNode(doc, pre, post)
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
