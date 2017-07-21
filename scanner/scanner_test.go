package scanner_test

import (
	"path/filepath"
	"testing"

	"github.com/jeffijoe/total-rename/casing"
	"github.com/jeffijoe/total-rename/lister"
	"github.com/jeffijoe/total-rename/scanner"
	"github.com/jeffijoe/total-rename/util"
	"github.com/stretchr/testify/assert"
)

func TestScanFile(t *testing.T) {
	test := func(file string, expectedOccurences []scanner.Occurence) {
		occurences, err := scanner.ScanFile(
			filepath.Join(util.GetWD(), "../_fixtures", file),
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

func TestScanFilePath(t *testing.T) {
	type args struct {
		filePath string
		variants casing.Variants
	}
	tests := []struct {
		name string
		args args
		want scanner.Occurences
	}{
		{
			name: "case 1",
			args: args{filePath: "/test/space-stuff/Space.js", variants: casing.GenerateCasings("space")},
			want: scanner.Occurences{
				&scanner.Occurence{
					Match:      "space",
					Casing:     casing.Original,
					StartIndex: 6,
				},
				&scanner.Occurence{
					Match:      "Space",
					Casing:     casing.TitleCase,
					StartIndex: 18,
				},
			},
		},
		{
			name: "case 2",
			args: args{filePath: "/test/api/repositories/spaces/SpaceRepository.js", variants: casing.GenerateCasings("space")},
			want: scanner.Occurences{
				&scanner.Occurence{
					Match:      "space",
					Casing:     casing.Original,
					StartIndex: 23,
				},
				&scanner.Occurence{
					Match:      "Space",
					Casing:     casing.TitleCase,
					StartIndex: 30,
				},
			},
		},
		{
			name: "case 3",
			args: args{filePath: "/test/api/consts/SPACE_TYPE.js", variants: casing.GenerateCasings("Space")},
			want: scanner.Occurences{
				&scanner.Occurence{
					Match:      "SPACE",
					Casing:     casing.UpperCase,
					StartIndex: 17,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := scanner.ScanFilePath(tt.args.filePath, tt.args.variants)
			for i, got := range res {
				expected := tt.want[i]
				assert.Equal(t, expected.Match, got.Match)
				assert.Equal(t, expected.Casing, got.Casing)
				assert.Equal(t, expected.StartIndex, got.StartIndex)
			}
		})
	}
}

func TestScanFileNodes(t *testing.T) {
	nodes := lister.FileNodes{
		&lister.FileNode{
			Path: filepath.Join(util.GetWD(), "../_fixtures/fixture1/input/space-repository.js"),
			Type: lister.NodeTypeFile,
		},
		&lister.FileNode{
			Path: filepath.Join(util.GetWD(), "../_fixtures/fixture1/input/spaces"),
			Type: lister.NodeTypeDir,
		},
	}

	expectedGroups := scanner.OccurenceGroups{
		&scanner.OccurenceGroup{
			Occurences: scanner.Occurences{
				&scanner.Occurence{Casing: casing.UpperCase, Match: "SPACE", LineNumber: 0},
				&scanner.Occurence{Casing: casing.Original, Match: "space", LineNumber: 0},
				&scanner.Occurence{Casing: casing.TitleCase, Match: "Space", LineNumber: 2},
				&scanner.Occurence{Casing: casing.TitleCase, Match: "Space", LineNumber: 3},
				&scanner.Occurence{Casing: casing.TitleCase, Match: "Space", LineNumber: 6},
			},
		},
		&scanner.OccurenceGroup{
			Occurences: scanner.Occurences{
				&scanner.Occurence{Casing: casing.Original, Match: "space", LineNumber: 0},
			},
		},
	}
	result, err := scanner.ScanFileNodes(nodes, "space")
	assert.NoError(t, err)
	for i, group := range result {
		exGroup := expectedGroups[i]
		for j, got := range group.Occurences {
			want := exGroup.Occurences[j]
			assert.Equal(t, want.Casing, got.Casing)
			assert.Equal(t, want.Match, got.Match)
		}
	}
}

func TestScanFileNodes_Error(t *testing.T) {
	nodes := lister.FileNodes{
		&lister.FileNode{
			Path: filepath.Join(util.GetWD(), "../_fixtures/fixture1/input/doesnotexist.js"),
			Type: lister.NodeTypeFile,
		},
		&lister.FileNode{
			Path: filepath.Join(util.GetWD(), "../_fixtures/fixture1/input/spaces"),
			Type: lister.NodeTypeDir,
		},
	}

	_, err := scanner.ScanFileNodes(nodes, "space")
	assert.Error(t, err)
}
