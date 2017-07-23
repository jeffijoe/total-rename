package replacer

import (
	"strings"

	"io/ioutil"
	"os"

	"unicode/utf8"

	"github.com/jeffijoe/total-rename/casing"
	"github.com/jeffijoe/total-rename/scanner"
)

// RenameFunc describes a function used to rename a file/folder.
type RenameFunc func(oldPath, newPath string) error

// ReplaceFileFunc describes a function used to replace the contents of a file.
type ReplaceFileFunc func(filePath, newContent string) error

// TotalRenameResult describes the result of calling TotalRename
type TotalRenameResult struct {
	OccurencesRenamed int
}

// TotalRename will rename files and paths.
func TotalRename(groups scanner.OccurenceGroups, replacement string, rename RenameFunc, replaceFile ReplaceFileFunc) (*TotalRenameResult, error) {
	renamed := 0
	replacementVariants := casing.GenerateCasings(replacement)
	for _, group := range groups {
		var count int
		var err error
		switch group.Type {
		case scanner.OccurenceGroupTypeContent:
			count, err = totalRenameFile(group, replacementVariants, replaceFile)
		case scanner.OccurenceGroupTypePath:
			count, err = totalRenamePath(group, replacementVariants, rename)
		}
		if err != nil {
			return nil, err
		}
		renamed = renamed + count
	}

	return &TotalRenameResult{
		OccurencesRenamed: renamed,
	}, nil
}

func totalRenameFile(group *scanner.OccurenceGroup, replacement casing.Variants, replaceFile ReplaceFileFunc) (int, error) {
	contentBytes, err := ioutil.ReadFile(group.Path)
	if err != nil {
		return 0, err
	}

	content := string(contentBytes)
	newContent := ReplaceText(content, group.Occurences, replacement)
	err = replaceFile(group.Path, newContent)
	if err != nil {
		return 0, err
	}
	return len(group.Occurences), nil
}

func totalRenamePath(group *scanner.OccurenceGroup, replacement casing.Variants, rename RenameFunc) (int, error) {
	newPath := ReplaceText(group.Path, group.Occurences, replacement)
	if err := rename(group.Path, newPath); err != nil {
		return 0, err
	}
	return len(group.Occurences), nil
}

// ReplaceText teplaces all occurences with their replacement variants
// Occurences should be ordered by StartIndex.
func ReplaceText(source string, occurences scanner.Occurences, replacementVariants casing.Variants) string {
	occurenceCount := len(occurences)
	if occurenceCount == 0 {
		return source
	}

	// String slices, cut at each index and leaving out the match.
	slices := make([]string, 0, len(occurences)+2)

	buf := make([]rune, 0, len(source))
	allRunes := []rune(source)
	ocIdx := 0
	for idx := 0; idx < len(allRunes); idx++ {
		charCode := allRunes[idx]
		var oc *scanner.Occurence
		if ocIdx != occurenceCount {
			oc = occurences[ocIdx]
		}

		if oc != nil && idx == oc.StartIndex {
			ocIdx = ocIdx + 1
			slices = append(slices, string(buf))
			buf = make([]rune, 0, len(source))
			idx = idx + utf8.RuneCountInString(oc.Match) - 1
		} else {
			buf = append(buf, charCode)
		}

	}
	slices = append(slices, string(buf))
	result := make([]string, 1, len(slices)+occurenceCount)
	result[0] = slices[0]
	for idx, oc := range occurences {
		v := replacementVariants.GetVariant(oc.Casing)
		result = append(result, v.Value, slices[idx+1])
	}

	return strings.Join(result, "")
}

// ReplaceFileContent replaces the file contents.
func ReplaceFileContent(filePath, newContent string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, []byte(newContent), fileInfo.Mode())
	return nil
}
