<p align="center"><img src="https://raw.githubusercontent.com/go-syndication/brand/main/social/go-syndication.png" alt="go-syndication/feed" width="720"></p>

# feed

[![CI](https://github.com/go-syndication/feed/actions/workflows/ci.yml/badge.svg)](https://github.com/go-syndication/feed/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-syndication/feed.svg)](https://pkg.go.dev/github.com/go-syndication/feed)
[![License: BSD-3-Clause](https://img.shields.io/badge/License-BSD--3--Clause-blue.svg)](LICENSE)

Pure-Go RSS 2.0 / Atom 1.0 / JSON Feed 1.1 parser and fetcher. `CGO_ENABLED=0`, zero third-party dependencies — standard library only.

The package auto-detects the source format from the bytes and normalizes RSS,
Atom and JSON Feed into a single `Feed` type, so downstream code never has to
branch on the wire format.

## Install

```sh
go get github.com/go-syndication/feed
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/go-syndication/feed"
)

func main() {
	// Parse bytes you already have (format is auto-detected):
	f, err := feed.Parse(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Title)
	for _, e := range f.Entries {
		fmt.Printf("%s — %s (%s)\n", e.Title, e.Link, e.Published)
	}

	// Or fetch over HTTP (nil client uses http.DefaultClient):
	f, err = feed.Fetch(context.Background(), nil, "https://example.com/feed.xml")
	if err != nil {
		panic(err)
	}
}
```

## API

```go
func Parse(data []byte) (*Feed, error)
func Fetch(ctx context.Context, client *http.Client, url string) (*Feed, error)

type Feed struct {
	Title   string
	Link    string
	Entries []Entry
}

type Entry struct {
	ID        string
	Title     string
	Author    string
	Summary   string
	Content   string
	Link      string
	Published time.Time
	Media     []Enclosure
}

type Enclosure struct {
	URL  string
	Type string
}
```

## Supported formats

- **RSS 2.0** — including `dc:creator`, `content:encoded`, `<enclosure>` and `media:content`.
- **Atom 1.0** — including `rel="alternate"` and `rel="enclosure"` links.
- **JSON Feed 1.1** — including `author`/`authors`, `content_html`/`content_text` and `attachments`.

## License

BSD-3-Clause. See [LICENSE](LICENSE).
