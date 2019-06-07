# Fitchner

[![GoDoc](https://godoc.org/github.com/danielkvist/fitchner?status.svg)](https://godoc.org/github.com/danielkvist/fitchner)
[![Go Report Card](https://goreportcard.com/badge/github.com/danielkvist/fitchner)](https://goreportcard.com/report/github.com/danielkvist/fitchner)

Fitchner provides utilities to make HTTP request and filter the response easily.

## Get

```bash
go get -u "github.com/danielkvist/fitchner"
```

## Example

```go
f, err := fitchner.NewFetcher(fitchner.WithSimpleGetRequest("https://github.com"))
if err != nil {
    log.Fatalf("while creating a new Fetcher: %v", err)
}

data, err := f.Do()
if err != nil {
    log.Fatalf("while making HTTP request: %v", err)
}

fmt.Println(string(data))

r := bytes.NewReader(data)
nodes, err := fitchner.Filter(r, fitchner.FilterByTag("h1"))
if err != nil {
    log.Fatalf("while filtering: %v", err)
}

for _, n := range nodes {
    fmt.Println(n.Data)
}
```

## Fetcher

A Fetcher is a simple struct with an *http.Client and an *http.Request:

```go
type Fetcher struct {
    c *http.Client
    r *http.Request
}
```

To create a new Fetcher you'll need to use the NewFetcher function:

```go
func NewFetcher(opts ...FetcherOption) (*Fetcher, error)
```

NewFetcher returns a pointer to a *Fetcher applying the options received.
It returns an error if:

* A FetcherOption function returns an error.
* The is no http.request provided.

Define with a FetcherOption an http.Client is optional. If no http.Client is provided a new one is created and assigned to the Fetcher client (f.c). But define an http.Request is obligatory.

You can pass an http.Request with the followings FetcherOption functions:

```go
    optReq := fitchner.WithSimpleGetRequest("https://google.com")
    f, _ := fitchner.NewFetcher(optReq)
```

Or:

```go
    req, _ := http.NewRequest(http.MethodGet, "https://google.com", nil)

    optReq := fitchner.WithRequest(req)
    f, _ := fitchner.NewFetcher(optReq)
```

## Do

As you can see at the example above. Do uses the http.Client and the http.Request of a *Fetcher and makes an HTTP request with them.

```go
func (f *Fetcher) Do() ([]byte, error)
```

It returns an error if:

* The status code of the response is not 200 (OK).
* The Content-Type of the response is not "text/html".
* There was an error making the request itself.

If nothing goes wrong, it returns a []byte with the body of the response.

## FetcherOption

A FetcherOption is a simple function that receives a *Fetcher and defines an option for NewFetcher:

```go
type FetcherOption func(f *Fetcher) error
```

An example is:

```go
func WithClient(c *http.Client) FetcherOption {
    return func(f *Fetcher) error {
        f.c = c
        return nil
    }
}
```

## Filter

Filter is a function that receives an io.Reader from which to extract HTML nodes.

```go
func Filter(r io.Reader, filters ...FilterFn) ([]*html.Node, error)
```

It returns an error if there is any problem extracting the nodes.

You can also pass none, one or more FilterFn functions to manipulate the final slice.

> The order of the filters is important!

## Links

Links receives an io.Reader and, using the Filter function and FilterByAttr FilterFn, returns a []string with all the links found.

```go
Links(r io.Reader) ([]string, error)
```

## Images

Images receives an io.Reader an, as Links, uses the Filter function and FilterByAttr FilterFn to extract a []string with all the images sources found.

```go
Images(r io.Reader) ([]string, error)
```

## FilterFn

A FilterFn defines a filter to filter HTML nodes. If the filter returns true, the HTML node is added to the result of the Filter function. And if the filter returns false the HTML node is excluded.

```go
type FilterFn func(n *html.Node) bool
```

Some predefined Filters that you can use are:

```go
func FilterByTag(tag string) FilterFn
func FilterByClass(class string) FilterFn
func FilterByID(id string) FilterFn
func FilterByAttr(attr string) FilterFn
```