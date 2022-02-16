package todo_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/todo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Examples
// ----------------------------------------------------------------------------

func ExampleNewSet() {
	// Path to search tasks
	pathToCurrentDir := "../../testdata/golden/working_with_task_and_conf/" // Usually "./"
	pathToUserHomeDir := "../../testdata/golden/user_home_with_task"        // Usually "~/"

	taskList, err := todo.NewSet(pathToCurrentDir, pathToUserHomeDir)
	if err != nil {
		log.Fatal(err)
	}

	// Get the first task from current dir (as local task).
	// Note that the 1st item is not "0" like the array key.
	taskLocal, err := taskList.Local.GetTask(1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(taskLocal.Todo)

	// Get the first task from the user home dir (as global task).
	// Note that the 1st item is not "0" like the array key.
	taskGlobal, err := taskList.Global.GetTask(1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(taskGlobal.Todo)

	// Output:
	// This is an un-done task in local (working dir) with conf in local.
	// This is an un-done task in global.
}

// ----------------------------------------------------------------------------
//  Tests
// ----------------------------------------------------------------------------

func TestNewSet(t *testing.T) {
	for _, test := range []struct {
		pathCurr         string
		pathHome         string
		expectTodoLocal  string
		expectTodoGlobal string
	}{
		{
			"testdata/golden/working_with_task_and_conf/",
			"testdata/golden/user_home_with_task",
			"This is an un-done task in local (working dir) with conf in local.",
			"This is an un-done task in global.",
		},
		{
			"testdata/golden/working_task_under_conf/",
			"testdata/golden/user_home_with_task",
			"This is an un-done task in local (working dir) under conf dir.",
			"This is an un-done task in global.",
		},
	} {
		pathCurr := GetPathFromRoot(t, test.pathCurr)
		require.NotEmpty(t, pathCurr)

		pathHome := GetPathFromRoot(t, test.pathHome)
		require.NotEmpty(t, pathHome)

		obj, err := todo.NewSet(pathCurr, pathHome)
		require.NoError(t, err)

		// Test task in current dir (as local task)
		objTaskLocal, err := obj.Local.GetTask(1)
		assert.NoError(t, err)

		expectTaskLocal := test.expectTodoLocal
		actualTaskLocal := objTaskLocal.Todo
		assert.Equal(t, expectTaskLocal, actualTaskLocal)

		// Test task in user home dir (as global task)
		objTaskGlobal, err := obj.Global.GetTask(1)
		assert.NoError(t, err)

		expectTaskGlobal := test.expectTodoGlobal
		actualTaskGlobal := objTaskGlobal.Todo
		assert.Equal(t, expectTaskGlobal, actualTaskGlobal)
	}
}

func TestNewSet_empty_arg(t *testing.T) {
	obj, err := todo.NewSet("", "")
	require.NoError(t, err, "empty args should not be error")

	{
		// Local task's path
		expect, err := os.Getwd()
		require.NoError(t, err)

		actual := obj.Local.Dir()
		assert.Equal(t, expect, actual, "empty 1st arg should set current dir by default")
	}
	{
		// Global task's path
		expect, err := os.UserHomeDir()
		require.NoError(t, err)

		actual := obj.Global.Dir()
		assert.Equal(t, expect, actual, "empty 2nd arg should set user home dir by default")
	}
}

func TestNewSet_empty_dir(t *testing.T) {
	pathDirCurrDummy := t.TempDir() // Assign empty dir
	pathDirHomeDummy := t.TempDir() // Assign empty dir
	obj, err := todo.NewSet(pathDirCurrDummy, pathDirHomeDummy)

	require.NoError(t, err)

	// Local (Current Dir) Task
	assert.Empty(t, obj.Local.FileUsed(),
		"assigning empty dir as curr dir, it should not use any file")
	assert.Equal(t, pathDirCurrDummy, obj.Local.Dir(),
		"it should return the provided directory")
	assert.Equal(t, todo.NameFile, obj.Local.File(),
		"it should return the same file name as todo.NameFile constant is")

	// Global (User Home Dir) Task
	assert.Empty(t, obj.Global.FileUsed(),
		"assigning empty dir as user home dir, it should not use any file")
	assert.Equal(t, pathDirHomeDummy, obj.Global.Dir(),
		"it should return the provided directory")
	assert.Equal(t, "todo.txt", obj.Global.File(),
		"it should return the same file name as todo.NameFile constant is")
}

func TestNewSet_fail_to_get_user_home_path(t *testing.T) {
	// Backup and defer restoration
	oldOsGetwd := todo.OsGetwd
	defer func() {
		todo.OsGetwd = oldOsGetwd
	}()

	// Mock os.Getwd
	todo.OsGetwd = func() (dir string, err error) {
		return "", errors.New("forced fail")
	}

	tasks, err := todo.NewSet("", "")

	require.Nil(t, tasks, "on error returned object should be nil")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get current working dir")
}

func TestNewSet_fail_to_get_working_path(t *testing.T) {
	// Backup and defer restoration
	oldOsUserHomeDir := todo.OsUserHomeDir
	defer func() {
		todo.OsUserHomeDir = oldOsUserHomeDir
	}()

	// Mock os.Getwd
	todo.OsUserHomeDir = func() (dir string, err error) {
		return "", errors.New("forced fail")
	}

	tasks, err := todo.NewSet("", "")

	require.Nil(t, tasks, "on error returned object should be nil")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user home directory")
}

func TestNewSet_invalid_absolute_path(t *testing.T) {
	// Backup and restore before mocking
	oldFilePathAbs := todo.FilePathAbs
	defer func() {
		todo.FilePathAbs = oldFilePathAbs
	}()

	// Mock to force error
	todo.FilePathAbs = func(path string) (string, error) {
		if path == "assert error" {
			return "", errors.New("dummy error")
		}

		return filepath.Abs(path)
	}

	{
		// Error on local task (current dir)
		result, err := todo.NewSet("assert error", ".")
		require.Nil(t, result, "on error the object should be nil")
		require.Error(t, err)
	}
	{
		// Error on global task (user home dir)
		result, err := todo.NewSet(".", "assert error")
		require.Nil(t, result, "on error the object should be nil")
		require.Error(t, err)
	}
}

func TestNewSet_not_singleton(t *testing.T) {
	obj1, err := todo.NewSet(t.TempDir(), t.TempDir())
	require.NoError(t, err)

	obj2, err := todo.NewSet(t.TempDir(), t.TempDir())
	require.NoError(t, err)

	assert.NotSame(t, obj1, obj2, "pointers should not reference the same object")
}

func TestNewSet_scan_error(t *testing.T) {
	for _, test := range []struct {
		pathLocal        string
		pathGlobal       string
		expectErrContain string
	}{
		{
			// Fail local task scan
			"testdata/error/scanner_error/",
			"testdata/golden/user_home_with_task",
			"task file found but another error was produced",
		},
		{
			// Fail global task scan
			"testdata/golden/user_home_with_task",
			"testdata/error/scanner_error/",
			"task file found but another error was produced",
		},
	} {
		pathCurr := GetPathFromRoot(t, test.pathLocal)
		require.NotEmpty(t, pathCurr)

		pathHome := GetPathFromRoot(t, test.pathGlobal)
		require.NotEmpty(t, pathHome)

		obj, err := todo.NewSet(pathCurr, pathHome)
		require.Error(t, err)
		require.Nil(t, obj, "on error object should be nil")
		assert.Contains(t, err.Error(), test.expectErrContain)
	}
}
