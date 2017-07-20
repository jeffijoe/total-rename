package lister

import (
	"os"
	"path/filepath"

	"fmt"

	zglob "github.com/mattn/go-zglob"
)

// NodeType determines whether a file is a directory or a file.
type NodeType uint8

const (
	NodeTypeFile = NodeType(1)
	NodeTypeDir  = NodeType(2)
)

func Tree(root, glob string) ([]*FileNode, error) {
	empty := make([]*FileNode, 0)
	path := filepath.Join(root, glob)
	files, err := zglob.Glob(path)
	if err != nil {
		return empty, err
	}

	seenFolders := make(map[string]bool)
	result := []*FileNode{}
	for _, file := range files {
		if err != nil {
			return empty, err
		}

		dir := filepath.Dir(file)
		if _, prs := seenFolders[dir]; !prs {
			seenFolders[dir] = true
			result = append(result, &FileNode{
				Path: dir,
				Type: NodeTypeDir,
			})
		}
		result = append(result, &FileNode{
			Path: file,
			Type: NodeTypeFile,
		})
	}

	return result, nil
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

type FileNode struct {
	Type NodeType
	Path string
}

func (f FileNode) String() string {
	return fmt.Sprintf("{Path = %s, Type = %d}", f.Path, f.Type)
}
