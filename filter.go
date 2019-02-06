package fitchner

import (
	"golang.org/x/net/html"
)

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

func Filter(nodes []*html.Node, f func(n *html.Node)) {
	if f == nil {
		return
	}

	for _, n := range nodes {
		f(n)
	}
}
