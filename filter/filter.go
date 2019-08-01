package filter

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// FilterFn defines a filter to filter HTML nodes.
// If the filter returns true, the HTML node is added
// to the result of the Filter function.
// If the filter returns false the HTML node is excluded.
type FilterFn func(n *html.Node) bool

// FilterByTag receives an HTML tag without "<" nor ">"
// and returns a FilterFn that can be used as a filter
// on the Filter function.
func FilterByTag(tag string) FilterFn {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}

		if n.Data == tag {
			return true
		}

		return false
	}
}

// FilterByClass receives a CSS class and
// returns a FilterFn that can be used as a filter
// on the Filter function.
func FilterByClass(class string) FilterFn {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}

		for _, a := range n.Attr {
			if a.Key != "class" {
				continue
			}

			if ok := strings.Contains(a.Val, class); ok {
				return true
			}
		}

		return false
	}
}

// FilterByID receives an ID and returns a FilterFn
// that can be used as a filter on the Filter function.
func FilterByID(id string) FilterFn {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}

		for _, a := range n.Attr {
			if a.Key != "id" {
				continue
			}

			if ok := strings.Contains(a.Val, id); ok {
				return true
			}
		}

		return false
	}
}

// FilterByAttr receives an HTML attribute and returns a
// FilterFn that can be used as a filter on the Filter function.
func FilterByAttr(attr string) FilterFn {
	return func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return false
		}

		for _, a := range n.Attr {
			if a.Key != attr {
				continue
			}

			return true
		}

		return false
	}
}

// Filter receives an io.Reader from which to extract the HTML
// nodes.
//
// It returns an error if there is any problem while extracting the tokens
// to create the HTML nodes.
//
// You can pass none, one or more FilterFn functions to apply.
// The order of the filters affects the result.
//
// It finally returns an []*html.Node.
func Filter(r io.Reader, filters ...FilterFn) ([]*html.Node, error) {
	nodes, err := tokens(r)
	if err != nil {
		return nil, fmt.Errorf("while tokenizing: %v", err)
	}

	if filters == nil {
		return nodes, nil
	}

	filtered := forEachNode(nodes, filters...)
	return filtered, nil
}

// Links receives an io.Reader from which to extract all the links.
// It returns an error if there is any problem while extrating the tokens
// to create the HTML nodes from which extract the links.
// If nothing goes wrong it returns a []string with all the links found.
func Links(r io.Reader) ([]string, error) {
	var links []string
	nodes, err := Filter(r, FilterByAttr("href"))
	if err != nil {
		return nil, fmt.Errorf("while extracting nodes with attribute \"href\": %v", err)
	}

	for _, n := range nodes {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}

			links = append(links, a.Val)
			break
		}
	}

	return links, nil
}

// Images receives an io.Reader from which to extract all the image sources.
// It returns an error if there is any problem while extrating the tokens
// to create the HTML nodes from which extract the images sources.
// If nothing goes wrong it returns a []string with all the image sources found.
func Images(r io.Reader) ([]string, error) {
	var images []string
	nodes, err := Filter(r, FilterByAttr("src"))
	if err != nil {
		return nil, fmt.Errorf("while extracting nodes with attribute \"src\": %v", err)
	}

	for _, n := range nodes {
		for _, a := range n.Attr {
			if a.Key != "src" {
				continue
			}

			images = append(images, a.Val)
			break
		}
	}

	return images, nil
}

func tokens(r io.Reader) ([]*html.Node, error) {
	var nodes []*html.Node
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			if z.Err() == io.EOF {
				return nodes, nil
			}
			return nodes, fmt.Errorf("while tokenizing: %v", z.Err())
		case tt == html.StartTagToken:
			token := z.Token()
			node := html.Node{
				Type:     tokenTypeToNodeType(token.Type),
				DataAtom: token.DataAtom,
				Data:     token.Data,
				Attr:     token.Attr,
			}

			nodes = append(nodes, &node)
		}
	}
}

func tokenTypeToNodeType(tt html.TokenType) html.NodeType {
	switch tt {
	case html.TextToken:
		return html.TextNode
	case html.StartTagToken:
		return html.ElementNode
	case html.EndTagToken:
		return html.ElementNode
	case html.SelfClosingTagToken:
		return html.ElementNode
	case html.CommentToken:
		return html.CommentNode
	case html.DoctypeToken:
		return html.DoctypeNode
	default:
		return html.ErrorNode
	}
}

func forEachNode(nodes []*html.Node, filters ...FilterFn) []*html.Node {
	var filtered []*html.Node

	for _, n := range nodes {
		pass := true
		for _, f := range filters {
			if ok := f(n); !ok {
				pass = false
			}
		}

		if pass {
			filtered = append(filtered, n)
		}
	}

	return filtered
}
