package fitchner

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

func BenchmarkFetch(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("error creating a new request: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := Fetch(req)
		if err != nil {
			b.Fatalf("error: %v", err)
		}
	}
}

func testHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; chatset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}
