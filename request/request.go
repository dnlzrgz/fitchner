package request

import (
	"fmt"
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
