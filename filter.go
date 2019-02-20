package fitchner

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Filter filters the body of an HTTP response depending on the received params.
// All three params (tag, attr or val) can be an empty string. In which case a
// slice of *html.Node with all the found html.ElementNode is returned.
// An error is returned if there is any problem while parsing.
func Filter(r io.Reader, tag, attr, val string) ([]*html.Node, error) {
	doc, err := html.Parse(r)
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
		if val == "" {
			return
		}

		if n.Type == html.ElementNode {
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

// Links extracts all the links found on the body of an HTTP response.
// It ignores the links with the prefix "tel:" and removes the "mailto:" prefix.
// An error is returned if there is any problem while parsing.
func Links(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("while parsing: %v", err)
	}

	var links []string
	var pre func(n *html.Node)

	pre = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "a" {
				for _, a := range n.Attr {
					if a.Key != "href" {
						continue
					}

					links = append(links, strings.TrimPrefix(a.Val, "mailto:"))
				}
			}
		}
	}

	forEachNode(doc, pre, nil)
	return links, nil
}
