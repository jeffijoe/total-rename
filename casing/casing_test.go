package casing_test

import (
	"reflect"
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
	test("use javascript", []string{"USE_JAVASCRIPT", "USE JAVASCRIPT", "use_javascript", "useJavascript", "UseJavascript"})
}

func TestVariants_GetVariant(t *testing.T) {
	type args struct {
		casing casing.Casing
	}
	tests := []struct {
		name     string
		variants casing.Variants
		args     args
		want     casing.Variant
	}{
		{
			name:     "case 1",
			variants: casing.GenerateCasings("space"),
			args: args{
				casing: casing.UpperCase,
			},
			want: casing.Variant{
				Value:  "SPACE",
				Casing: casing.UpperCase,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.variants.GetVariant(tt.args.casing); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variants.GetVariant() = %v, want %v", got, tt.want)
			}
		})
	}
}
