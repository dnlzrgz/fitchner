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
		t.Fatalf("while creating a new request: %v", err)
	}

	client := &http.Client{}
	_, err = Fetch(client, req)
	if err != nil {
		t.Errorf("%v", err)
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

func testHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; chatset=UTF-8")
		w.WriteHeader(http.StatusOK)
	}
}
