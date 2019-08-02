package request

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	req *http.Request
}

type Option func(r *Request) error

func WithMethod(method string) Option {
	return func(r *Request) error {
		if method == "" {
			r.req.Method = http.MethodGet
			return nil
		}

		r.req.Method = method
		return nil
	}
}

func WithURL(baseURL string) Option {
	return func(r *Request) error {
		if baseURL == "" {
			return fmt.Errorf("URL for request not provided")
		}

		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}

		r.req.URL = parsedURL
		return nil
	}
}

func WithAgent(agent string) Option {
	return func(r *Request) error {
		if agent == "" {
			return fmt.Errorf("user agent for request not provided")
		}

		r.req.Header.Set("User-Agent", agent)
		return nil
	}
}

func WithBasicAuth(username string, password string) Option {
	return func(r *Request) error {
		r.req.SetBasicAuth(username, password)
		return nil
	}
}

func WithBody(body []byte) Option {
	return func(r *Request) error {
		if body == nil {
			return nil
		}

		r.req.ContentLength = int64(len(body))
		r.req.Body = ioutil.NopCloser(bytes.NewReader(body))

		return nil
	}
}

func New(opts ...Option) (*Request, error) {
	r := &Request{
		req: &http.Request{},
	}

	for _, option := range opts {
		err := option(r)
		if err != nil {
			return nil, fmt.Errorf("while applying option to Request: %v", err)
		}
	}

	return r, nil
}

func (r *Request) Req() *http.Request {
	return r.req
}

func Get(baseURL string) (*Request, error) {
	return New(WithMethod(http.MethodGet), WithURL(baseURL))
}

func Head(baseURL string) (*Request, error) {
	return New(WithMethod(http.MethodHead), WithURL(baseURL))
}

func Post(baseURL string, body []byte) (*Request, error) {
	return New(WithMethod(http.MethodPost), WithURL(baseURL), WithBody(body))
}

func Put(baseURL string, body []byte) (*Request, error) {
	return New(WithMethod(http.MethodPut), WithURL(baseURL), WithBody(body))
}

func Delete(baseURL string) (*Request, error) {
	return New(WithMethod(http.MethodDelete), WithURL(baseURL))
}
