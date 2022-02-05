package todo

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/1set/todotxt"
	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------

// NameFile はタスクのファイル名です。この定数値がファイルの読み込みに使われます。
const NameFile = "todo.txt"

// ----------------------------------------------------------------------------
//  型の定義
// ----------------------------------------------------------------------------

// Todo は todotxt.TaskList 型（github.com/1set/todotxt）を埋め込んだ拡張型です。
type Todo struct {
	*todotxt.TaskList
	pathDir  string // タスクの保存先ディレクトリ
	pathFile string // タスク・ファイルのパス（存在した場合）
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New は Todo の新規オブジェクトを返します。pathDir のディレクトリ以下をサーチ
// し、最初に見つけた "todo.txt" をロードします。
//
// 通常、この関数を使ってタスクオブジェクトを生成することはなく、 NewList() 関
// 数を使ってローカルおよびグローバルにあるタスクがセットになったオブジェクトを
// 使います。
func New(pathDir string) (*Todo, error) {
	taskList := todotxt.NewTaskList()

	obj := &Todo{&taskList, "", ""}

	if err := obj.loadTask(pathDir); err != nil {
		return nil, errors.Wrap(err, "fail to load task")
	}

	return obj, nil
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// CustomSort は isALessThanB を使ってタスクをソートします。
//
// ソートは安定ソートです。また isALessThanB の関数は A が B より小さい場合に
// true を返す必要があります。
func (t *Todo) CustomSort(isALessThanB func(a int, b int) bool) {
	sort.SliceStable(*t.TaskList, isALessThanB)
}

// Dir はタスクの読み込み・書き込みをする際に使われるディレクトリ名を返します。
func (t *Todo) Dir() string {
	return t.pathDir
}

// File はタスクの読み込み・書き込みをする際に使われるファイル名を返します。
func (t *Todo) File() string {
	return NameFile
}

// FileUsed は読み込んだタスク・ファイルのパスを返します。
// ファイルが存在しない場合は "" にセットされています。
func (t *Todo) FileUsed() string {
	return t.pathFile
}

// findFileTask は Todo.pathDir 以下にあるタスク・ファイルを検索して、最初に見つ
// けたファイルのパスを返します。ファイルが見つからない場合は "" を返します。
func (t *Todo) findFileTask() (pathFile string) {
	// 現在のディレクトリ直下
	pathFile = filepath.Join(t.pathDir, NameFile)
	if util.IsFile(pathFile) {
		return pathFile
	}

	// .qiitask 下
	pathFile = filepath.Join(t.pathDir, ".qiitask", NameFile)
	if util.IsFile(pathFile) {
		return pathFile
	}

	// .config/qiitask 下
	pathFile = filepath.Join(t.pathDir, ".config", "qiitask", NameFile)
	if util.IsFile(pathFile) {
		return pathFile
	}

	return ""
}

// GetTodoByKey は t.TaskList[key] にある Todo フィールドの文字列値を返します。
//
// GetTask と違いタスク ID ではなく、スライスのキー番号で指定します（ID は 1 ス
// タート、Key は 0 スタートです）。指定された key が範囲外の場合はエラーを返し
// ます。
func (t *Todo) GetTodoByKey(key int) (string, error) {
	if key < 0 || len(*t.TaskList) < key+1 {
		return "", errors.New("task not found (out of range)")
	}

	return []todotxt.Task(*t.TaskList)[key].Todo, nil
}

// IsKeyDone は t.TaskList[key] にあるタスクが完了済み（"x" 付き）の場合に
// true を返します。key が範囲外/存在しない場合は false を返します。
func (t *Todo) IsKeyDone(key int) bool {
	if key < 0 || len(*t.TaskList) < key+1 {
		return false
	}

	return []todotxt.Task(*t.TaskList)[key].Completed
}

// Len は現在のタスクの長さ（len）を返します。
// タスクもタスクのファイルも存在しない場合は -1 を返します。
func (t *Todo) Len() int {
	length := -1

	if t.TaskList != nil {
		if length = len([]todotxt.Task(*t.TaskList)); length == 0 && t.FileUsed() == "" {
			length = -1
		}
	}

	return length
}

// loadTask はタスク・ファイルを読み込みます。
// ファイルが存在するも、読み込みに失敗した場合は error を返します。ファイルが存
// 在しない場合は error を返さず何もしません。
func (t *Todo) loadTask(pathDir string) error {
	// 保存先・読み込み先のディレクトリとファイルのパスをリセット
	t.pathDir = pathDir
	t.pathFile = ""

	// タスク・ファイルの検索（pathDir 以下を検索します）
	pathFileTarget := t.findFileTask()
	if pathFileTarget == "" {
		return nil
	}

	// タスク・ファイルが存在した場合は読み込み
	tasklist, err := todotxt.LoadFromPath(pathFileTarget)
	if err != nil {
		return errors.Wrap(err,
			"task file found but another error was produced")
	}

	t.TaskList = &tasklist
	t.pathDir = filepath.Dir(pathFileTarget)
	t.pathFile = pathFileTarget

	return nil
}

// OverWrite は現在状態のタスクを読み込み元のファイルに上書きします。
// 新規保存の場合は SaveAs を使います。
func (t *Todo) OverWrite(ui *cui.UI) error {
	if t.FileUsed() == "" || !util.IsFile(t.FileUsed()) {
		pathFileTask := NameFile

		// 保存の問い合わせ
		yes, err := ui.Confirm(fmt.Sprintf(
			"タスク・ファイルがありません。\nタスクを保存しますか? (%v)",
			pathFileTask,
		))
		if err != nil {
			return errors.Wrap(err, "task save has been canceled")
		}

		if !yes {
			return errors.New("保存しませんでした。（タスクは破棄されました）")
		}

		return t.WriteToPath(pathFileTask)
	}

	return t.WriteToPath(t.FileUsed())
}
