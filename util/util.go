package util

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetWD returns the current working directory.
func GetWD() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.FromSlash(dir)
}

// CopyDir copies a directory
func CopyDir(source, target string) error {
	source = filepath.FromSlash(source)
	target = filepath.FromSlash(target)
	return filepath.Walk(source, func(dir string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			newDir := filepath.Join(target, strings.TrimPrefix(dir, source))
			err := os.MkdirAll(newDir, fi.Mode())
			if err != nil {
				return err
			}
		} else {
			newFile := filepath.Join(
				target,
				strings.TrimPrefix(dir, source),
			)
			err := CopyFile(
				dir,
				newFile,
			)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// CopyFile copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	source = filepath.FromSlash(source)
	dest = filepath.FromSlash(dest)
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err == nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return err
}
