![Version](https://img.shields.io/badge/version-0.0.0-orange.svg)
![Go](https://img.shields.io/github/go-mod/go-version/vigo/accept)
[![Documentation](https://godoc.org/github.com/vigo/accept?status.svg)](https://pkg.go.dev/github.com/vigo/accept)
[![Go Report Card](https://goreportcard.com/badge/github.com/vigo/accept)](https://goreportcard.com/report/github.com/vigo/accept)
[![golangci-lint](https://github.com/vigo/accept/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/vigo/accept/actions/workflows/golangci-lint.yml)
[![Test go code](https://github.com/vigo/accept/actions/workflows/test.yml/badge.svg)](https://github.com/vigo/accept/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/vigo/accept/graph/badge.svg?token=88BNNWA3K0)](https://codecov.io/gh/vigo/accept)
![Powered by Rake](https://img.shields.io/badge/powered_by-rake-blue?logo=ruby)

# accept

A lightweight Go library for parsing HTTP `Accept` headers and selecting the
most suitable `Content-Type`.

Documents followed for `Accept` http reques header:

- https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept
- https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
- https://developer.mozilla.org/en-US/docs/Web/HTTP/MIME_types
- https://www.iana.org/assignments/media-types/media-types.xhtml#application
- https://developer.mozilla.org/en-US/docs/Glossary/Quality_values


---

## Installation

```bash
go get -u github.com/vigo/accept
```

---

## Usage

By default, unless otherwise specified, the fallback content type is always
set to `application/json`. If desired, you can customize the fallback value
using `WithDefaultMediaType` method. Also, if request `Accept` header is `*/*`,
library matches first supported media type too.

```go
// your server supports: application/json and text/html
contentNegotiator := accept.New(
    accept.WithSupportedMediaTypes("application/json", "text/html"),
)

contentNegotiator := accept.New(
    accept.WithSupportedMediaTypes("application/json", "text/html"),
    accept.WithDefaultMediaType("text/plain"),
)

// in your http handler
acceptHeader := r.Header.Get("Accept")
contentType := cn.Negotiate(acceptHeader)
```

Full example, your server supports: `application/json`, `text/html` and 
`text/plain` and for unmatched `Accept` header, will use `text/plain` as
default fallback value.

```go
# main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vigo/accept"
)

func handlerFunc(cn *accept.ContentNegotiation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acceptHeader := r.Header.Get("Accept")

		contentType := cn.Negotiate(acceptHeader)
		w.Header().Set("Content-Type", contentType)

		switch contentType {
		case "application/json":
			response := map[string]string{"message": "OK"}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
			}
		case "text/html":
			_, _ = w.Write([]byte("<p>\n\t<span>message<span>\n\t<pre>OK</pre>\n</p>\n"))
		default:
			_, _ = w.Write([]byte("message: OK\n"))
		}
	}
}

func main() {
	contentNegotiator := accept.New(
		accept.WithSupportedMediaTypes("application/json", "text/html", "text/plain"),
		accept.WithDefaultMediaType("text/plain"),
	)

	log.Println("starting server at :8080")
	http.HandleFunc("/", handlerFunc(contentNegotiator))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Run the server:

```bash
go run main.go
```

Test the requests:

```bash
curl localhost:8080 -H 'Accept: text/html'
curl localhost:8080 -H 'Accept: text/plain'
curl localhost:8080 -H 'Accept: text/markdown'     # fallback to "text/plain"
curl localhost:8080 -H 'Accept: application/json'
curl localhost:8080                                # curl sends */*, first match is "application/json"
```

---

## Rake Tasks

```bash
rake -T

rake coverage  # show test coverage
rake test      # run test
```

---

## Contributor(s)

* [Uğur Özyılmazel](https://github.com/vigo) - Creator, maintainer

---

## Contribute

All PR’s are welcome!

1. `fork` (https://github.com/vigo/accept/fork)
1. Create your `branch` (`git checkout -b my-feature`)
1. `commit` yours (`git commit -am 'add some functionality'`)
1. `push` your `branch` (`git push origin my-feature`)
1. Than create a new **Pull Request**!

---

## License

This project is licensed under MIT (MIT)

---

This project is intended to be a safe, welcoming space for collaboration, and
contributors are expected to adhere to the [code of conduct][coc].

[coc]: https://github.com/vigo/accept/blob/main/CODE_OF_CONDUCT.md
