package util

import "os"
import "path/filepath"

// GetWD returns the current working directory.
func GetWD() string {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	return exPath
}
