# Fitchner

[![GoDoc](https://godoc.org/github.com/danielkvist/fitchner?status.svg)](https://godoc.org/github.com/danielkvist/fitchner)

Fitchner provides utilities to make HTTP requests and to extract information from the responses.

## Import

```go
import "github.com/danielkvist/fitchner"
```

## Utilities

### Fetch

Fetch makes an HTTP request and returns the response.
Receives an http.Client and an http.Request for it.
An error is returned if:

- The client fails to make the request.
- There is some problem while reading the response body.
- There is a non-2xxx response.

When there is no error, returns an io.Reader with the body of the response.

```go
func Fetch(c *http.Client, req *http.Request) (*bytes.Reader, error)
```

### Filter

Filter filters the body of an HTTP response depending on the received params.
All three params (tag, attr or val) can be an empty string. In which case a
slice of \*html.Node with all the found html.ElementNode is returned.
An error is returned if there is any problem while parsing.

> NOTE: The tag doesn't has to include "<" neither ">".

```go
func Filter(r io.Reader, tag, attr, val string) ([]*html.Node, error)
```

### Links

Links extracts all the links found on the body of an HTTP response.
It ignores the links with the prefix "tel:" and removes the "mailto:" prefix.
An error is returned if there is any problem while parsing.

```go
func Links(r io.Reader) ([]string, error)
```
