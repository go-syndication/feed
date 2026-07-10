package feed

import "encoding/xml"

// atomFeed models the subset of Atom 1.0 the parser understands.
type atomFeed struct {
	XMLName xml.Name    `xml:"feed"`
	Title   string      `xml:"title"`
	Links   []atomLink  `xml:"link"`
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID        string     `xml:"id"`
	Title     string     `xml:"title"`
	Summary   string     `xml:"summary"`
	Content   string     `xml:"content"`
	Updated   string     `xml:"updated"`
	Published string     `xml:"published"`
	Author    atomAuthor `xml:"author"`
	Links     []atomLink `xml:"link"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

func parseAtom(data []byte) (*Feed, error) {
	var doc atomFeed
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	f := &Feed{
		Title: doc.Title,
		Link:  atomAlternate(doc.Links),
	}
	for _, ae := range doc.Entries {
		e := Entry{
			ID:        ae.ID,
			Title:     ae.Title,
			Author:    ae.Author.Name,
			Summary:   ae.Summary,
			Content:   ae.Content,
			Link:      atomAlternate(ae.Links),
			Published: parseDate(firstNonEmpty(ae.Published, ae.Updated)),
		}
		for _, l := range ae.Links {
			if l.Rel == "enclosure" && l.Href != "" {
				e.Media = append(e.Media, Enclosure{URL: l.Href, Type: l.Type})
			}
		}
		f.Entries = append(f.Entries, e)
	}
	return f, nil
}

// atomAlternate returns the href of the link with rel="alternate", falling
// back to the first link with an empty rel, then the first link with any
// non-empty href.
func atomAlternate(links []atomLink) string {
	var firstEmptyRel, firstAny string
	for _, l := range links {
		if l.Href == "" {
			continue
		}
		if l.Rel == "alternate" {
			return l.Href
		}
		if firstAny == "" {
			firstAny = l.Href
		}
		if l.Rel == "" && firstEmptyRel == "" {
			firstEmptyRel = l.Href
		}
	}
	if firstEmptyRel != "" {
		return firstEmptyRel
	}
	return firstAny
}
