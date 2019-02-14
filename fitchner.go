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
// An error is returned if the client fails to make the request, if there is
// some problem while reading the response body or if there
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

// Filter receives a []byte and returns an []*html.Node of nodes
// with the tag, the attribute or attribute's value specified. The
// three params are optional. If while parsing is any error returns the error.
func Filter(b []byte, tag, attr, val string) ([]*html.Node, error) {
	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("while parsing: %v", err)
	}

	var pre func(n *html.Node)
	var post func(n *html.Node)
	parsed := []*html.Node{}

	pre = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if tag != "" {
				if n.Data == strings.ToLower(tag) {
					parsed = append(parsed, n)
				}

				if attr != "" {
					tmp := parsed[:0]
					for _, n := range parsed {
						for _, a := range n.Attr {
							if a.Key != attr {
								continue
							}

							tmp = append(tmp, n)
						}
					}

					parsed = tmp
					return
				}

				return
			}

			if attr != "" {
				for _, a := range n.Attr {
					if a.Key != attr {
						continue
					}

					parsed = append(parsed, n)
				}

				return
			}

			parsed = append(parsed, n)
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

// Links receives a []byte and returns all links found except
// the links with prefix "tel:". Also removes the "mailto:"
// prefix if found. Returns an error if any.
func Links(b []byte) ([]string, error) {
	nodes, err := Filter(b, "a", "href", "")
	if err != nil {
		return nil, err
	}

	var links []string
	for _, n := range nodes {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if strings.Contains(a.Val, "tel:") {
					continue
				}

				links = append(links, strings.TrimPrefix(a.Val, "mailto:"))
			}
		}
	}

	return links, nil
}
