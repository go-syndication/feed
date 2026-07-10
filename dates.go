package feed

import (
	"strings"
	"time"
)

// dateLayouts is the ordered set of layouts tried by parseDate. RSS uses
// RFC822/RFC1123 variants; Atom and JSON Feed use RFC3339.
var dateLayouts = []string{
	time.RFC1123Z,
	time.RFC1123,
	time.RFC3339,
	time.RFC3339Nano,
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"02 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05 MST",
	"2006-01-02T15:04:05Z07:00",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// parseDate tries each known layout in turn and returns the zero time.Time
// when the input is empty or matches none of them. It never reports an error
// so that a single malformed date cannot fail an entire feed parse.
func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t
		}
	}
	return time.Time{}
}
