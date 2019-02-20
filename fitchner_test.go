package fitchner

import (
	"bufio"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const tpl = `<!DOCTYPE HTML>
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

func TestFetch(t *testing.T) {
	r := testFetch(t)
	tests := strings.Split(tpl, "\n")

	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		l := scanner.Text()
		if ok := strings.Contains(l, tests[i]); !ok {
			t.Errorf("expected to find %q on line %v. got: %q", tests[i], i, l)
		}
		i++
	}
}

func testHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, tpl)
	}
}

func testFetch(t *testing.T) io.Reader {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("while creating a new request for testing: %v", err)
	}

	client := &http.Client{}
	b, err := Fetch(client, req)
	if err != nil {
		t.Errorf("while fetching from the test server: %v", err)
	}

	return b
}

func BenchmarkFetch(b *testing.B) {
	handler := testHandler()
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		b.Fatalf("while creating a new request for benchmarking: %v", err)
	}

	client := &http.Client{}
	for i := 0; i < b.N; i++ {
		Fetch(client, req)
	}
}
