package main

import (
	"flag"
	"os"
	"runtime"
	"strconv"

	"fmt"

	"github.com/fatih/color"
	"github.com/jeffijoe/total-rename/casing"
	"github.com/jeffijoe/total-rename/cli"
	"github.com/jeffijoe/total-rename/lister"
	"github.com/jeffijoe/total-rename/replacer"
	"github.com/jeffijoe/total-rename/scanner"
	"github.com/jeffijoe/total-rename/util"
)

func main() {
	help := flag.Bool("help", false, "Shows the help menu")
	dryRun := flag.Bool("dry", false, "If set, won't rename anything.")
	force := flag.Bool("force", false, "Replaces all occurences without asking")
	flag.Parse()
	fmt.Println("total-rename - case-preserving renaming utility")
	fmt.Println("Copyright © Jeff Hansen 2017 to present. All rights reserved.")
	fmt.Println()
	if *help {
		printHelp()
		return
	}

	if *dryRun {
		fmt.Println("--dry active; won't rename anything.")
	}

	if *force {
		fmt.Println("--force active; won't prompt for confirmation")
	}

	fmt.Println()
	if flag.NArg() < 3 {
		fmt.Println("Not enough arguments, expects 3: <path> <needle> <replacement>")
		return
	}
	replacement := flag.Arg(2)
	path := flag.Arg(0)
	needle := flag.Arg(1)
	nodes, err := lister.ListFileNodes(util.GetWD(), path)
	if err != nil {
		panic(err)
	}
	var groups scanner.OccurenceGroups
	if *force {
		groups, err = scanner.ScanFileNodes(nodes, needle)
	} else {
		groups, err = promptOccurences(nodes, needle, replacement)
	}
	if err != nil {
		panic(err)
	}

	rename := os.Rename
	replace := replacer.ReplaceFileContent
	if *dryRun {
		rename = func(p1, p2 string) error {
			return nil
		}
		replace = func(p1, p2 string) error {
			return nil
		}
	}

	result, err := replacer.TotalRename(groups, replacement, rename, replace)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Done! Renamed %d occurences!", result.OccurencesRenamed)
	fmt.Println()
}

func promptOccurences(nodes lister.FileNodes, needle, replacement string) (scanner.OccurenceGroups, error) {
	groups, err := scanner.ScanFileNodes(nodes, needle)
	if err != nil {
		return nil, err
	}
	replacementVariants := casing.GenerateCasings(replacement)
	result := scanner.OccurenceGroups{}
	for _, group := range groups {
		var newGroup *scanner.OccurenceGroup
		switch group.Type {
		case scanner.OccurenceGroupTypeContent:
			newGroup, err = promptGroup(group, replacementVariants, promptContentOccurence)
		case scanner.OccurenceGroupTypePath:
			newGroup, err = promptGroup(group, replacementVariants, promptPathOccurence)
		}
		if err != nil {
			return nil, err
		}
		if newGroup != nil {
			result = append(result, newGroup)
		}
	}

	return result, nil
}

// OccurencePrompter is a function that prompts the user whether the occurence should be replaced or not.
type OccurencePrompter func(occurence *scanner.Occurence, replacementVariants casing.Variants, w *cli.Wrapper) (bool, error)

func promptGroup(group *scanner.OccurenceGroup, replacementVariants casing.Variants, promptOccurence OccurencePrompter) (*scanner.OccurenceGroup, error) {
	w := cli.Clearable()

	occurences := scanner.Occurences{}
	result := &scanner.OccurenceGroup{
		Path: group.Path,
		Type: group.Type,
	}
	countReplaced := 0
	countSkipped := 0
	printFileStatus := func(printf func(string, ...interface{}) (int, error)) {
		color.Set(color.BgWhite)
		color.Set(color.FgBlack)
		printf(group.Path)
		color.Set(color.BgGreen)

		if countReplaced > 0 {
			printf(" %d replaced", countReplaced)
		}
		if countSkipped > 0 {
			if countReplaced > 0 {
				printf(", %d skipped", countSkipped)
			} else {
				printf(": %d skipped", countSkipped)

			}
		}
		color.Unset()
		printf("\n")
	}
	for _, oc := range group.Occurences {
		printFileStatus(w.Printf)
		w.Println()
		shouldReplace, err := promptOccurence(oc, replacementVariants, w)
		if err != nil {
			return nil, err
		}
		if shouldReplace {
			countReplaced = countReplaced + 1
			occurences = append(occurences, oc)
		} else {
			countSkipped = countSkipped + 1
		}

		w.Clear()
	}
	if runtime.GOOS != "windows" {
		printFileStatus(fmt.Printf)
	}
	if len(occurences) == 0 {
		return nil, nil
	}
	result.Occurences = occurences
	return result, nil
}

func promptPathOccurence(occurence *scanner.Occurence, replacementVariants casing.Variants, w *cli.Wrapper) (bool, error) {
	color.Set(color.FgHiBlack)
	beforeMatch := occurence.Line[:occurence.LineStartIndex]
	afterMatch := occurence.Line[occurence.LineStartIndex+len(occurence.Match):]
	w.Println("Occurence in path:")
	w.Print("   ")
	w.Printf(beforeMatch)
	color.Set(color.FgYellow)
	w.Print(occurence.Match)
	color.Set(color.FgHiBlack)
	w.Println(afterMatch)

	w.Println()
	color.Set(color.FgWhite)
	w.Print("Replace ")
	color.Set(color.FgYellow)
	w.Print(occurence.Match)
	color.Set(color.FgWhite)
	w.Print(" with ")
	color.Set(color.FgGreen)
	w.Print(replacementVariants.GetVariant(occurence.Casing).Value)
	color.Set(color.FgWhite)
	w.Println("? [Y/n] ")
	response, err := w.Confirm(true)
	return response, err
}

func promptContentOccurence(occurence *scanner.Occurence, replacementVariants casing.Variants, w *cli.Wrapper) (bool, error) {
	color.Set(color.FgHiBlack)
	for i, ln := range occurence.SurroundingLinesBefore {
		lineNum := occurence.LineNumber + i + 1 - len(occurence.SurroundingLinesBefore)
		w.Println(formatLine(lineNum, ln))
	}
	beforeMatch := occurence.Line[:occurence.LineStartIndex]
	afterMatch := occurence.Line[occurence.LineStartIndex+len(occurence.Match):]
	w.Printf(formatLine(occurence.LineNumber+1, beforeMatch))
	color.Set(color.FgYellow)
	w.Print(occurence.Match)
	color.Set(color.FgHiBlack)
	w.Println(afterMatch)
	for i, ln := range occurence.SurroundingLinesAfter {
		lineNum := occurence.LineNumber + i + 2
		w.Println(formatLine(lineNum, ln))
	}
	color.Unset()
	w.Println()
	color.Set(color.FgWhite)
	w.Print("Replace ")
	color.Set(color.FgYellow)
	w.Print(occurence.Match)
	color.Set(color.FgWhite)
	w.Print(" with ")
	color.Set(color.FgGreen)
	w.Print(replacementVariants.GetVariant(occurence.Casing).Value)
	color.Set(color.FgWhite)
	w.Print("? [Y/n] ")
	response, err := w.Confirm(true)
	return response, err
}

func formatLine(lineNum int, str string) string {
	return fmt.Sprintf("%6s: %s", strconv.Itoa(lineNum), str)
}

func printHelp() {
	fmt.Println("OPTIONS:")
	fmt.Println("    Options must be specified before arguments.")
	fmt.Println("")
	fmt.Println("    --dry         If set, won't rename anything")
	fmt.Println("    --force       Replaces all occurences without asking")
	fmt.Println("    --help        Shows this help text")
	fmt.Println("")
	fmt.Println("ARGUMENTS:")
	fmt.Println("")
	fmt.Println("    <pattern>  Search pattern (glob). Relative to working")
	fmt.Println("               directory unless rooted (absolute path).")
	fmt.Println("    <find>     The string to find. If multiple words,")
	fmt.Println("               please use camelCase.")
	fmt.Println("    <replace>  The string to replace occurences with.")
	fmt.Println("               If multiple words, please use camelCase.")
	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("")
	fmt.Println("    total-rename \"**/*.txt\" \"awesome\" \"excellent\"")
	fmt.Println("")
	fmt.Println("    Rename all occurences of \"awesome\" to \"excellent\" in")
	fmt.Println("    all .txt files (and folders) recursively from the")
	fmt.Println("    current working directory:")
	fmt.Println("")
	fmt.Println("EXAMPLE:")
	fmt.Println("")
	fmt.Println("    total-rename --force \"/Users/jeff/projects/my-app/src/**/*.*\" \"awesome\" \"excellent\"")
	fmt.Println("")
	fmt.Println("    Like the first example, but from an absolute path, and match all")
	fmt.Println("    file extensions and don't ask for confirmation for each occurence.")
	fmt.Println("")
}
