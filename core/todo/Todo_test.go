package todo_test

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/1set/todotxt"
	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTodoByKey_key_not_exist(t *testing.T) {
	pathDirTask := GetPathFromRoot(t, "testdata/sort/random_num_no_priority/")

	tasks, err := todo.New(pathDirTask)
	require.NoError(t, err)

	// Request negative
	task1, err := tasks.GetTodoByKey(-1)

	require.Empty(t, task1, "on error the value should be empty")
	assert.Error(t, err, "requesting excessive number of the slice should be an error")

	// Request over number
	lenTask := len([]todotxt.Task(*tasks.TaskList))
	task2, err := tasks.GetTodoByKey(lenTask + 1)

	require.Empty(t, task2, "on error the value should be empty")
	assert.Error(t, err, "requesting excessive number of the slice should be an error")
}

func TestGetTodoByKey_sort_with_no_priority(t *testing.T) {
	pathDirTask := GetPathFromRoot(t, "testdata/sort/random_num_no_priority/")

	expectOut, err := os.ReadFile(filepath.Join(pathDirTask, "expect_out.txt"))
	require.NoError(t, err, "failed to read file for expected results")

	task, err := todo.New(pathDirTask)
	require.NoError(t, err)

	isALessThanB := func(a, b int) bool {
		taskA, err := task.GetTodoByKey(a)
		require.NoError(t, err, "Index of A: %v", a)

		A, err := strconv.Atoi(taskA)
		require.NoError(t, err, "Todo A: %v", A)

		taskB, err := task.GetTodoByKey(b)
		require.NoError(t, err, "Index of B: %v", b)

		B, err := strconv.Atoi(taskB)
		require.NoError(t, err, "Todo B: %v", B)

		result := A < B
		t.Logf("Is %v(%v) less than %v(%v) -> %v", A, a, B, b, result)

		return result
	}

	task.CustomSort(isALessThanB)

	expect := strings.TrimSpace(string(expectOut))
	actual := strings.TrimSpace(task.String())

	assert.Equal(t, expect, actual)
}

func TestIsKeyDone(t *testing.T) {
	pathDirTask := GetPathFromRoot(t, "testdata/golden/working_with_task_and_conf/")

	obj, err := todo.New(pathDirTask)
	require.NoError(t, err)

	for _, test := range []struct {
		msgErr string
		key    int
		expect bool
	}{
		{
			"the IsKeyDone method should return false if the task is undone.\nthe 1st task in %v/todo.txt is undone",
			0,
			false,
		},
		{
			"the IsKeyDone method should return true if the task is done.\nthe 2nd task in %v/todo.txt is done",
			1,
			true,
		},
		{
			"the IsKeyDone method should return false if the key if out-of-range",
			100,
			false,
		},
	} {
		if test.expect {
			require.True(t, obj.IsKeyDone(test.key), test.msgErr)
		} else {
			require.False(t, obj.IsKeyDone(test.key), test.msgErr)
		}
	}
}

func TestLen(t *testing.T) {
	pathDirTask := GetPathFromRoot(t, "testdata/golden/working_with_task_and_conf/")

	obj, err := todo.New(pathDirTask)
	require.NoError(t, err)

	expect := 2
	actual := obj.Len()
	assert.Equal(t, expect, actual)
}

func TestLen_no_task_loaded(t *testing.T) {
	obj, err := todo.New(t.TempDir())
	require.NoError(t, err)

	expect := -1
	actual := obj.Len()
	assert.Equal(t, expect, actual, "it should be -1 if no task file was loaded")
}

func TestOverWrite(t *testing.T) {
	pathDirTask := GetPathFromRoot(t, "testdata/golden/working_with_task_and_conf/")

	obj, err := todo.New(pathDirTask)
	require.NoError(t, err, "failed to create object during test")

	ui := cui.New()
	ui.ForceTrue = true

	err = obj.OverWrite(ui)
	require.NoError(t, err, "overwriting the loaded file should not be an error")
}

func TestOverWrite_accept_overwrite(t *testing.T) {
	pathDirTmp := t.TempDir()

	deferReturn := util.ChDir(pathDirTmp)
	defer deferReturn()

	obj, err := todo.New(pathDirTmp)
	require.NoError(t, err, "failed to create object during test")

	ui := cui.New()
	ui.ForceTrue = true

	err = obj.OverWrite(ui)

	require.NoError(t, err, "if user seletcs yes(y/Y) then it should write a file")
	assert.FileExists(t, todo.NameFile)
}

func TestOverWrite_deny_overwrite(t *testing.T) {
	obj, err := todo.New(t.TempDir())
	require.NoError(t, err, "failed to create object during test")

	ui := cui.New()
	ui.ForceFalse = true

	err = obj.OverWrite(ui)

	require.Error(t, err,
		"if user choice was no(n/N) then it should return an error")
	assert.Contains(t, err.Error(), "保存しませんでした")
	assert.Contains(t, err.Error(), "タスクは破棄されました")
}

func TestOverWrite_task_not_loaded(t *testing.T) {
	obj, err := todo.New(t.TempDir())
	require.NoError(t, err, "failed to create object during test")

	ui := cui.New()
	ui.ForceError = true

	err = obj.OverWrite(ui)

	require.Error(t, err,
		"using OverWrite method when no task was loaded then it should be an error")
	assert.Contains(t, err.Error(), "task save has been canceled: forced error")
}
