/*
Package cmdlist defines the "list" command.
*/
package cmdlist

import (
	"io"

	"github.com/1set/todotxt"
	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ----------------------------------------------------------------------------
//  Commnad Struct
// ----------------------------------------------------------------------------

// Command is the struct to hold cobra.Command and it's flag options.
type Command struct {
	*cobra.Command
	AppInfo    *appinfo.AppInfo
	styleTable string // flag for "--style" option
	addNoEdit  bool   // flag for "--add-no-edit" option
	isGlobal   bool   // flag for "--global" option
	showAll    bool   // flag for "--all" option
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New returns the newly created object pointer of the "list" command.
func New(appInfo *appinfo.AppInfo) *cobra.Command {
	// Instantiate new object
	cmdList := new(Command)

	// Set command
	cmdList.Command = &cobra.Command{
		Use:   "list [option]",
		Short: "タスクの一覧を表示します",
		Long: util.HereDoc(`
				About:
				  'list' は現在のタスク一覧を表示します。デフォルトで未完成のタスク
				  のみが表示されます。
			`),
		Example: util.HereDoc(`
				qiitask list
				qiitask list --all
			`, "  "),
	}

	// Set app info (conf and tasks)
	cmdList.AppInfo = appInfo

	// Assign the method to command's RunE()
	cmdList.Command.RunE = cmdList.List

	// Define flags for `list` command.
	cmdList.Flags().StringVarP(
		&cmdList.styleTable, "style", "s", "text", "テーブルの表示スタイルを指定します。(text, color, markdown, html, csv)",
	)
	cmdList.Flags().BoolVar(
		&cmdList.addNoEdit, "add-no-edit", false, "出力時に DO-NOT-EDIT を追加します。（自動生成された旨を加えます）",
	)
	cmdList.Flags().BoolVarP(
		&cmdList.isGlobal, "global", "g", false, "グローバル・タスクを表示します",
	)
	cmdList.Flags().BoolVarP(
		&cmdList.showAll, "all", "a", false, "完了済みのタスクも表示します",
	)

	return cmdList.Command
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

func (c *Command) getTaskList() *todo.Todo {
	taskList := c.AppInfo.Tasks.Local

	if c.isGlobal || taskList.FileUsed() == "" {
		taskList = c.AppInfo.Tasks.Global
	}

	return taskList
}

func (c *Command) drawTable(mirror io.Writer, taskList *todo.Todo) {
	var (
		appendSeparator bool
		style           cui.TableStyle
		tableTmp        = table.NewWriter()
		ui              = cui.New()
	)

	ui.MirrorIO = mirror

	switch c.styleTable {
	case "color":
		style = cui.AsColoredTable
		appendSeparator = false
	case "markdown":
		style = cui.AsMarkdownTable
	case "html":
		style = cui.AsHTMLTable
	case "csv":
		style = cui.AsCSVTable
	default:
		style = cui.AsDefaultTable // same as text
		appendSeparator = true
	}

	// テーブルデータ作成
	tableTmp.AppendHeader(table.Row{"#", "title"})

	tmpTask := taskList.Filter(todotxt.FilterNotCompleted) // 未完成タスクのみ
	intSeparator := c.AppInfo.Config.GetInt("separator_interval")

	for i := range tmpTask {
		if i%intSeparator == 0 && appendSeparator {
			tableTmp.AppendSeparator() // 区切り線の挿入
		}

		tableTmp.AppendRow([]interface{}{
			[]todotxt.Task(tmpTask)[i].ID, []todotxt.Task(tmpTask)[i].Todo,
		})
	}

	ui.DrawTable(tableTmp, style) // テーブルの描画
}

// List は "list" コマンドの本体です。
func (c *Command) List(cmd *cobra.Command, args []string) error {
	taskList := c.getTaskList()

	if taskList.Len() < 1 {
		return errors.Errorf("まだタスクはありません")
	}

	if c.addNoEdit {
		switch c.styleTable {
		case "text":
			cmd.Println("// Code generated by QiiTask; DO NOT EDIT.")
		case "color":
			cmd.Println("// Code generated by QiiTask; DO NOT EDIT.")
		case "markdown":
			cmd.Println("<!-- // Code generated by QiiTask; DO NOT EDIT. -->")
		case "html":
			cmd.Println("<!-- // Code generated by QiiTask; DO NOT EDIT. -->")
		}
	}

	c.drawTable(cmd.OutOrStdout(), taskList)

	return nil
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------
