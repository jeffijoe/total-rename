package casing

import (
	"strings"

	"github.com/mgutz/str"
)

// Casing type
type Casing uint8

// Types of casing
const (
	NoMatch        = iota
	LowerCase      = iota
	UpperCase      = iota
	CamelCase      = iota
	TitleCase      = iota
	SnakeCase      = iota
	KebabCase      = iota
	UpperSnakeCase = iota
	UpperKebabCase = iota
)

// DetermineCasing determines a strings casing
func DetermineCasing(s string) Casing {
	matches := []Casing{NoMatch}
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

	return matches[len(matches)-1]
}

// GenerateCasings generates casings for the specified string
func GenerateCasings(s string) Variants {
	underscored := str.Underscore(s)
	dasherized := str.Dasherize(s)
	return Variants{
		NoMatch:        s,
		LowerCase:      strings.ToLower(s),
		UpperCase:      strings.ToUpper(s),
		CamelCase:      str.Camelize(s),
		TitleCase:      str.Classify(s),
		SnakeCase:      underscored,
		KebabCase:      dasherized,
		UpperSnakeCase: strings.ToUpper(underscored),
		UpperKebabCase: strings.ToLower(dasherized),
	}
}

// Variants contains variations of a string in different casings.
type Variants map[Casing]string
