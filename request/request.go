// Package request provides small utilities to work with http.Response easily.
package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Option defines an option for a new *http.Request.
type Option func(r *http.Request) error

// WithMethod receives an HTTP method
// and assigns it to the *http.Request. If the received
// method is empty it assigns an http.MethodGet.
func WithMethod(method string) Option {
	return func(r *http.Request) error {
		if method == "" {
			r.Method = http.MethodGet
			return nil
		}

		r.Method = method
		return nil
	}
}

// WithURL receives a base URL and assigns it
// to the *http.Request's URL. If not URL is provided
// it returns a non-nil error.
func WithURL(baseURL string) Option {
	return func(r *http.Request) error {
		if baseURL == "" {
			return fmt.Errorf("URL for request not provided")
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		r.URL = parsedURL
		return nil
	}
}

// WithAgent receives an user agent and uses it
// as an "User-Agent" head for the *http.Request.
func WithAgent(agent string) Option {
	return func(r *http.Request) error {
		if agent == "" {
			return fmt.Errorf("user agent for request not provided")
		}

		r.Header.Set("User-Agent", agent)
		return nil
	}
}

// WithBasicAuth receives an username and a password
// and assigns it to the *http.Request to use
// basic authentication.
func WithBasicAuth(username string, password string) Option {
	return func(r *http.Request) error {
		r.SetBasicAuth(username, password)
		return nil
	}
}

// WithBody receives a []byte and uses it to set the
// ContentLength of the *http.Request and the Body itself.
func WithBody(body []byte) Option {
	return func(r *http.Request) error {
		if body == nil {
			return nil
		}

		r.ContentLength = int64(len(body))
		r.Body = ioutil.NopCloser(bytes.NewReader(body))

		return nil
	}
}

// New returns a new *http.Request or an error if
// any error occurs while applying the received Options.
func New(opts ...Option) (*http.Request, error) {
	r := &http.Request{}

	for _, option := range opts {
		err := option(r)
		if err != nil {
			return nil, fmt.Errorf("while applying option to Request: %v", err)
		}
	}

	return r, nil
}

// Get returns a new *http.Request with an HTTP GET method
// to the received URL or an error if something goes wrong.
func Get(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, baseURL, nil)
}

// Head returns a new *http.Request with an HTTP HEAD method
// to the received URL or an error if something goes wrong.
func Head(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodHead, baseURL, nil)
}

// Post returns a new *http.Method with an HTTP POST method
// using the received body to the received URL
// or returns an error if something goes wrong.
func Post(baseURL string, body []byte) (*http.Request, error) {
	return New(WithMethod(http.MethodPost), WithURL(baseURL), WithBody(body))
}

// Put returns a new *http.Method with an HTTP PUT method
// using the received body to the received URL
// or returns an error if something goes wrong.
func Put(baseURL string, body []byte) (*http.Request, error) {
	return New(WithMethod(http.MethodPut), WithURL(baseURL), WithBody(body))
}

// Delete returns a new *http.Method with an HTTP DELETE method
// to the received URL or an error if something goes wrong.
func Delete(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodDelete, baseURL, nil)
}
