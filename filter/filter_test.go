package filter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielkvist/fitchner/client"
	"github.com/danielkvist/fitchner/request"
)

func TestFilter(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	tests := []struct {
		name             string
		filters          []FilterFn
		expectedTagItems []string
	}{
		{
			name:    "no filters",
			filters: nil,
			expectedTagItems: []string{
				"html",
				"head",
				"title",
				"body",
				"h1",
				"a",
				"section",
				"a",
				"div",
				"img",
				"span",
				"img",
			},
		},
		{
			name:             "filtering by tag",
			filters:          []FilterFn{FilterByTag("h1")},
			expectedTagItems: []string{"h1"},
		},
		{
			name:             "filtering by class",
			filters:          []FilterFn{FilterByClass("link")},
			expectedTagItems: []string{"a", "a"},
		},
		{
			name:             "filtering by id",
			filters:          []FilterFn{FilterByID("title")},
			expectedTagItems: []string{"h1"},
		},
		{
			name: "filtering by tag and class",
			filters: []FilterFn{
				FilterByClass("home"),
				FilterByTag("a"),
			},
			expectedTagItems: []string{"a"},
		},
		{
			name:             "filtering by attribute",
			filters:          []FilterFn{FilterByAttr("href")},
			expectedTagItems: []string{"a", "a"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := request.Get(server.URL)
			if err != nil {
				t.Fatalf("while creating a new request: %v", err)
			}

			c, err := client.New()
			if err != nil {
				t.Fatalf("while creating a new client: %v", err)
			}

			resp, err := c.Do(req)
			if err != nil {
				t.Fatalf("while making a request: %v", err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("while reading the response body: %v", err)
			}

			r := bytes.NewReader(b)
			filtered, err := Filter(r, tc.filters...)
			if err != nil {
				t.Fatalf("while filtering: %v", err)
			}

			if len(filtered) != len(tc.expectedTagItems) {
				t.Fatalf("expected %v items filtered. got=%v", len(tc.expectedTagItems), len(filtered))
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
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	expectedLinks := []string{"/", "https://golang.org"}

	req, err := request.Get(server.URL)
	if err != nil {
		t.Fatalf("while creating a new request: %v", err)
	}

	c, err := client.New()
	if err != nil {
		t.Fatalf("while creating a new client: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("while making a request: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("while reading the response body: %v", err)
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

func TestImages(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	expectedImages := []string{"cat.jpeg", "btn.png"}

	req, err := request.Get(server.URL)
	if err != nil {
		t.Fatalf("while creating a new request: %v", err)
	}

	c, err := client.New()
	if err != nil {
		t.Fatalf("while creating a new client: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		t.Fatalf("while making a request: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("while reading the response body: %v", err)
	}

	r := bytes.NewReader(b)
	images, err := Images(r)
	if err != nil {
		t.Fatalf("while filtering: %v", err)
	}

	for i, ei := range images {
		if ei != expectedImages[i] {
			t.Fatalf("expected image source to be %q. got=%q", ei, expectedImages[i])
		}
	}
}

func testHandler() func(w http.ResponseWriter, r *http.Request) {
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
			<div>
				<img src="cat.jpeg">
				<span>
					<img src="btn.png">
				</span>
			</div>
		</body>
		</html>`

		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, tpl)
	}
}

func BenchmarkFilterWithoutFilters(b *testing.B) {
	benchmarkFilter(b)
}

func BenchmarkFilterByClass(b *testing.B) {
	benchmarkFilter(b, FilterByClass("link"))
}

func BenchmarkFilterByTagAndClass(b *testing.B) {
	benchmarkFilter(b, FilterByTag("a"), FilterByClass("home"))
}

func BenchmarkFilterByClassAndTag(b *testing.B) {
	benchmarkFilter(b, FilterByClass("home"), FilterByTag("a"))
}

func benchmarkFilter(b *testing.B, filters ...FilterFn) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := request.Get(server.URL)
	if err != nil {
		b.Fatalf("while creating a new request: %v", err)
	}

	c, err := client.New()
	if err != nil {
		b.Fatalf("while creating a new client: %v", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		b.Fatalf("while making a request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		b.Fatalf("while reading the response body: %v", err)
	}

	r := bytes.NewReader(body)
	for i := 0; i < b.N; i++ {
		Filter(r, filters...)
	}
}
