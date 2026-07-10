package feed

import "encoding/json"

// jsonFeed models the subset of JSON Feed 1.1 the parser understands.
type jsonFeed struct {
	Title string     `json:"title"`
	Home  string     `json:"home_page_url"`
	Items []jsonItem `json:"items"`
}

type jsonItem struct {
	ID            string           `json:"id"`
	Title         string           `json:"title"`
	URL           string           `json:"url"`
	Summary       string           `json:"summary"`
	ContentHTML   string           `json:"content_html"`
	ContentText   string           `json:"content_text"`
	DatePublished string           `json:"date_published"`
	Author        *jsonAuthor      `json:"author"`
	Authors       []jsonAuthor     `json:"authors"`
	Attachments   []jsonAttachment `json:"attachments"`
}

type jsonAuthor struct {
	Name string `json:"name"`
}

type jsonAttachment struct {
	URL      string `json:"url"`
	MIMEType string `json:"mime_type"`
}

func parseJSONFeed(data []byte) (*Feed, error) {
	var doc jsonFeed
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	f := &Feed{
		Title: doc.Title,
		Link:  doc.Home,
	}
	for _, it := range doc.Items {
		e := Entry{
			ID:        it.ID,
			Title:     it.Title,
			Link:      it.URL,
			Summary:   it.Summary,
			Content:   firstNonEmpty(it.ContentHTML, it.ContentText),
			Author:    jsonItemAuthor(it),
			Published: parseDate(it.DatePublished),
		}
		for _, a := range it.Attachments {
			if a.URL != "" {
				e.Media = append(e.Media, Enclosure{URL: a.URL, Type: a.MIMEType})
			}
		}
		f.Entries = append(f.Entries, e)
	}
	return f, nil
}

// jsonItemAuthor resolves the author name from the singular author object or
// the first element of the authors array (JSON Feed 1.1).
func jsonItemAuthor(it jsonItem) string {
	if it.Author != nil && it.Author.Name != "" {
		return it.Author.Name
	}
	if len(it.Authors) > 0 {
		return it.Authors[0].Name
	}
	return ""
}
