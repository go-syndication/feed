package feed

import (
	"time"

	"github.com/go-datetime/dates"
)

// parseDate parses a feed date — RSS RFC 822/1123 variants, Atom and JSON Feed
// RFC 3339, and bare `YYYY-MM-DD` dates — returning the zero time.Time for
// empty or unparseable input so a single malformed date never fails an entire
// feed parse. The real-world format zoo (including named-zone abbreviations and
// 2-digit years) lives in the shared go-datetime/dates library so every fleet
// consumer reuses one table instead of maintaining its own.
func parseDate(s string) time.Time {
	t, _ := dates.Parse(s)
	return t
}
