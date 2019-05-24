package fitchner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetcherDo(t *testing.T) {
	handler := testDoHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("while creating a new request for testing: %v", err)
	}

	tests := []struct {
		name           string
		options        []FetcherOption
		expectedToFail bool
	}{
		{
			name: "client and request options",
			options: []FetcherOption{
				WithClient(client),
				WithRequest(req),
			},
			expectedToFail: false,
		},
		{
			name:           "no options",
			options:        []FetcherOption{},
			expectedToFail: true,
		},
		{
			name: "simple get request",
			options: []FetcherOption{
				WithSimpleGetRequest(server.URL),
			},
			expectedToFail: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewFetcher(tc.options...)
			if err != nil {
				if tc.expectedToFail {
					t.Skipf("expected to fail test: %v", err)
				}

				t.Fatalf("while creating a new Fetcher not expected to fail: %v", err)
			}

			if _, err := f.Do(); err != nil {
				t.Fatalf("while fetching: %v", err)
			}
		})
	}
}

func testDoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello, World!")
	}
}

func BenchmarkFetcherDo(b *testing.B) {
	handler := testDoHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	bc := func(f *Fetcher) error {
		client := &http.Client{}
		f.c = client
		return nil
	}

	br := func(f *Fetcher) error {
		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			b.Fatalf("while creating a new request for benchmark Fetcher: %v", err)
		}

		f.r = req
		return nil
	}

	f, err := NewFetcher(bc, br)
	if err != nil {
		b.Fatalf("while creating a new Fetchner for benchmark: %v", err)
	}

	for i := 0; i < b.N; i++ {
		f.Do()
	}
}
