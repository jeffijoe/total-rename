package replacer

import (
	"fmt"
	"testing"

	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"strings"

	"github.com/jeffijoe/total-rename/casing"
	"github.com/jeffijoe/total-rename/lister"
	"github.com/jeffijoe/total-rename/scanner"
	"github.com/jeffijoe/total-rename/util"
	"github.com/stretchr/testify/assert"
)

func TestReplaceText(t *testing.T) {
	type args struct {
		source              string
		occurences          scanner.Occurences
		replacementVariants casing.Variants
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			want: "board is great, Boards Are Great, BOARDMEMBERS SUCK! board_snakes are the worst, but BOARD_UPPER_SNAKES SUCK EVEN MORE!",
			args: args{
				source: "space is great, Spaces Are Great, SPACEMEMBERS SUCK! space_snakes are the worst, but SPACE_UPPER_SNAKES SUCK EVEN MORE!",
				occurences: scanner.Occurences{
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 0,
					},
					&scanner.Occurence{
						Casing:     casing.TitleCase,
						Match:      "Space",
						StartIndex: 16,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "SPACE",
						StartIndex: 34,
					},
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 53,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "space",
						StartIndex: 85,
					},
				},
				replacementVariants: casing.GenerateCasings("board"),
			},
		},
		{
			name: "case 2",
			want: "board is great, Boards Are Great, BOARDMEMBERS SUCK! board_snakes are the worst, but BOARD",
			args: args{
				source: "space is great, Spaces Are Great, SPACEMEMBERS SUCK! space_snakes are the worst, but SPACE",
				occurences: scanner.Occurences{
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 0,
					},
					&scanner.Occurence{
						Casing:     casing.TitleCase,
						Match:      "Space",
						StartIndex: 16,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "SPACE",
						StartIndex: 34,
					},
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 53,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "space",
						StartIndex: 85,
					},
				},
				replacementVariants: casing.GenerateCasings("board"),
			},
		},
		{
			name: "case 3",
			want: "the board is great, Boards Are Great, BOARDMEMBERS SUCK! board_snakes are the worst, but BOARD",
			args: args{
				source: "the space is great, Spaces Are Great, SPACEMEMBERS SUCK! space_snakes are the worst, but SPACE",
				occurences: scanner.Occurences{
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 4,
					},
					&scanner.Occurence{
						Casing:     casing.TitleCase,
						Match:      "Space",
						StartIndex: 20,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "SPACE",
						StartIndex: 38,
					},
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "space",
						StartIndex: 57,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "space",
						StartIndex: 89,
					},
				},
				replacementVariants: casing.GenerateCasings("board"),
			},
		},
		{
			name: "case 4",
			want: "timeSpace is time_space with TIME_SPACE for TimeSpace and TIMESPACE with timespace",
			args: args{
				source: "spaceTime is space_time with SPACE_TIME for SpaceTime and SPACETIME with spacetime",
				occurences: scanner.Occurences{
					&scanner.Occurence{
						Casing:     casing.Original,
						Match:      "spaceTime",
						StartIndex: 0,
					},
					&scanner.Occurence{
						Casing:     casing.SnakeCase,
						Match:      "space_time",
						StartIndex: 13,
					},
					&scanner.Occurence{
						Casing:     casing.UpperSnakeCase,
						Match:      "SPACE_TIME",
						StartIndex: 29,
					},
					&scanner.Occurence{
						Casing:     casing.TitleCase,
						Match:      "SpaceTime",
						StartIndex: 44,
					},
					&scanner.Occurence{
						Casing:     casing.UpperCase,
						Match:      "SPACETIME",
						StartIndex: 58,
					},
					&scanner.Occurence{
						Casing:     casing.LowerCase,
						Match:      "spacetime",
						StartIndex: 73,
					},
				},
				replacementVariants: casing.GenerateCasings("timeSpace"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceText(tt.args.source, tt.args.occurences, tt.args.replacementVariants); got != tt.want {
				t.Errorf("ReplaceText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceFileContent(t *testing.T) {
	now := time.Now().UTC().Unix()
	file := filepath.Join(os.TempDir(), "total-rename-test-"+string(now)+".txt")
	ioutil.WriteFile(file, []byte("plz"), 0644)
	content, _ := ioutil.ReadFile(file)
	assert.Equal(t, "plz", string(content))

	ReplaceFileContent(file, "haha")
	content, _ = ioutil.ReadFile(file)
	assert.Equal(t, "haha", string(content))
}

func TestTotalRename(t *testing.T) {
	tempDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("%d", time.Now().Unix()),
	)
	fixtureInputDir, _ := filepath.Abs(filepath.Join(util.GetWD(), "../_fixtures/fixture1/input"))
	util.CopyDir(fixtureInputDir, tempDir)
	nodes, _ := lister.ListFileNodes(
		tempDir,
		"**/*.js",
	)

	groups, _ := scanner.ScanFileNodes(nodes, "space")

	_, err := TotalRename(groups, os.Rename, ReplaceFileContent)
	assert.NoError(t, err)
	expectedDir, _ := filepath.Abs(filepath.Join(util.GetWD(), "../_fixtures/fixture1/expected"))
	expectedNodes, _ := lister.ListFileNodes(
		expectedDir,
		"**/*.js",
	)

	for _, node := range expectedNodes {
		tmpPath := filepath.Join(
			tempDir,
			strings.TrimPrefix(
				node.Path,
				strings.TrimSuffix(
					fixtureInputDir,
					"/input",
				)+"/expected",
			),
		)
		t.Log(node.Path)
		t.Log(fixtureInputDir)
		t.Log(tmpPath)
		fi, err := os.Stat(tmpPath)
		assert.NoError(t, err)
		if node.Type == lister.NodeTypeDir {
			assert.True(t, fi.IsDir())
		} else {
			assert.False(t, fi.IsDir())
			expectedContent, _ := ioutil.ReadFile(node.Path)
			actualContent, _ := ioutil.ReadFile(tmpPath)
			assert.Equal(t, expectedContent, actualContent)
		}
	}
}
