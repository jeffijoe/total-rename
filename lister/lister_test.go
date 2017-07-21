package lister_test

import (
	"testing"

	"strings"

	"github.com/jeffijoe/total-rename/lister"
	"github.com/jeffijoe/total-rename/util"
)

func TestListFileNodes(t *testing.T) {
	fixturePath := "../_fixtures/fixture1/input/**/*.js"
	result, _ := lister.ListFileNodes(util.GetWD(), fixturePath)

	containsFile(t, result, "spaces", lister.NodeTypeDir)
	containsFile(t, result, "space-repository.js", lister.NodeTypeFile)
	containsFile(t, result, "spaceTypes.js", lister.NodeTypeFile)
	containsFile(t, result, "SPACE_STUFFS.js", lister.NodeTypeFile)
}

func containsFile(t *testing.T, result []*lister.FileNode, name string, nodeType lister.NodeType) {
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
