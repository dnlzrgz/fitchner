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

To create a new Fetcher you'll need to use NewFetcher:

```go
func NewFetcher(opts ...FetcherOption) (*Fetcher, error)
```

NewFetche returns a pointer to a *Fetcher applying the options received.
It returns an error if:

* A FetcherOption returns an error.
* The is no http.request provided.

Define with a FetcherOption an http.Client is optional. If no http.Client is provided a new one is created and assigned to the Fetcher client (f.c). But define an http.request is obligatory.

You can pass an http.request witht the following FetcherOption:

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

As you can see at the example above. Do uses the http.client and the http.request of a *Fetcher and makes an HTTP request with them.

```go
func (f *Fetcher) Do() ([]byte, error)
```

It returns an error if:

* The status code of the response is not 200 (OK).
* The Content-Type of the response is not "text/html;".
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

Filter is a function that receives an io.Reader from which to extract the HTML nodes.

```go
func Filter(r io.Reader, filters ...FilterFn) ([]*html.Node, error)
```

It returns an error if there is any problem extracting the nodes.

You can also pass none, one or more FilterFn to manipulate the final slice.

> The order of the filters is important!

## Links

Links receives an io.Reader and, using Filter and FilterByAttr FilterFn, returns a []string with all the links found.

```go
Links(r io.Reader) ([]string, error)
```

It returns an error if there is any problem extracting the nodes.

## FilterFn

A FilterFn defines a filter for Filter.

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

## Help

Any help is welcome. So please, if there's anything I can improve on, don't hesitate to let me know.