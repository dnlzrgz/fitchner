package fitchner

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFilter(t *testing.T) {
	handler := testFilter()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tests := []struct {
		name                string
		filters             []FilterFn
		expectedNumberItems int
		expectedTagItems    []string
	}{
		{
			name:                "no filters",
			filters:             nil,
			expectedNumberItems: 8,
			expectedTagItems:    []string{"html", "head", "title", "body", "h1", "a", "section", "a"},
		},
		{
			name:                "filtering by tag",
			filters:             []FilterFn{FilterByTag("h1")},
			expectedNumberItems: 1,
			expectedTagItems:    []string{"h1"},
		},
		{
			name:                "filtering by class",
			filters:             []FilterFn{FilterByClass("link")},
			expectedNumberItems: 2,
			expectedTagItems:    []string{"a", "a"},
		},
		{
			name:                "filtering by id",
			filters:             []FilterFn{FilterByID("title")},
			expectedNumberItems: 1,
			expectedTagItems:    []string{"h1"},
		},
		{
			name: "filtering by tag and class",
			filters: []FilterFn{
				FilterByClass("home"),
				FilterByTag("a"),
			},
			expectedNumberItems: 1,
			expectedTagItems:    []string{"a"},
		},
		{
			name:                "filtering by attribute",
			filters:             []FilterFn{FilterByAttr("href")},
			expectedNumberItems: 2,
			expectedTagItems:    []string{"a", "a"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFetcher(WithSimpleGetRequest(server.URL))
			if err != nil {
				t.Fatalf("while creating a new Fetcher: %v", err)
			}

			b, err := f.Do()
			if err != nil {
				t.Fatalf("while fetching for testing: %v", err)
			}

			r := bytes.NewReader(b)
			filtered, err := Filter(r, tc.filters...)
			if err != nil {
				t.Fatalf("while filtering: %v", err)
			}

			if len(filtered) != tc.expectedNumberItems {
				t.Fatalf("expected %v items filtered. got=%v", tc.expectedNumberItems, len(filtered))
			}

			for i, n := range filtered {
				if n.Data != tc.expectedTagItems[i] {
					t.Fatalf("expected item with tag %q. got item with tag %q", tc.expectedTagItems[i], n.Data)
				}
			}
		})
	}
}

func TestLinks(t *testing.T) {
	handler := testFilter()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	expectedLinks := []string{"/", "https://golang.org"}

	f, err := NewFetcher(WithSimpleGetRequest(server.URL))
	if err != nil {
		t.Fatalf("while creating a new Fetcher: %v", err)
	}

	b, err := f.Do()
	if err != nil {
		t.Fatalf("while fetching for testing: %v", err)
	}

	r := bytes.NewReader(b)
	links, err := Links(r)
	if err != nil {
		t.Fatalf("while filtering: %v", err)
	}

	for i, el := range expectedLinks {
		if el != links[i] {
			t.Fatalf("expected link to be %q. got=%q", el, links[i])
		}
	}

}

func testFilter() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl := `
		<!DOCTYPE HTML>
		<html>
		<head>
			<title>Testing</title>
		</head>
		<body>
			<h1 id="title">Testing</h1>
			<a href="/" class="home link">Home</a>
			<section>
				<a href="https://golang.org" alt="google" class="link">Golang</a>
			</section>
		</body>
		</html>`

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, tpl)
	}
}

func BenchmarkFilterNoFilters(b *testing.B) {
	handler := testFilter()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	f, err := NewFetcher(WithSimpleGetRequest(server.URL))
	if err != nil {
		b.Fatalf("while creating a new Fetcher: %v", err)
	}

	data, err := f.Do()
	if err != nil {
		b.Fatalf("while fetching for testing: %v", err)
	}

	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		Filter(r)
	}
}

func BenchmarkFilterByClass(b *testing.B) {
	handler := testFilter()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	f, err := NewFetcher(WithSimpleGetRequest(server.URL))
	if err != nil {
		b.Fatalf("while creating a new Fetcher: %v", err)
	}

	data, err := f.Do()
	if err != nil {
		b.Fatalf("while fetching for testing: %v", err)
	}

	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		Filter(r, FilterByClass("link"))
	}
}

func BenchmarkFilterByTagAndClass(b *testing.B) {
	handler := testFilter()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	f, err := NewFetcher(WithSimpleGetRequest(server.URL))
	if err != nil {
		b.Fatalf("while creating a new Fetcher: %v", err)
	}

	data, err := f.Do()
	if err != nil {
		b.Fatalf("while fetching for testing: %v", err)
	}

	r := bytes.NewReader(data)
	for i := 0; i < b.N; i++ {
		Filter(r, []FilterFn{FilterByTag("a"), FilterByClass("home")}...)
	}
}
