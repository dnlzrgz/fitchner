// Package client provides small utilities to work with http.Client easily.
package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Client is a simple wrapper
// for an *http.Client.
type Client struct {
	client *http.Client
}

// Option defines an option for a new Client.
type Option func(c *Client) error

// WithTimeout receives a time duration
// and assigns it to the *http.Client's Timeout
// inside the Client.
func WithTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.client.Timeout = t
		return nil
	}
}

// WithProxy receives a proxy's URL
// and assigns it to the *http.Client's Transport
// inside the Client. If the proxy's URL is an empty string
// it returns nil directly.
func WithProxy(proxy string) Option {
	return func(c *Client) error {
		if proxy == "" {
			return nil
		}

		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return err
		}

		tr := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		c.client.Transport = tr
		return nil
	}
}

// New returns a new *Client or an error
// if any error occurs while applying the received Options.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		client: &http.Client{},
	}

	for _, option := range opts {
		err := option(c)
		if err != nil {
			return nil, fmt.Errorf("while applying option to Client: %v", err)
		}
	}

	return c, nil
}

// Do sends an HTTP request using the *http.Client inside the Client
// and returns the HTTP response or a non-nil error if something goes wrong.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// DoAfter sends an HTTP request using the *http.Client inside the Client
// after the received time duration and returns the HTTP response or a non-nil
// error if something goes wrong.
func (c *Client) DoAfter(req *http.Request, t time.Duration) (*http.Response, error) {
	time.Sleep(t)
	return c.Do(req)
}

// Ping sends an HTTP request using the *http.Client inside the Client
// and returns the status and the status code of the HTTP response
// or a non-nil error if something goes wrong.
func (c *Client) Ping(req *http.Request) (int, string, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, resp.Status, nil
}

// PingAfter sends an HTTP request using the *http.Client inside the Client
// and returns the status and the status code of the HTTP response
// or a non-nil error if something goes wrong.
func (c *Client) PingAfter(req *http.Request, t time.Duration) (int, string, error) {
	time.Sleep(t)
	return c.Ping(req)
}
