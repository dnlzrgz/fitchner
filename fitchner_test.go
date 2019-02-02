package fitchner

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func TestFetch(t *testing.T) {
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

	if len(b) <= 0 {
		t.Fatalf("expected []byte to not be empty")
	}

	tests := []string{
		"<!DOCTYPE HTML>",
		"<head>",
		"<title>Testing</title>",
		"<h1 class=\"title\">",
		"</h1>",
		"<a href=\"https://www.google.com\"",
		"<a href=\"mailto:",
		"</html>",
	}

	for _, tt := range tests {
		if !bytes.Contains(b, []byte(tt)) {
			t.Errorf("expected to find %s on response body", tt)
		}
	}
}

func TestExtractNodes(t *testing.T) {
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

	nodes, err := ExtractNodes(b)
	if err != nil {
		t.Errorf("while extracting nodes: %v", err)
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
					Val: "mail",
				},
				html.Attribute{
					Key: "id",
					Val: "mail",
				},
			},
		},
	}

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
				t.Errorf("expected node %q to have attribute %v. got %v", n.Data, attr.Key, n.Attr[j].Key)
			}

			if attr.Val != n.Attr[j].Val {
				t.Errorf("expected node %q to have attribute %v with value %v. got %v", n.Data, attr.Key, attr.Val, n.Attr[j].Val)
			}
		}
	}
}

func TestExtractLinks(t *testing.T) {
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

	links, err := ExtractLinks(b)
	if err != nil {
		t.Errorf("while extracting nodes: %v", err)
	}

	tests := []string{"https://www.google.com", "testing@test.com"}

	if len(links) != len(tests) {
		t.Errorf("expected links to have len %v. got: %v", len(tests), len(links))
	}

	for i, tt := range tests {
		if tt != links[i] {
			t.Errorf("expected link to be %q. got: %q", tt, links[i])
		}
	}
}

func BenchmarkFetch(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		Fetch(client, req)
	}
}

func BenchmarkExtractNodes(b *testing.B) {
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
		b.Fatalf("while making a new fetch: %v", err)
	}

	for i := 0; i < b.N; i++ {
		ExtractNodes(body)
	}
}

func BenchmarkExtractLinks(b *testing.B) {
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
		b.Fatalf("while making a new fetch: %v", err)
	}

	for i := 0; i < b.N; i++ {
		ExtractLinks(body)
	}
}

func testHandler() func(w http.ResponseWriter, r *http.Request) {
	tpl := `<!DOCTYPE HTML>
	<html>
	<head>
	<title>Testing</title>
	</head>
	<body>
	<header>
	<h1 class="title">Testing</h1>
	</header>
	<a href="https://www.google.com" class="link" id="link">Links</a>
	<a href="mailto:testing@test.com" class="mail" id="mail">Mail</a>
	</body>
	</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; chatset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, tpl)
	}
}
