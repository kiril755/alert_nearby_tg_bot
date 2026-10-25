package matcher

import "strings"

// Matcher reports whether a message is relevant, based on a keyword list.
type Matcher struct {
	keywords []string
}

// New creates a Matcher. Keywords are expected to be lowercased already.
func New(keywords []string) *Matcher {
	return &Matcher{keywords: keywords}
}

// Match reports whether text contains any configured keyword (case-insensitive).
func (m *Matcher) Match(text string) bool {
	lower := strings.ToLower(text)
	for _, kw := range m.keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}