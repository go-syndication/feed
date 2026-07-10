// Package feed provides a dependency-free parser and fetcher for RSS 2.0,
// Atom 1.0 and JSON Feed 1.1 documents. It normalizes each source format into
// a single [Feed] representation.
//
// The package uses only the Go standard library and builds with CGO disabled.
package feed

import "time"

// Feed is a normalized feed regardless of source format.
type Feed struct {
	Title   string
	Link    string // site link
	Entries []Entry
}

// Entry is one item/entry/post.
type Entry struct {
	ID        string
	Title     string
	Author    string
	Summary   string    // short text/description (may be HTML)
	Content   string    // full content if present (may be HTML)
	Link      string    // permalink to the entry
	Published time.Time // zero if absent/unparseable
	Media     []Enclosure
}

// Enclosure is an attached media object (RSS enclosure, media:content,
// Atom link rel=enclosure).
type Enclosure struct {
	URL  string
	Type string // MIME type, e.g. "image/jpeg"
}
