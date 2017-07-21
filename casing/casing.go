package casing

import (
	"strings"

	"github.com/mgutz/str"
)

// Casing type
type Casing uint8

// Types of casing
const (
	Original       = iota
	LowerCase      = iota
	UpperCase      = iota
	CamelCase      = iota
	TitleCase      = iota
	SnakeCase      = iota
	KebabCase      = iota
	UpperSnakeCase = iota
	UpperKebabCase = iota
)

// Variants contains variations of a string in different casings.
type Variants []Variant

// Variant represents a string in a different casing variant.
type Variant struct {
	Casing Casing
	Value  string
}

// DetermineCasing determines a strings casing
func DetermineCasing(s string) Casing {
	matches := []Casing{Original}
	hasUnderscore := strings.Contains(s, "_")
	hasDash := strings.Contains(s, "-")
	if s == str.Camelize(s) {
		matches = append(matches, CamelCase)
	}
	if s == strings.ToLower(s) {
		matches = append(matches, LowerCase)
	}
	if s == str.Classify(s) {
		matches = append(matches, TitleCase)
	}
	if s == strings.ToUpper(s) {
		matches = append(matches, UpperCase)
	}
	if s == strings.ToUpper(str.Underscore(s)) && hasUnderscore {
		matches = append(matches, UpperSnakeCase)
	}
	if s == strings.ToUpper(str.Dasherize(s)) && hasDash {
		matches = append(matches, UpperKebabCase)
	}
	if s == str.Underscore(s) && hasUnderscore {
		matches = append(matches, SnakeCase)
	}
	if s == str.Dasherize(s) && hasDash {
		matches = append(matches, KebabCase)
	}

	// There are going to be multiple matches,
	// the above is ordered in a way that the
	// last match is the most specific.
	return matches[len(matches)-1]
}

// GenerateCasings generates casings for the specified string
func GenerateCasings(s string) Variants {
	underscored := strings.Trim(str.Underscore(s), "_")
	dasherized := strings.Trim(str.Dasherize(s), "-")
	return Variants{
		Variant{Original, s},
		Variant{LowerCase, strings.ToLower(s)},
		Variant{UpperCase, strings.ToUpper(s)},
		Variant{CamelCase, str.Camelize(s)},
		Variant{TitleCase, str.Classify(s)},
		Variant{SnakeCase, underscored},
		Variant{KebabCase, dasherized},
		Variant{UpperSnakeCase, strings.ToUpper(underscored)},
		Variant{UpperKebabCase, strings.ToLower(dasherized)},
	}
}

// GetVariant returns the variant for the specified casing.
func (variants Variants) GetVariant(casing Casing) Variant {
	var orig Variant
	for _, v := range variants {
		if v.Casing == casing {
			return v
		}
		if v.Casing == Original {
			orig = v
		}
	}
	return orig
}
