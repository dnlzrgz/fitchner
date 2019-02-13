package fitchner

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"
)

func TestFetch(t *testing.T) {
	b := testFetch(t)

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
	<a href="mailto:testing@test.com" class="link mail" id="mail">Mail</a>
	</body>
	</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; chatset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, tpl)
	}
}

func testFetch(t *testing.T) []byte {
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

	return b
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
