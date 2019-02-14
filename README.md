# Fitchner

[![GoDoc](https://godoc.org/github.com/danielkvist/fitchner?status.svg)](https://godoc.org/github.com/danielkvist/fitchner)

Fitchner is a Go module which facilities the work with HTTP requests and the extraction of information from HTTP responses.

## Utilities

### Fetch

Fetch needs and http.Client and a http.Request to make an HTTP request with the specified client.
It returns an error in case of:

- If the client fails to make the requests.
- If it gets a non-2xx response.
- If there is a error while reading the response body with ioutil.ReadAll.

If neither of the three cases described above occurs returns a []byte with the body of the response to que client's request.

```go
func Fetch(c *http.Client, req *http.Request) ([]byte, error)
```

### Filter

Filter need to receive a []byte with a response body.
The following params are optional and each of them modifies the returned slice of nodes ([]\*html.Node):

- `tag`.
- `attr`.
- `val`.

> NOTE: The tag doesn't has to include "<" neither ">".

It returns an error if there is any problem while parsing the body of the response.

```go
func Filter(b []byte, tag, attr, val string) ([]*html.Node, error)
```

### Links

Links returns all the links found on the body of the response that needs to receive. Or, an error if there is any problem while parsing the body of the response.

> Note: Links ignores links with the "tel:" prefix and, also removes the "mailto:" prefix if any.

```go
func Links(b []byte) ([]string, error)
```
