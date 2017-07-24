package simplematch

import (
	"strings"
)

// Matcher is a matcher.
type Matcher struct {
	splat []string
	empty bool
}

// NewMatcher returns a new matcher.
func NewMatcher(pattern string) *Matcher {
	splat := strings.Split(pattern, "|")
	for i, s := range splat {
		splat[i] = strings.ToUpper(s)
	}
	m := &Matcher{
		splat: splat,
		empty: len(splat) == 0 || (len(splat) == 1 && splat[0] == ""),
	}
	return m
}

// Matches checks if a string matches the pattern.
func (m *Matcher) Matches(s string) bool {
	s = strings.ToUpper(s)
	if m.empty {
		return false
	}

	for _, p := range m.splat {
		if strings.Contains(s, p) {
			return true
		}
	}
	return false
}
