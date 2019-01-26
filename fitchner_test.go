package fitchner

import (
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
		t.Errorf("error creating a new request: %v", err)
	}

	resp, err := Fetch(req)
	if err != nil {
		t.Errorf("bad response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v as status code but got %v instead", http.StatusOK, resp.StatusCode)
	}
}

func TestHTMLParser(t *testing.T) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Errorf("error creating a new request: %v", err)
	}

	resp, err := Fetch(req)
	if err != nil {
		t.Errorf("bad response: %v", err)
	}

	tests := []struct {
		expectedType html.TokenType
		expectedTag  string
		expectedAttr []html.Attribute
	}{
		{
			expectedType: html.StartTagToken,
			expectedTag:  "h1",
			expectedAttr: []html.Attribute{
				html.Attribute{
					Key: "class",
					Val: "title",
				},
			},
		},
		{
			expectedType: html.StartTagToken,
			expectedTag:  "a",
			expectedAttr: []html.Attribute{
				html.Attribute{
					Key: "href",
					Val: "https://www.google.com",
				},
				html.Attribute{
					Key: "class",
					Val: "link",
				},
			},
		},
	}

	parsed := HTMLParser(resp)
	for i, tp := range parsed {
		tt := tests[i]
		if tt.expectedType != tp.Type {
			t.Errorf("expected %s as type. got %s", tt.expectedType, tp.Type)
		}

		if tt.expectedTag != tp.Data {
			t.Errorf("expected %s as tag. got %s", tt.expectedTag, tp.Data)
		}

		if len(tt.expectedAttr) != len(tp.Attr) {
			t.Errorf("expected %v attributes. got %v", len(tt.expectedAttr), len(tp.Attr))
		}

		for j := 0; j < len(tt.expectedAttr); j++ {
			if tt.expectedAttr[j] != tp.Attr[j] {
				t.Errorf("expected attr %q with the value %q. got %q with value %q", tt.expectedAttr[j].Key, tt.expectedAttr[j].Val, tp.Attr[j].Key, tp.Attr[j].Val)
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
		b.Fatalf("error creating a new request: %v", err)
	}

	for i := 0; i < b.N; i++ {
		resp, err := Fetch(req)
		if err != nil {
			b.Fatalf("error: %v", err)
			resp.Body.Close()
		}
		resp.Body.Close()
	}
}

func BenchmarkHTMLParser(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("error creating a new request: %v", err)
	}

	resp, err := Fetch(req)
	if err != nil {
		b.Errorf("bad response: %v", err)
	}

	for i := 0; i < b.N; i++ {
		HTMLParser(resp)
	}
}

func testHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; chatset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<h1 class=\"title\">test server</h1>")
		fmt.Fprintf(w, "<a href=\"https://www.google.com\" class=\"link\">link</a>")
	}
}
