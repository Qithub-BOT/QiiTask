package cmdsort

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/1set/todotxt"
	"github.com/pkg/errors"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/core/config"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/Qithub-BOT/QiiTask/core/query"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/spf13/cobra"
)

// Command is the struct to hold cobra.Command and it's flag options.
type Command struct {
	*cobra.Command
	AppInfo  *appinfo.AppInfo
	CUI      *cui.UI
	numSort  int
	isGlobal bool
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns the newly created object pointer of the "list" command.
func New(appInfo *appinfo.AppInfo) *cobra.Command {
	// Instantiate new object
	cmdSort := new(Command)

	// Add command properties
	cmdSort.Command = &cobra.Command{
		Use:   "sort",
		Short: "タスクを対話式で並べ替えます",
		Long: util.HereDoc(`
				About:
				  'sort' は現在のタスク一覧を対話式でソートします。（未完成タスクのみ）

				  デフォルトで、客観的な質問を使い全件ソートします。"--top" オプション
				  で件数指定があった場合は、主観的な質問を使って n 件ぶんを部分ソートし
				  ます。
			`),
		Example: util.HereDoc(`
				qiitask sort          // 全件ソート
				qiitask sort --top 10 // 部分ソート（上位 10 件）
			`, "  "),
	}

	// Add app properties
	cmdSort.AppInfo = appInfo

	// Add CUI object
	cmdSort.CUI = cui.New()

	// Assign method to execute command
	cmdSort.RunE = cmdSort.Sort

	// Define flags for `list` command.
	cmdSort.Flags().IntVarP(
		&cmdSort.numSort, "top", "t", 0, "一覧のトップ n 件をソートします（部分ソート）",
	)
	cmdSort.Flags().BoolVarP(
		&cmdSort.isGlobal, "golobal", "g", false, "強制的にグローバル・タスクをソートします",
	)

	return cmdSort.Command
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

func (c *Command) askIsALessThanB(a, b string, indexQ int) bool {
	questions := c.AppInfo.Config.GetQueryObjective()

	// Reset index
	if len(questions) <= indexQ {
		indexQ = 0
	}

	idontknow := "質問を変える"
	selections := []string{
		a, b, idontknow,
	}

	answer, err := c.CUI.Select(questions[indexQ], selections, idontknow, "")
	// TODO: fatal はよろしくない
	if err != nil {
		log.Fatalf("強制終了されました（タスクに変更はありません）: %v\n", err)
	}

	if answer == idontknow {
		indexQ++

		c.CUI.DrawHR()

		return c.askIsALessThanB(a, b, indexQ)
	}

	return answer == a
}

// GetTaks は現在のタスク一覧のを返します。完了済のタスクはソートされた状態で返されます。
func (c *Command) GetTaks() (*todo.Todo, error) {
	taskList := c.AppInfo.Tasks.Local

	if c.isGlobal {
		taskList = c.AppInfo.Tasks.Global
	}

	if taskList.Len() < 1 {
		return nil, errors.New("task is empty. No tasks to sort")
	}

	if err := taskList.Sort(todotxt.SortCompletedDateAsc); err != nil {
		return nil, errors.Wrap(err, "failed to soft by completed date")
	}

	return taskList, nil
}

func (c *Command) isQueryReady() bool {
	return !c.AppInfo.Config.IsDefaultConf()
}

// SurvaySave は質問集が存在しない場合、ユーザに問い合わせて設定ファイルを新規保存します。
func (c *Command) SurvaySave() error {
	// アプリの設定ファイルの更新
	if c.AppInfo.Config.FileUsed() == "" && c.AppInfo.Tasks.Local.FileUsed() != "" {
		msg := "設定ファイルが見つかりません。タスクと同じ階層に作成しますか？"

		yes, err := c.CUI.Confirm(msg)

		switch {
		case yes && err == nil:
			pathDirConf := filepath.Dir(c.AppInfo.Tasks.Local.FileUsed())

			if !strings.Contains(pathDirConf, ".qiitask") {
				pathDirConf = filepath.Join(pathDirConf, ".qiitask")
			}

			pathFileConf := filepath.Join(pathDirConf, config.NameConf)

			return c.AppInfo.Config.WriteConfigAs(pathFileConf)
		case err != nil:
			return errors.Wrap(err, "error during confirmation")
		default:
			fmt.Println("保存しませんでした。（質問タイプを毎回選択する必要があります）")
		}

		c.CUI.DrawHR()
	}

	return nil
}

// Sort は cmdsort の本体です。タスクを対話式でソートします。
func (c *Command) Sort(cmd *cobra.Command, args []string) error {
	// 質問集の読み込み準備
	if !c.isQueryReady() {
		if err := c.SurvayQueryType(); err != nil {
			return err
		}

		if err := c.SurvaySave(); err != nil {
			return err
		}
	}

	answers, err := c.GetTaks()
	if err != nil {
		return errors.Wrap(err, "fail to get task list")
	}

	indexQ := 0

	answers.CustomSort(func(a int, b int) bool {
		var result bool

		// Both A and B is undone so ask user which is prior
		A, _ := answers.GetTodoByKey(a)
		B, _ := answers.GetTodoByKey(b)

		switch {
		case answers.IsKeyDone(a) && answers.IsKeyDone(b):
			// Both task is done so keep the order
			result = true
		case answers.IsKeyDone(a) && !answers.IsKeyDone(b):
			// Task A is done but B is not so B is prior
			result = false
		case !answers.IsKeyDone(a) && answers.IsKeyDone(b):
			// Task A is undone but B is done so A is prior
			result = true
		default:
			result = c.askIsALessThanB(A, B, indexQ)
		}

		return result
	})

	return answers.OverWrite(c.CUI)
}

// SurvayQueryType はユーザに質問集のタイプ（タスク整理向け、コレクション整理向
// けなど）を問い合わせ、アプリの設定ファイルの queries フィールドにセットします。
//
// このメソッドは保存は行いません。呼び出し側で必要にあわせて保存処理を別途実行
// する必要があります。
func (c *Command) SurvayQueryType() error {
	oldTimeout := c.CUI.Timeout
	defer func() {
		c.CUI.Timeout = oldTimeout
	}()

	c.CUI.Timeout = 0

	listKey := query.List()
	helpDescription := util.HereDoc(`
		アプリの設定ファイルが作成されていないため、デフォルトの質問集を使います。
		現在のタスクの種類にマッチした質問のタイプを選択してください。

		■ 質問タイプの説明:
	`)

	for _, key := range listKey {
		desctiption := query.Description(key)

		helpDescription += fmt.Sprintf("%v: %v\n", key, desctiption)
	}

	message := "ソート時に使う質問タイプを選択してください。"

	selectedKey, err := c.CUI.Select(message, listKey, "task", helpDescription)
	if err != nil {
		return errors.Wrap(err, "error during sort process")
	}

	selectedQuery, err := query.New(selectedKey)
	if err != nil {
		return errors.Wrap(err, "error during query type selection")
	}

	c.AppInfo.Config.Set("queries", selectedQuery)

	return nil
}
