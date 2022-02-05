package appinfo_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	exitOnErr := func(err error) {
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Empty args will set the default
	appInfo, err := appinfo.New("", "", "")
	exitOnErr(err)

	// Get expect path (as default path for local task)
	pathDirCurr, err := os.Getwd()
	exitOnErr(err)

	// Get expect path (as default path for global task)
	pathDirHome, err := os.UserHomeDir()
	exitOnErr(err)

	// Assertion
	if appInfo.Version == "" {
		fmt.Println("ok (app version whould be as is)")
	}

	if pathDirCurr == appInfo.Tasks.Local.Dir() {
		fmt.Println("ok")
	}

	if pathDirHome == appInfo.Tasks.Global.Dir() {
		fmt.Println("ok")
	}
	// Output:
	// ok (app version whould be as is)
	// ok
	// ok
}

func TestNew_non_singleton(t *testing.T) {
	obj1, err := appinfo.New("", "", "")
	require.NoError(t, err, "emmpty args should not be error")

	obj2, err := appinfo.New("", "", "")
	require.NoError(t, err, "emmpty args should not be error")

	require.NotSame(t, obj1, obj2, "pointers should not reference the same object")
}

func TestNew_fail_to_get_working_path(t *testing.T) {
	// Backup and defer restoration
	oldOsGetwd := todo.OsGetwd
	defer func() {
		todo.OsGetwd = oldOsGetwd
	}()

	// Mock to force fail
	todo.OsGetwd = func() (dir string, err error) {
		return "", errors.New("forced error")
	}

	appInfo, err := appinfo.New("", "", "")

	require.Error(t, err)
	assert.Nil(t, appInfo, "on error returned object should be nil")
	assert.Contains(t, err.Error(), "failed to instantiate tasks object")
}

func TestNew_load_mal_format_conf(t *testing.T) {
	pathDirApp := filepath.Join(util.GetPathDirRepo(), "testdata", "error", "malformed_conf")

	obj, err := appinfo.New(pathDirApp, t.TempDir(), "")
	require.Error(t, err)
	require.Nil(t, obj, "return object should be nil on error")

	assert.Contains(t, err.Error(), "failed to instantiate config object")
}
