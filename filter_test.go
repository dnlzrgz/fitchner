package fitchner

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func TestFilter(t *testing.T) {
	b := testFetch(t)
	nodes, err := Filter(b, "a", "class", "mail")
	if err != nil {
		t.Errorf("while filtering: %v", err)
	}

	tests := []struct {
		data string
		attr []html.Attribute
	}{
		{
			data: "a",
			attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "mailto:testing@test.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link mail",
				},
				html.Attribute{
					Key: "id",
					Val: "mail",
				},
			},
		},
	}

	testFilter(t, nodes, tests)
}

func TestFilterEmpty(t *testing.T) {
	b := testFetch(t)
	nodes, err := Filter(b, "", "", "")
	if err != nil {
		t.Errorf("while filtering: %v", err)
	}

	tests := []struct {
		data string
		attr []html.Attribute
	}{
		{data: "html"},
		{data: "head"},
		{data: "title"},
		{data: "body"},
		{data: "header"},
		{
			data: "h1",
			attr: []html.Attribute{
				html.Attribute{
					Key: "class",
					Val: "title",
				},
			},
		},
		{
			data: "a",
			attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "https://www.google.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link",
				},
				html.Attribute{
					Key: "id",
					Val: "link",
				},
			},
		},
		{
			data: "a",
			attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "mailto:testing@test.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link mail",
				},
				html.Attribute{
					Key: "id",
					Val: "mail",
				},
			},
		},
	}

	testFilter(t, nodes, tests)
}

func TestFilterTag(t *testing.T) {
	b := testFetch(t)
	nodes, err := Filter(b, "h1", "", "")
	if err != nil {
		t.Errorf("while filtering: %v", err)
	}

	tests := []struct {
		data string
		attr []html.Attribute
	}{
		{
			data: "h1",
			attr: []html.Attribute{
				html.Attribute{
					Key: "class",
					Val: "title",
				},
			},
		},
	}

	testFilter(t, nodes, tests)
}

func TestFilterAttr(t *testing.T) {
	b := testFetch(t)
	nodes, err := Filter(b, "", "class", "")
	if err != nil {
		t.Errorf("while filtering: %v", err)
	}

	tests := []struct {
		data string
		attr []html.Attribute
	}{
		{
			data: "h1",
			attr: []html.Attribute{
				html.Attribute{
					Key: "class",
					Val: "title",
				},
			},
		},
		{
			data: "a",
			attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "https://www.google.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link",
				},
				html.Attribute{
					Key: "id",
					Val: "link",
				},
			},
		},
		{
			data: "a",
			attr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "mailto:testing@test.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link mail",
				},
				html.Attribute{
					Key: "id",
					Val: "mail",
				},
			},
		},
	}

	testFilter(t, nodes, tests)
}

func TestFilterVal(t *testing.T) {
	b := testFetch(t)
	nodes, err := Filter(b, "", "", "title")
	if err != nil {
		t.Errorf("while filtering: %v", err)
	}

	tests := []struct {
		data string
		attr []html.Attribute
	}{
		{
			data: "h1",
			attr: []html.Attribute{
				html.Attribute{
					Key: "class",
					Val: "title",
				},
			},
		},
	}

	testFilter(t, nodes, tests)
}

func TestLinks(t *testing.T) {
	b := testFetch(t)
	links, err := Links(b)
	if err != nil {
		t.Errorf("while extracting links: %v", err)
	}

	tests := []string{"https://www.google.com", "testing@test.com"}
	if len(tests) != len(links) {
		t.Errorf("links expected to have len %v. got %v", len(tests), len(links))
	}

	for i, tt := range tests {
		if tt != links[i] {
			t.Errorf("expected %q link. got %q instead", tt, links[i])
		}
	}
}

func testFilter(t *testing.T, nodes []*html.Node, tests []struct {
	data string
	attr []html.Attribute
}) {
	for i, tt := range tests {
		n := nodes[i]

		if tt.data != n.Data {
			t.Errorf("expected node %q. got %q", tt.data, n.Data)
		}

		if len(tt.attr) != len(n.Attr) {
			t.Errorf("expected node %q to have %v attributes. got %v", n.Data, len(tt.attr), len(n.Attr))
		}

		for j, attr := range tt.attr {
			if attr.Key != n.Attr[j].Key {
				t.Errorf("expected node %q to have attribute %q. got %q", n.Data, attr.Key, n.Attr[j].Key)
			}

			if attr.Val != n.Attr[j].Val {
				t.Errorf("expected node %q to have attribute %q with value %q. got %q", n.Data, attr.Key, attr.Val, n.Attr[j].Val)
			}
		}
	}
}

func BenchmarkFilter(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	body, err := Fetch(client, req)
	if err != nil {
		b.Errorf("while making a new fetch: %v", err)
	}

	for i := 0; i < b.N; i++ {
		Filter(body, "a", "class", "link")
	}
}

func BenchmarkLinks(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	body, err := Fetch(client, req)
	if err != nil {
		b.Errorf("while making a new fetch: %v", err)
	}

	for i := 0; i < b.N; i++ {
		Links(body)
	}
}
