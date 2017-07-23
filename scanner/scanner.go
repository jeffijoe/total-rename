package scanner

import (
	"io/ioutil"
	"sort"
	"strings"
	"sync"

	"fmt"

	"path/filepath"

	"os"

	"github.com/jeffijoe/total-rename/casing"
	"github.com/jeffijoe/total-rename/lister"
	"github.com/mgutz/str"
)

// OccurenceGroupType determines what kind of occurence group it is.
type OccurenceGroupType uint8

// Types of occurence group types.
const (
	OccurenceGroupTypeContent = iota
	OccurenceGroupTypePath    = iota
)

// OccurenceGroups is a list of occurence groups.
// When sorted, files come first.
type OccurenceGroups []*OccurenceGroup

// OccurenceGroup is a grouping of occurences by file path and type.
type OccurenceGroup struct {
	Path       string
	Occurences Occurences
	Type       OccurenceGroupType
}

// Occurences is a slice of occurences.
type Occurences []*Occurence

// Occurence is an occurence of the search text in a file.
type Occurence struct {
	Casing                 casing.Casing
	Match                  string
	Line                   string
	StartIndex             int
	LineStartIndex         int
	SurroundingLinesBefore []string
	SurroundingLinesAfter  []string
	LineNumber             int
}

// ScanFileNodes will scan files and folders for occurences of the specified string.
func ScanFileNodes(nodes lister.FileNodes, needle string) (OccurenceGroups, error) {
	variants := casing.GenerateCasings(needle)
	ch := make(chan *OccurenceGroup, 20)
	che := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(len(nodes))
	for _, node := range nodes {
		n := node
		go func() {
			defer wg.Done()
			if n.Type == lister.NodeTypeFile {
				occurences, err := ScanFile(n.Path, variants)
				if err != nil {
					che <- err
					return
				}
				if len(occurences) > 0 {
					ch <- &OccurenceGroup{
						Path:       n.Path,
						Occurences: occurences,
						Type:       OccurenceGroupTypeContent,
					}
				}
			}
			pathOccurences := ScanFilePath(n.Path, variants)
			if len(pathOccurences) > 0 {
				ch <- &OccurenceGroup{
					Path:       n.Path,
					Occurences: pathOccurences,
					Type:       OccurenceGroupTypePath,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
		close(che)
	}()

	for err := range che {
		return nil, err
	}

	result := OccurenceGroups{}
	for group := range ch {
		result = append(result, group)
	}
	sort.Stable(result)
	return result, nil
}

// ScanFilePath scans a file path name for occurences.
func ScanFilePath(filePath string, variants casing.Variants) Occurences {
	used := map[int]struct{}{}
	result := Occurences{}
	dirLen := len(filepath.Dir(filePath)) + 1
	fileName := filepath.Base(filePath)
	for _, variant := range variants {
		occurences := getOccurences(fileName, variant.Value)
		if len(occurences) == 0 {
			continue
		}
		for _, startIndex := range occurences {
			if _, ok := used[startIndex]; ok {
				continue
			}

			occurence := &Occurence{
				Casing:         variant.Casing,
				Match:          variant.Value,
				StartIndex:     dirLen + startIndex,
				LineStartIndex: dirLen + startIndex,
				Line:           filePath,
			}
			used[startIndex] = struct{}{}
			result = append(result, occurence)
		}
	}
	sort.Sort(result)
	return result
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
					Casing:         variant.Casing,
					Match:          variant.Value,
					StartIndex:     totalIndex + startIndex,
					LineStartIndex: startIndex,
					Line:           line,
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

	for i := lineIdx - count; ; i++ {
		if i < 0 {
			continue
		}
		if i >= lineIdx {
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

func (slice OccurenceGroups) Len() int {
	return len(slice)
}

func (slice OccurenceGroups) Less(i int, j int) bool {
	left := slice[i]
	right := slice[j]
	if left.Type < right.Type {
		return true
	}
	if left.Type > right.Type {
		return false
	}
	leftPathSegmentCount := len(strings.Split(left.Path, string(os.PathSeparator)))
	rightPathSegmentCount := len(strings.Split(right.Path, string(os.PathSeparator)))
	return leftPathSegmentCount > rightPathSegmentCount
}

func (slice OccurenceGroups) Swap(i int, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (g OccurenceGroup) String() string {
	result := []string{}
	result = append(result, "{")
	result = append(result, "Path: "+g.Path)
	result = append(result, "Type: "+fmt.Sprintf("%d", g.Type))
	result = append(result, "Occurences: [")
	for _, o := range g.Occurences {
		result = append(result, "  "+o.String())
	}
	result = append(result, "]")
	result = append(result, "}")
	return strings.Join(result, "\n")
}

func (slice OccurenceGroups) String() string {
	result := []string{}
	for _, g := range slice {
		result = append(result, g.String()+"\n")
	}
	return strings.Join(result, "\n")
}
