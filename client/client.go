package client

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client *http.Client
}

type Option func(c *Client) error

func WithTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.client.Timeout = t
		return nil
	}
}

func WithProxy(proxy string) Option {
	return func(c *Client) error {
		if proxy == "" {
			c.client.Transport = nil
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

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) DoAfter(req *http.Request, t time.Duration) (*http.Response, error) {
	time.Sleep(t)
	return c.Do(req)
}

func (c *Client) Ping(req *http.Request) (int, string, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, resp.Status, nil
}

func (c *Client) PingAfter(req *http.Request, t time.Duration) (int, string, error) {
	time.Sleep(t)
	return c.Ping(req)
}
