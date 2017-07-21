package casing_test

import (
	"testing"

	"github.com/jeffijoe/total-rename/casing"
	"github.com/stretchr/testify/assert"
)

func TestDetermineCasing(t *testing.T) {
	assert.EqualValues(t, casing.LowerCase, casing.DetermineCasing("hello"))
	assert.EqualValues(t, casing.UpperCase, casing.DetermineCasing("HELLO"))
	assert.EqualValues(t, casing.CamelCase, casing.DetermineCasing("helloThere"))
	assert.EqualValues(t, casing.TitleCase, casing.DetermineCasing("HelloThere"))
	assert.EqualValues(t, casing.TitleCase, casing.DetermineCasing("Hello"))
	assert.EqualValues(t, casing.SnakeCase, casing.DetermineCasing("hello_there"))
}

func TestGenerateCasings(t *testing.T) {
	test := func(input string, expectedStrings []string) {
		result := casing.GenerateCasings(input)
		resultValues := []string{}
		for _, v := range result {
			resultValues = append(resultValues, v.Value)
		}
		for _, expected := range expectedStrings {
			assert.Contains(t, resultValues, expected)
		}
	}

	test("space stuff", []string{"SPACE_STUFF", "SPACE STUFF", "space_stuff", "spaceStuff", "SpaceStuff"})
}
