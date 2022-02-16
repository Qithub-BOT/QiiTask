package todo_test

import (
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-utiles/util"
)

// ----------------------------------------------------------------------------
//  Helper Functions
// ----------------------------------------------------------------------------

// GetPathFromRoot はリポジトリのルートから見て相対パスで指定されたディレクトリ
// の絶対パスを返します。
func GetPathFromRoot(t *testing.T, path string) string {
	t.Helper()

	pathRoot := util.GetPathDirRepo()
	if pathRoot == "" {
		t.Fatal("'.git' not found: failed to get repository root directory")
	}

	pathTarget := filepath.Clean(filepath.Join(pathRoot, path))
	if !util.IsDir(pathTarget) {
		t.Fatalf("dir not found: %v", pathTarget)
	}

	return pathTarget
}
