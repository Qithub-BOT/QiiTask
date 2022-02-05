package cmdsort_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1set/todotxt"
	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsort"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/kami-zh/go-capturer"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	appInfo, err := appinfo.New(t.TempDir(), t.TempDir(), "")
	require.NoError(t, err)

	obj1 := cmdsort.New(appInfo)
	obj2 := cmdsort.New(appInfo)

	assert.NotSame(t, obj1, obj2, "it should not reference the same object")
}

func createTask(t *testing.T, todo string, asDone bool) todotxt.Task {
	t.Helper()

	newTask := todotxt.Task{
		Todo: todo,
	}

	if asDone {
		newTask.Completed = asDone
	}

	return newTask
}

func TestSort(t *testing.T) {
	tmpDir := t.TempDir()
	appInfo, err := appinfo.New(tmpDir, tmpDir, "")
	require.NoError(t, err)

	retrunOrigin := util.ChDir(tmpDir)
	defer retrunOrigin()

	appInfo.Tasks.Global.TaskList = &todotxt.TaskList{
		createTask(t, "global task 1", true),
		createTask(t, "global task 2", true),
	}
	appInfo.Tasks.Local.TaskList = &todotxt.TaskList{
		createTask(t, "local task 1", true),
		createTask(t, "local task 2", true),
		createTask(t, "local task 3", false),
		createTask(t, "local task 4", false),
		createTask(t, "local task 5", true),
	}

	obj := new(cmdsort.Command)
	obj.Command = new(cobra.Command)
	obj.AppInfo = appInfo
	obj.CUI = cui.New()

	obj.CUI.ForceTrue = true
	obj.CUI.ForceString = "task"

	cmdRoot := new(cobra.Command)
	cmdRoot.AddCommand(obj.Command)

	out := capturer.CaptureOutput(func() {
		err = obj.Sort(cmdRoot, []string{})
		require.NoError(t, err)
	})

	assert.Empty(t, out)

	require.FileExists(t, todo.NameFile)

	savedTask, err := os.ReadFile(todo.NameFile)
	require.NoError(t, err)

	assert.Contains(t, string(savedTask), util.HereDoc(`
		local task 3
		local task 4
		x local task 5
		x local task 2
		x local task 1
	`))
}

func TestSort_query_not_ready(t *testing.T) {
	tmpDir := t.TempDir()

	appInfo, err := appinfo.New(tmpDir, tmpDir, "")
	require.NoError(t, err)

	retrunOrigin := util.ChDir(tmpDir)
	defer retrunOrigin()

	appInfo.Tasks.Global.TaskList = &todotxt.TaskList{
		createTask(t, "global task 1", true),
	}
	appInfo.Tasks.Local.TaskList = &todotxt.TaskList{
		createTask(t, "local task 1", true),
	}

	obj := new(cmdsort.Command)
	obj.Command = new(cobra.Command)
	obj.AppInfo = appInfo

	obj.CUI = cui.New()
	obj.CUI.ForceError = false
	obj.CUI.ForceString = "unknown query type"

	cmdRoot := new(cobra.Command)
	cmdRoot.AddCommand(obj.Command)

	out := capturer.CaptureOutput(func() {
		err := obj.Sort(cmdRoot, []string{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error during query type selection")
	})

	assert.Empty(t, out)
}

func TestSort_query_not_ready_and_fail_save(t *testing.T) {
	tmpDir := t.TempDir()

	tmpFileTask := filepath.Join(tmpDir, "todo.txt")
	os.WriteFile(tmpFileTask, []byte{}, 0o600)

	appInfo, err := appinfo.New(tmpDir, t.TempDir(), "")
	require.NoError(t, err)

	retrunOrigin := util.ChDir(tmpDir)
	defer retrunOrigin()

	appInfo.Tasks.Global.TaskList = &todotxt.TaskList{
		createTask(t, "global task 1", true),
	}
	appInfo.Tasks.Local.TaskList = &todotxt.TaskList{
		createTask(t, "local task 1", true),
	}

	obj := new(cmdsort.Command)
	obj.Command = new(cobra.Command)
	obj.AppInfo = appInfo

	// Select query type as "task"
	obj.CUI = cui.New()
	obj.CUI.ForceError = false
	obj.CUI.ForceTrue = false
	obj.CUI.ForceString = "task"

	cmdRoot := new(cobra.Command)
	cmdRoot.AddCommand(obj.Command)

	out := capturer.CaptureOutput(func() {
		err := obj.Sort(cmdRoot, []string{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error during confirmation: unknown error occurred during confirmation: EOF")
	})

	assert.Contains(t, out, "設定ファイルが見つかりません。タスクと同じ階層に作成しますか？")
}

func TestSort_no_task_found(t *testing.T) {
	tmpDir := t.TempDir()
	appInfo, err := appinfo.New(tmpDir, tmpDir, "")
	require.NoError(t, err)

	retrunOrigin := util.ChDir(tmpDir)
	defer retrunOrigin()

	obj := new(cmdsort.Command)
	obj.Command = new(cobra.Command)
	obj.AppInfo = appInfo

	obj.CUI = cui.New()
	obj.CUI.ForceError = false
	obj.CUI.ForceString = "task"

	cmdRoot := new(cobra.Command)
	cmdRoot.AddCommand(obj.Command)

	out := capturer.CaptureOutput(func() {
		err := obj.Sort(cmdRoot, []string{})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "task is empty. No tasks to sort")
	})

	assert.Empty(t, out)
}

func TestSurvayQueryType(t *testing.T) {
	tmpDir := t.TempDir()
	appInfo, err := appinfo.New(tmpDir, tmpDir, "")
	require.NoError(t, err)

	obj := new(cmdsort.Command)
	obj.Command = cmdsort.New(appInfo)
	obj.AppInfo = appInfo
	obj.CUI = cui.New()

	{
		obj.CUI.ForceError = true

		err := obj.SurvayQueryType()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error during sort process")
	}
	{
		obj.CUI.ForceError = false
		obj.CUI.ForceString = "unknown type"

		err := obj.SurvayQueryType()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "error during query type selection")
	}
	{
		obj.CUI.ForceError = false
		obj.CUI.ForceString = "task"

		err := obj.SurvayQueryType()

		require.NoError(t, err, "'task' should be an existing query type key")
	}
}

func TestSurvaySave(t *testing.T) {
	tmpDir := t.TempDir()

	tmpFileTask := filepath.Join(tmpDir, "todo.txt")
	os.WriteFile(tmpFileTask, []byte{}, 0o600)

	appInfo, err := appinfo.New(tmpDir, t.TempDir(), "")
	require.NoError(t, err)

	retrunOrigin := util.ChDir(tmpDir)
	defer retrunOrigin()

	appInfo.Tasks.Global.TaskList = &todotxt.TaskList{
		createTask(t, "global task 1", true),
	}
	appInfo.Tasks.Local.TaskList = &todotxt.TaskList{
		createTask(t, "local task 1", true),
	}

	obj := new(cmdsort.Command)

	obj.Command = new(cobra.Command)
	obj.AppInfo = appInfo
	obj.CUI = cui.New()

	{
		obj.CUI.ForceError = false
		obj.CUI.ForceTrue = true
		obj.CUI.ForceFalse = false
		obj.CUI.ForceString = "task"

		out := capturer.CaptureOutput(func() {
			err := obj.SurvaySave()

			require.Error(t, err)
			assert.Contains(t, err.Error(), ".qiitask/config.json: no such file or directory")
		})

		assert.Empty(t, out)
	}
	{
		obj.CUI.ForceError = false
		obj.CUI.ForceTrue = false
		obj.CUI.ForceFalse = true
		obj.CUI.ForceString = "task"

		out := capturer.CaptureOutput(func() {
			err := obj.SurvaySave()

			require.NoError(t, err)
		})

		assert.Contains(t, out, "保存しませんでした。（質問タイプを毎回選択する必要があります）")
	}
}
