package fitchner

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func TestFilterAttr(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	b, err := Fetch(client, req)
	if err != nil {
		t.Errorf("while making a new fetch: %v", err)
	}

	nodes, err := Nodes(b)
	if err != nil {
		t.Errorf("while extracting nodes: %v", err)
	}

	tests := []string{"h1", "a", "a"}
	filtered := FilterAttr(nodes, "class")
	for i, tt := range tests {
		if tt != filtered[i].Data {
			t.Errorf("expected %q. got %q", tt, filtered[i].Data)
		}
	}

	noFiltered := FilterAttr(nodes, "")
	if !testEqual(nodes, noFiltered) {
		t.Errorf("FilterAttr returns a different []*html.Node when no attr specified")
	}
}

func TestFilterTag(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	b, err := Fetch(client, req)
	if err != nil {
		t.Errorf("while making a new fetch: %v", err)
	}

	nodes, err := Nodes(b)
	if err != nil {
		t.Errorf("while extracting nodes: %v", err)
	}

	filtered := FilterTag(nodes, "h1")
	if len(filtered) > 1 {
		t.Errorf("filtered should contain only 1 element. got %v", len(filtered))
	}

	if filtered[0].Data != "h1" {
		t.Errorf("filtered should contain %q. got %q", "h1", filtered[0].Data)
	}

	noFiltered := FilterTag(nodes, "")
	if !testEqual(nodes, noFiltered) {
		t.Errorf("FilterTag returns a different []*html.Node when no tag specified")
	}
}

func TestLinks(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	b, err := Fetch(client, req)
	if err != nil {
		t.Errorf("while making a new fetch: %v", err)
	}

	nodes, err := Nodes(b)
	if err != nil {
		t.Errorf("while extracting nodes: %v", err)
	}

	tests := []string{"https://www.google.com", "mailto:testing@test.com"}
	links := Links(nodes)
	for i, tt := range tests {
		if tt != links[i] {
			t.Errorf("expected %s. got %s", tt, links[i])
		}
	}
}

func testEqual(a, b []*html.Node) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
