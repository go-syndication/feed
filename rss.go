package feed

import "encoding/xml"

// rssDocument models the subset of RSS 2.0 the parser understands, including
// the dc:, content: and media: namespace extensions.
type rssDocument struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Link  string    `xml:"link"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	GUID        string         `xml:"guid"`
	Description string         `xml:"description"`
	Author      string         `xml:"author"`
	Creator     string         `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Encoded     string         `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	PubDate     string         `xml:"pubDate"`
	Date        string         `xml:"http://purl.org/dc/elements/1.1/ date"`
	Enclosures  []rssEnclosure `xml:"enclosure"`
	MediaGroups []rssMedia     `xml:"http://search.yahoo.com/mrss/ content"`
}

type rssEnclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type rssMedia struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

func parseRSS(data []byte) (*Feed, error) {
	var doc rssDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	f := &Feed{
		Title: doc.Channel.Title,
		Link:  doc.Channel.Link,
	}
	for _, it := range doc.Channel.Items {
		e := Entry{
			Title:     it.Title,
			Link:      it.Link,
			Summary:   it.Description,
			Content:   it.Encoded,
			Published: parseDate(firstNonEmpty(it.PubDate, it.Date)),
		}
		if it.GUID != "" {
			e.ID = it.GUID
		} else {
			e.ID = it.Link
		}
		e.Author = firstNonEmpty(it.Creator, it.Author)
		for _, enc := range it.Enclosures {
			if enc.URL != "" {
				e.Media = append(e.Media, Enclosure{URL: enc.URL, Type: enc.Type})
			}
		}
		for _, m := range it.MediaGroups {
			if m.URL != "" {
				e.Media = append(e.Media, Enclosure{URL: m.URL, Type: m.Type})
			}
		}
		f.Entries = append(f.Entries, e)
	}
	return f, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
