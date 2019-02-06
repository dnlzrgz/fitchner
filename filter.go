package fitchner

import (
	"golang.org/x/net/html"
)

// FilterAttr returns a new []*html.Node with all the *html.Node
// with the specified attribute. If no attribute is passed it returns
// the received []*html.Node.
func FilterAttr(nodes []*html.Node, attr string) []*html.Node {
	if attr == "" {
		return nodes
	}

	var f func(n *html.Node)
	var filtered []*html.Node

	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != attr {
				continue
			}

			filtered = append(filtered, n)
		}
	}

	Filter(nodes, f)
	return filtered
}

// FilterTag returns a new []*html.Node with all the *html.Node
// that satisfy the specified tag. If no tag is passed it returns
// the received []*html.Node.
func FilterTag(nodes []*html.Node, tag string) []*html.Node {
	if tag == "" {
		return nodes
	}

	var f func(n *html.Node)
	var filtered []*html.Node

	f = func(n *html.Node) {
		if n.Data == tag {
			filtered = append(filtered, n)
		}

	}

	Filter(nodes, f)
	return filtered
}

// Links receives an []*html.Node and extracts all the links
// found in a []string.
func Links(nodes []*html.Node) []string {
	var f func(n *html.Node)
	var links []string

	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			links = append(links, a.Val)
		}
	}

	Filter(nodes, f)
	return links
}

// Filter applies for each *html.Node the function received
// as argument. If the function is nil does nothing.
func Filter(nodes []*html.Node, f func(n *html.Node)) {
	if f == nil {
		return
	}

	for _, n := range nodes {
		f(n)
	}
}
