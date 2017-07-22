package lister_test

import (
	"testing"

	"strings"

	"github.com/jeffijoe/total-rename/lister"
	"github.com/jeffijoe/total-rename/util"
)

func TestListFileNodes(t *testing.T) {
	fixturePath := "../_fixtures/fixture1/input/**/*.*"
	result, _ := lister.ListFileNodes(util.GetWD(), fixturePath)

	contains(t, result, "spaces", lister.NodeTypeDir)
	contains(t, result, "space-repository.js", lister.NodeTypeFile)
	contains(t, result, "spaceTypes.js", lister.NodeTypeFile)
	contains(t, result, "SPACE_STUFFS.js", lister.NodeTypeFile)
}

func TestListFileNodes_Nested(t *testing.T) {
	fixturePath := "../_fixtures/fixture2/input/**/*.*"
	result, _ := lister.ListFileNodes(util.GetWD(), fixturePath)

	contains(t, result, "space-a", lister.NodeTypeDir)
	contains(t, result, "SPACE-a-a", lister.NodeTypeDir)
	contains(t, result, "get_spaces.js", lister.NodeTypeFile)
	contains(t, result, "Spaces-b", lister.NodeTypeDir)
	contains(t, result, "findSpaces.js", lister.NodeTypeFile)
	contains(t, result, "to-space-we.go", lister.NodeTypeFile)
}

func contains(t *testing.T, result []*lister.FileNode, name string, nodeType lister.NodeType) {
	for _, f := range result {
		if strings.HasSuffix(f.Path, name) {
			if f.Type == nodeType {
				return
			}
			t.Errorf("File %s did not have expected node type %d, but was %d", f.Path, nodeType, f.Type)
		}
	}
	t.Errorf("Did not find file %s", name)
}
