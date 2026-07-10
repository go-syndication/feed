package feed

import (
	"bytes"
	"errors"
)

// ErrUnrecognized is returned by [Parse] when the input does not look like a
// supported feed format.
var ErrUnrecognized = errors.New("feed: unrecognized format")

// Parse detects the format (RSS/Atom/JSONFeed) from the bytes and parses it.
func Parse(data []byte) (*Feed, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, ErrUnrecognized
	}
	if trimmed[0] == '{' {
		return parseJSONFeed(data)
	}
	kind, err := detectXML(trimmed)
	if err != nil {
		return nil, err
	}
	switch kind {
	case xmlRSS:
		return parseRSS(data)
	case xmlAtom:
		return parseAtom(data)
	default:
		return nil, ErrUnrecognized
	}
}
