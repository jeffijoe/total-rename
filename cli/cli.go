package cli

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	tm "github.com/buger/goterm"
)

// Wrapper is a proxy around fmt so we can manage clearing lines.
type Wrapper struct {
	NewlineCount int
}

// Clearable lets you write to the terminal and clear it afterwards.
func Clearable() *Wrapper {
	w := &Wrapper{}
	return w
}

// Printf does what you expect.
func (w *Wrapper) Printf(fmtString string, args ...interface{}) (int, error) {
	str := fmt.Sprintf(fmtString, args...)
	w.SyncNewlines(str)
	return fmt.Printf(str)
}

// Print does what you expect.
func (w *Wrapper) Print(args ...interface{}) (int, error) {
	str := fmt.Sprint(args...)
	w.SyncNewlines(str)
	return fmt.Print(str)
}

// Println does what you expect.
func (w *Wrapper) Println(args ...interface{}) (int, error) {
	str := fmt.Sprintln(args...)
	w.SyncNewlines(str)
	return fmt.Printf(str)
}

// ReadLine does what you expect.
func (w *Wrapper) ReadLine() (string, error) {
	scanner := bufio.NewReader(os.Stdin)
	str, err := scanner.ReadString('\n')
	if err != nil {
		return "", err
	}
	w.SyncNewlines(str)
	return str, nil
}

// Clear will clear all written lines.
func (w *Wrapper) Clear() {
	// This shit won't work on Windows.
	if runtime.GOOS == "windows" {
		w.Println()
		w.NewlineCount = 0
		return
	}
	// tmWidth := tm.Width()
	if w.NewlineCount == 0 {
		return
	}
	for index := 0; index < w.NewlineCount; index++ {
		// fmt.Print("\033[1A")
		fmt.Printf("\033[2K")
		// for i := 0; i < tmWidth; i++ {
		// 	fmt.Print(" ")
		// }

		fmt.Print("\033[1A")
	}

	w.NewlineCount = 0
}

// Confirm asks for user confirmation.
func (w *Wrapper) Confirm(defaultValue bool) (bool, error) {
	response, err := w.ReadLine()
	if err != nil {
		return defaultValue, err
	}
	response = strings.Trim(response, "\n")
	if response == "y" || response == "Y" {
		return true, nil
	}
	if response == "n" || response == "N" {
		return false, nil
	}

	return defaultValue, nil
}

// SyncNewlines increments the internal newline counter by counting newlines in the specified string.
func (w *Wrapper) SyncNewlines(str string) {
	nlCount := strings.Count(str, "\n")
	splat := strings.Split(str, "\n")
	width := tm.Width()
	for _, s := range splat {
		if len(s) > width {
			nlCount = nlCount + 1
		}
	}
	w.NewlineCount = w.NewlineCount + nlCount
}
