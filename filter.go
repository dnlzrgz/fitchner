package fitchner

import (
	"strings"

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

// FilterClass returns a new []*html.Node with all the *html.Node
// that has the specified class. If no class is passed it returns
// the received []*html.Node. If no *html.Node is found returns an empty
// []*html.Node.
func FilterClass(nodes []*html.Node, class string) []*html.Node {
	if class == "" {
		return nodes
	}

	var f func(n *html.Node)
	var filtered []*html.Node

	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != "class" {
				continue
			}

			if strings.Contains(a.Val, class) {
				filtered = append(filtered, n)
			}
		}
	}

	Filter(nodes, f)
	return filtered
}

// FilterID returns a new *html.Node if any *html.Node with the specified
// id is found. If no id is passed it returns nil. If no *html.Node is found
// also returns nil.
func FilterID(nodes []*html.Node, id string) *html.Node {
	if id == "" {
		return nil
	}

	var f func(n *html.Node)
	var filtered *html.Node

	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key != "id" {
				continue
			}

			if strings.Contains(a.Val, id) {
				filtered = n
				break
			}
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
