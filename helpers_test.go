package feed

import "testing"

func TestParseDate(t *testing.T) {
	cases := map[string]bool{ // input -> expect zero
		"":                                true,
		"not a date":                      true,
		"Mon, 02 Jan 2006 15:04:05 -0700": false,
		"Mon, 02 Jan 2006 15:04:05 GMT":   false,
		"2006-01-02T15:04:05Z":            false,
		"2006-01-02T15:04:05.123456789Z":  false,
		"02 Jan 2006 15:04:05 -0700":      false,
		"02 Jan 2006 15:04:05 GMT":        false,
		"2006-01-02 15:04:05":             false,
		"2006-01-02":                      false,
	}
	for in, wantZero := range cases {
		if got := parseDate(in).IsZero(); got != wantZero {
			t.Errorf("parseDate(%q).IsZero() = %v, want %v", in, got, wantZero)
		}
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := firstNonEmpty("", "b", "c"); got != "b" {
		t.Errorf("got %q", got)
	}
	if got := firstNonEmpty("", ""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestAtomAlternate(t *testing.T) {
	cases := []struct {
		name  string
		links []atomLink
		want  string
	}{
		{"alternate", []atomLink{{Href: "a", Rel: "self"}, {Href: "b", Rel: "alternate"}}, "b"},
		{"empty-rel", []atomLink{{Href: "x", Rel: "self"}, {Href: "y", Rel: ""}}, "y"},
		{"only-self", []atomLink{{Href: "s", Rel: "self"}}, "s"},
		{"skip-empty-href", []atomLink{{Href: "", Rel: "alternate"}, {Href: "z", Rel: "self"}}, "z"},
		{"two-empty-rel", []atomLink{{Href: "p", Rel: ""}, {Href: "q", Rel: ""}}, "p"},
		{"none", nil, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := atomAlternate(c.links); got != c.want {
				t.Errorf("atomAlternate = %q, want %q", got, c.want)
			}
		})
	}
}
