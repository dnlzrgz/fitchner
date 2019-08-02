package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Option func(r *http.Request) error

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

func WithAgent(agent string) Option {
	return func(r *http.Request) error {
		if agent == "" {
			return fmt.Errorf("user agent for request not provided")
		}

		r.Header.Set("User-Agent", agent)
		return nil
	}
}

func WithBasicAuth(username string, password string) Option {
	return func(r *http.Request) error {
		r.SetBasicAuth(username, password)
		return nil
	}
}

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

func Get(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, baseURL, nil)
}

func Head(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodHead, baseURL, nil)
}

func Post(baseURL string, body []byte) (*http.Request, error) {
	return New(WithMethod(http.MethodPost), WithURL(baseURL), WithBody(body))
}

func Put(baseURL string, body []byte) (*http.Request, error) {
	return New(WithMethod(http.MethodPut), WithURL(baseURL), WithBody(body))
}

func Delete(baseURL string) (*http.Request, error) {
	return http.NewRequest(http.MethodDelete, baseURL, nil)
}
