package scanner_test

import (
	"testing"

	"path/filepath"

	"github.com/jeffijoe/total-replace/casing"
	"github.com/jeffijoe/total-replace/scanner"
	"github.com/jeffijoe/total-replace/util"
	"github.com/stretchr/testify/assert"
)

func TestScanFile(t *testing.T) {
	test := func(file string, expectedOccurences []scanner.Occurence) {
		occurences, err := scanner.ScanFile(
			filepath.Join(util.GetWD(), "_fixtures", file),
			casing.GenerateCasings("space"),
		)
		assert.NoError(t, err)
		for idx, expected := range expectedOccurences {
			actual := occurences[idx]
			assert.Equal(t, expected.Match, actual.Match, "String did not match")
			assert.Equal(t, expected.Casing, actual.Casing, "Casing did not match for "+actual.Match)
			assert.Equal(t, expected.StartIndex, actual.StartIndex, "Start index did not match")
			assert.Equal(t, expected.LineNumber, actual.LineNumber, "Line number did not match")
		}
	}

	test("fixture1/input/space-repository.js", []scanner.Occurence{
		scanner.Occurence{Casing: casing.UpperCase, Match: "SPACE", StartIndex: 7, LineNumber: 0},
		scanner.Occurence{Casing: casing.Original, Match: "space", StartIndex: 31, LineNumber: 0},
		scanner.Occurence{Casing: casing.TitleCase, Match: "Space", StartIndex: 66, LineNumber: 2},
		scanner.Occurence{Casing: casing.TitleCase, Match: "Space", StartIndex: 106, LineNumber: 3},
		scanner.Occurence{Casing: casing.TitleCase, Match: "Space", StartIndex: 133, LineNumber: 6},
	})
}

func TestGetSurroundingLines(t *testing.T) {
	src := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	before, after := scanner.GetSurroundingLines(
		src,
		0,
		2,
	)
	assert.Equal(t, 0, len(before))
	assert.Equal(t, 2, len(after))
	assert.Equal(t, "2", after[0])
	assert.Equal(t, "3", after[1])
}
