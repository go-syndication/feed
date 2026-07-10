package feed

import (
	"errors"
	"testing"
	"time"
)

const rssSample = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0"
     xmlns:dc="http://purl.org/dc/elements/1.1/"
     xmlns:content="http://purl.org/rss/1.0/modules/content/"
     xmlns:media="http://search.yahoo.com/mrss/">
  <channel>
    <title>Example Channel</title>
    <link>https://example.com/</link>
    <item>
      <title>First Post</title>
      <link>https://example.com/first</link>
      <guid>urn:uuid:first</guid>
      <description>Short summary</description>
      <content:encoded><![CDATA[<p>Full body</p>]]></content:encoded>
      <dc:creator>Ada Lovelace</dc:creator>
      <pubDate>Mon, 02 Jan 2006 15:04:05 -0700</pubDate>
      <enclosure url="https://example.com/a.mp3" type="audio/mpeg"/>
      <enclosure url="" type="audio/mpeg"/>
      <media:content url="https://example.com/img.jpg" type="image/jpeg"/>
      <media:content url="" type="image/jpeg"/>
    </item>
    <item>
      <title>Second Post</title>
      <link>https://example.com/second</link>
      <description>Another</description>
      <author>fallback@example.com</author>
      <pubDate>Tue, 03 Jan 2006 15:04:05 GMT</pubDate>
    </item>
  </channel>
</rss>`

func TestParseRSS(t *testing.T) {
	f, err := Parse([]byte(rssSample))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.Title != "Example Channel" || f.Link != "https://example.com/" {
		t.Fatalf("channel mismatch: %+v", f)
	}
	if len(f.Entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(f.Entries))
	}
	e := f.Entries[0]
	if e.ID != "urn:uuid:first" {
		t.Errorf("ID = %q", e.ID)
	}
	if e.Author != "Ada Lovelace" {
		t.Errorf("Author = %q", e.Author)
	}
	if e.Summary != "Short summary" {
		t.Errorf("Summary = %q", e.Summary)
	}
	if e.Content != "<p>Full body</p>" {
		t.Errorf("Content = %q", e.Content)
	}
	if e.Published.IsZero() {
		t.Errorf("Published is zero")
	}
	if len(e.Media) != 2 {
		t.Fatalf("want 2 media, got %d: %+v", len(e.Media), e.Media)
	}
	if e.Media[0].URL != "https://example.com/a.mp3" || e.Media[0].Type != "audio/mpeg" {
		t.Errorf("media0 = %+v", e.Media[0])
	}
	if e.Media[1].URL != "https://example.com/img.jpg" {
		t.Errorf("media1 = %+v", e.Media[1])
	}
	// Second entry: no guid -> ID falls back to link; author fallback.
	e2 := f.Entries[1]
	if e2.ID != "https://example.com/second" {
		t.Errorf("e2.ID = %q", e2.ID)
	}
	if e2.Author != "fallback@example.com" {
		t.Errorf("e2.Author = %q", e2.Author)
	}
}

func TestParseRSSMalformed(t *testing.T) {
	// Detects as RSS, then fails to unmarshal (unclosed tags).
	_, err := Parse([]byte(`<rss><channel><title>x</title>`))
	if err == nil {
		t.Fatal("expected error for malformed RSS")
	}
}

const atomSample = `<?xml version="1.0" encoding="utf-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Atom Example</title>
  <link rel="alternate" href="https://atom.example.com/"/>
  <entry>
    <id>tag:example.com,2006:1</id>
    <title>Atom Entry</title>
    <author><name>Grace Hopper</name></author>
    <summary>Summary text</summary>
    <content>Content body</content>
    <published>2006-01-02T15:04:05Z</published>
    <updated>2007-01-02T15:04:05Z</updated>
    <link rel="alternate" href="https://atom.example.com/entry"/>
    <link rel="enclosure" href="https://atom.example.com/f.pdf" type="application/pdf"/>
    <link rel="enclosure" href=""/>
  </entry>
</feed>`

func TestParseAtom(t *testing.T) {
	f, err := Parse([]byte(atomSample))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.Title != "Atom Example" || f.Link != "https://atom.example.com/" {
		t.Fatalf("feed mismatch: %+v", f)
	}
	if len(f.Entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(f.Entries))
	}
	e := f.Entries[0]
	if e.ID != "tag:example.com,2006:1" {
		t.Errorf("ID = %q", e.ID)
	}
	if e.Author != "Grace Hopper" {
		t.Errorf("Author = %q", e.Author)
	}
	if e.Link != "https://atom.example.com/entry" {
		t.Errorf("Link = %q", e.Link)
	}
	if e.Content != "Content body" {
		t.Errorf("Content = %q", e.Content)
	}
	if want := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC); !e.Published.Equal(want) {
		t.Errorf("Published = %v, want %v", e.Published, want)
	}
	if len(e.Media) != 1 || e.Media[0].URL != "https://atom.example.com/f.pdf" {
		t.Errorf("Media = %+v", e.Media)
	}
}

func TestParseAtomMalformed(t *testing.T) {
	_, err := Parse([]byte(`<feed><entry><title>x</title>`))
	if err == nil {
		t.Fatal("expected error for malformed Atom")
	}
}

const jsonSample = `{
  "version": "https://jsonfeed.org/version/1.1",
  "title": "JSON Example",
  "home_page_url": "https://json.example.com/",
  "items": [
    {
      "id": "1",
      "title": "JSON Entry",
      "url": "https://json.example.com/1",
      "summary": "A summary",
      "content_html": "<p>hi</p>",
      "date_published": "2010-02-07T14:04:00-05:00",
      "author": {"name": "Alan Turing"},
      "attachments": [
        {"url": "https://json.example.com/a.mp4", "mime_type": "video/mp4"},
        {"url": "", "mime_type": "video/mp4"}
      ]
    },
    {
      "id": "2",
      "title": "Text Entry",
      "content_text": "plain text",
      "authors": [{"name": "Katherine Johnson"}]
    },
    {
      "id": "3",
      "title": "No Author"
    }
  ]
}`

func TestParseJSONFeed(t *testing.T) {
	f, err := Parse([]byte(jsonSample))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.Title != "JSON Example" || f.Link != "https://json.example.com/" {
		t.Fatalf("feed mismatch: %+v", f)
	}
	if len(f.Entries) != 3 {
		t.Fatalf("want 3 entries, got %d", len(f.Entries))
	}
	e := f.Entries[0]
	if e.Author != "Alan Turing" {
		t.Errorf("Author = %q", e.Author)
	}
	if e.Content != "<p>hi</p>" {
		t.Errorf("Content = %q", e.Content)
	}
	if e.Published.IsZero() {
		t.Errorf("Published is zero")
	}
	if len(e.Media) != 1 || e.Media[0].Type != "video/mp4" {
		t.Errorf("Media = %+v", e.Media)
	}
	// Second: content_text fallback + authors[0].
	if f.Entries[1].Content != "plain text" {
		t.Errorf("e2 Content = %q", f.Entries[1].Content)
	}
	if f.Entries[1].Author != "Katherine Johnson" {
		t.Errorf("e2 Author = %q", f.Entries[1].Author)
	}
	// Third: no author at all.
	if f.Entries[2].Author != "" {
		t.Errorf("e3 Author = %q", f.Entries[2].Author)
	}
}

func TestParseJSONFeedAuthorEmptyName(t *testing.T) {
	// author present but empty name -> fall through to authors array.
	data := `{"title":"t","items":[{"id":"1","author":{"name":""},"authors":[{"name":"Fallback"}]}]}`
	f, err := Parse([]byte(data))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if f.Entries[0].Author != "Fallback" {
		t.Errorf("Author = %q", f.Entries[0].Author)
	}
}

func TestParseJSONMalformed(t *testing.T) {
	_, err := Parse([]byte(`{ this is not json }`))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestParseUnrecognized(t *testing.T) {
	cases := map[string][]byte{
		"empty":        []byte("   \n\t"),
		"html":         []byte(`<html><body>nope</body></html>`),
		"comment":      []byte(`<!-- just a comment -->`),
		"malformedxml": []byte(`< not valid xml`),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Parse(data)
			if err == nil {
				t.Fatal("expected error")
			}
			if name != "malformedxml" && !errors.Is(err, ErrUnrecognized) {
				t.Fatalf("want ErrUnrecognized, got %v", err)
			}
		})
	}
}
