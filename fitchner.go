// Package fitchner provides utilities to make HTTP request
// and filter the response easily.
package fitchner

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Fetcher is a simple struct with an *http.Client
// and an *http.Request.
type Fetcher struct {
	c *http.Client
	r *http.Request
}

// FetcherOption defines an option for a new Fetcher.
type FetcherOption func(f *Fetcher) error

// WithClient receives an *http.Client and returns a FetcherOption
// that applies it.
func WithClient(c *http.Client) FetcherOption {
	return func(f *Fetcher) error {
		f.c = c
		return nil
	}
}

// WithRequest receives an *http.Request and returns a FetcherOption
// that applies it.
func WithRequest(r *http.Request) FetcherOption {
	return func(f *Fetcher) error {
		f.r = r
		return nil
	}
}

// WithSimpleGetRequest receives an URL and creates a simple HTTP GET request
// using it and returns a FetcherOption that applies it.
func WithSimpleGetRequest(url string) FetcherOption {
	return func(f *Fetcher) error {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		f.r = req
		return nil
	}
}

// NewFetcher returns a pointer to a Fetcher applying the options received.
// It returns an error if:
// 	* A FetcherOption returns an error.
//	* There is no http.Request provided.
// So you'll need to pass an http.Request using the WithRequest FetcherOption or
// using the WithSimpleGetRequest FetcherOption.
// If no http.Client is provided it creates and assigns a new one.
func NewFetcher(opts ...FetcherOption) (*Fetcher, error) {
	f := Fetcher{}
	for _, option := range opts {
		err := option(&f)
		if err != nil {
			return nil, fmt.Errorf("while applying option %T to Fetcher: %v", option, err)
		}
	}

	if f.r == nil {
		return nil, fmt.Errorf("while creating new Fetcher: HTTP request not provided")
	}

	if f.c == nil {
		f.c = &http.Client{}
	}

	return &f, nil
}

// Do uses the http.Client and the http.Request of a *Fetcher
// and makes an HTTP request.
// It returns an error if:
//	* The status code of the response is not 200 (OK).
//	* The Content-Type is not of type "text/html".
//	* There was an error making the request itself.
// If nothing goes wrong, it returns a []byte with the body of the response.
func (f *Fetcher) Do() ([]byte, error) {
	resp, err := f.c.Do(f.r)
	if err != nil {
		return nil, fmt.Errorf("while making a request to %q using method %q: %v", f.r.URL, f.r.Method, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code: %v", resp.StatusCode)
	}

	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, fmt.Errorf("response Content-Type is %q. expected Content-Type %q", ctype, "text/html;")

	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("while reading the response body of %q: %v", f.r.URL, err)
	}

	return b, nil
}
