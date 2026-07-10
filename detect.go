package feed

import (
	"bytes"
	"encoding/xml"
	"io"
)

type xmlKind int

const (
	xmlUnknown xmlKind = iota
	xmlRSS
	xmlAtom
)

// detectXML inspects the first XML start element to decide whether the
// document is an RSS or Atom feed. It returns an error only for malformed XML
// (a missing root element yields xmlUnknown with a nil error so Parse can
// report ErrUnrecognized).
func detectXML(data []byte) (xmlKind, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return xmlUnknown, nil
		}
		if err != nil {
			return xmlUnknown, err
		}
		if se, ok := tok.(xml.StartElement); ok {
			switch se.Name.Local {
			case "rss":
				return xmlRSS, nil
			case "feed":
				return xmlAtom, nil
			default:
				return xmlUnknown, nil
			}
		}
	}
}
