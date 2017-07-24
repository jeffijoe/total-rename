package lister

import (
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"

	"fmt"

	"github.com/jeffijoe/total-rename/simplematch"
	zglob "github.com/mattn/go-zglob"
)

// NodeType determines whether a file is a directory or a file.
type NodeType uint8

// Node types
const (
	NodeTypeFile = NodeType(1)
	NodeTypeDir  = NodeType(2)
)

// FileNodes is a list of file nodes.
type FileNodes []*FileNode

// FileNode is a file node.
type FileNode struct {
	Type NodeType
	Path string
}

// ListFileNodes lists file nodes relative from root matching the specified glob.
func ListFileNodes(root, glob, ignorePattern string) (FileNodes, error) {
	root = filepath.Clean(filepath.FromSlash(root))
	empty := FileNodes{}
	var path string
	var err error
	ignore := simplematch.NewMatcher(ignorePattern)
	if filepath.IsAbs(glob) {
		//path = filepath.FromSlash(glob)
		path = glob
	} else {
		if strings.HasPrefix(glob, "~") {
			user, err := user.Current()
			if err != nil {
				return empty, err
			}
			path = filepath.Join(user.HomeDir, glob[1:])
		} else {
			path = filepath.FromSlash(filepath.Join(root, glob))
		}
	}
	path, err = filepath.Abs(path)
	files, err := zglob.Glob(path)

	if err != nil {
		return empty, err
	}

	seenFolders := make(map[string]struct{})
	result := FileNodes{}
	for _, file := range files {
		if err != nil {
			return empty, err
		}
		fi, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return empty, err
		}

		if !fi.IsDir() {
			result = gatherDirectories(root, filepath.Dir(file), result, seenFolders, ignore)
			if ignore.Matches(file) {
				continue
			}
			result = append(result, &FileNode{
				Path: filepath.FromSlash(file),
				Type: NodeTypeFile,
			})
		}
	}
	sort.Sort(result)
	return result, nil
}

func gatherDirectories(root, dir string, result FileNodes, seenFolders map[string]struct{}, ignore *simplematch.Matcher) FileNodes {
	dir = filepath.Clean(dir)
	for {
		if dir == root {
			return result
		}
		if _, prs := seenFolders[dir]; prs {
			return result
		}

		seenFolders[dir] = struct{}{}
		if ignore.Matches(dir) {
			dir = filepath.Clean(filepath.Dir(dir))
			continue
		}

		result = append(result, &FileNode{
			Path: filepath.FromSlash(dir),
			Type: NodeTypeDir,
		})
		dir = filepath.Clean(filepath.Dir(dir))
	}
}

func getNodeType(path string) (NodeType, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return NodeTypeFile, err
	}

	if stat.IsDir() {
		return NodeTypeDir, nil
	}
	return NodeTypeFile, nil
}

func (f FileNode) String() string {
	return fmt.Sprintf("{Path = %s, Type = %d}", f.Path, f.Type)
}

func (slice FileNodes) Len() int {
	return len(slice)
}

func (slice FileNodes) Less(i, j int) bool {
	return slice[i].Type < slice[j].Type
}

func (slice FileNodes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
