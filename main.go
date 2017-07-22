package main

import (
	"flag"

	"fmt"

	tm "github.com/buger/goterm"
)

func main() {
	help := flag.Bool("help", false, "Shows the help menu")
	dryRun := flag.Bool("dry", false, "If set, won't rename anything.")
	force := flag.Bool("force", false, "Replaces all occurences without asking")
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *dryRun {
		fmt.Println("Dry run active, won't rename anything.")
	}

	if *force {
		fmt.Println("Not gonna ask for permission")
	}

	//tm.Flush()
}

func clearLine() {
	tm.ResetLine("")
	for i := 0; i < tm.Width(); i++ {
		tm.Print(" ")
	}
	tm.MoveCursorUp(1)
}

func printHelp() {
	fmt.Println("total-rename - case-preserving renaming utility")
	fmt.Println("")
	fmt.Println("Example: rename all occurences of \"awesome\" to \"excellent\"")
	fmt.Println("")
	fmt.Println("    total-rename **/*.txt awesome excellent")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("    --dry   - If set, won't rename anything")
	fmt.Println("    --force - Replaces all occurences without asking")
	fmt.Println("    --help  - Shows this help text")
}
