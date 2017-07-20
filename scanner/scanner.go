package scanner

import (
	"io/ioutil"
	"sort"
	"strings"

	"fmt"

	"github.com/jeffijoe/total-replace/casing"
	"github.com/mgutz/str"
)

// Occurences is a slice of occurences.
type Occurences []*Occurence

// Occurence is an occurence of the search text in a file.
type Occurence struct {
	Casing                 casing.Casing
	Match                  string
	StartIndex             int
	Path                   string
	SurroundingLinesBefore []string
	SurroundingLinesAfter  []string
	LineNumber             int
}

// ScanFile scans a single file and returns the occurences of the
// specified variants.
func ScanFile(filePath string, variants casing.Variants) (Occurences, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")
	result := Occurences{}
	totalIndex := 0
	for lineIdx, line := range lines {
		used := map[int]struct{}{}
		for _, variant := range variants {
			lineOccurences := getOccurences(line, variant.Value)
			if len(lineOccurences) == 0 {
				continue
			}

			for _, startIndex := range lineOccurences {
				if _, ok := used[startIndex]; ok {
					continue
				}

				linesBefore, linesAfter := GetSurroundingLines(lines, lineIdx, 3)
				occurence := &Occurence{
					Casing:                 variant.Casing,
					Match:                  variant.Value,
					Path:                   filePath,
					StartIndex:             totalIndex + startIndex,
					SurroundingLinesBefore: linesBefore,
					SurroundingLinesAfter:  linesAfter,
					LineNumber:             lineIdx,
				}
				used[startIndex] = struct{}{}
				result = append(result, occurence)
			}
		}
		totalIndex = totalIndex + len(line) + 1
	}
	sort.Sort(result)
	return result, nil
}

// GetSurroundingLines returns the surrounding lines
func GetSurroundingLines(lines []string, lineIdx int, count int) (before []string, after []string) {
	length := len(lines)
	before = []string{}
	after = []string{}

	for i := lineIdx - 1; ; i-- {
		if i < 0 {
			break
		}
		before = append(before, lines[i])
		if len(before) == count {
			break
		}
	}

	for i := lineIdx + 1; ; i++ {
		if i >= length {
			break
		}
		after = append(after, lines[i])
		if len(after) == count {
			break
		}
	}
	return before, after
}

// Returns a slice of index occurences
func getOccurences(s string, needle string) []int {
	buf := []int{}
	last := 0
	for {
		index := str.IndexOf(s, needle, last)
		if index == -1 {
			return buf
		}
		buf = append(buf, index)
		last = index + 1
	}
}

func (o Occurence) String() string {
	return fmt.Sprintf("{ StartIndex: %d, Match: %s, Casing: %d }", o.StartIndex, o.Match, o.Casing)
}

func (slice Occurences) Len() int {
	return len(slice)
}

func (slice Occurences) Less(i, j int) bool {
	return slice[i].StartIndex < slice[j].StartIndex
}

func (slice Occurences) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
