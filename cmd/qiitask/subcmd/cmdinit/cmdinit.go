package cmdinit

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/core/config"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// ----------------------------------------------------------------------------
//  Commnad Struct
// ----------------------------------------------------------------------------

// Command は cobra.Command 型の拡張型です。cobra.Command に加えフラグの設定値を
// 保持するためのフィールドを持ちます。
type Command struct {
	*cobra.Command
	AppInfo   *appinfo.AppInfo
	CUI       *cui.UI
	typeQuery string
}

// ----------------------------------------------------------------------------
//  Variables
// ----------------------------------------------------------------------------

// OsUserHomeDir は os.UserHomeDir のコピーです。テストで os.UserHomeDir の挙動
// を変えたい場合に代替関数を割り当ててモックしてください。
var OsUserHomeDir = os.UserHomeDir

// OsMkdirAll は os.MkdirAll のコピーです。テストで os.MkdirAll の挙動を変えたい
// 場合に代替関数を割り当ててモックしてください。
var OsMkdirAll = os.MkdirAll

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New は "say" コマンドの新規オブジェクト（のポインタ）を返します。
func New(appInfo *appinfo.AppInfo) *cobra.Command {
	// コマンドのインスタンス生成
	cmdInit := new(Command)

	// コマンドの割り当て
	cmdInit.Command = &cobra.Command{
		Use:   "init",
		Short: "グローバルなタスクを作成します。",
		Long: util.HereDoc(`
				About:
				  'init' コマンドは、どのディレクトリにいても確認できるグローバルなタスクをユー
				  ザーのホームディレクトリに作成します。カレント・ディレクトリに todo.txt がある
				  場合は、そちらを優先します。
			`),
		Example: util.HereDoc(`
				qiitask init
				qiitask init --force             // 既存の設定があっても強制的に初期化します
				qiitask init --query-type task   // 質問集を作業の整理向けにセットします
				qiitask init --query-type likes  // 質問集をコレクションの整理向けにセットします
			`, "  "),
	}

	// Set app info (conf and tasks)
	cmdInit.AppInfo = appInfo

	// Add CUI object
	cmdInit.CUI = cui.New()

	// RunE function
	cmdInit.Command.RunE = cmdInit.Init

	// Define flags for `init` command.
	cmdInit.Flags().StringVarP(
		&cmdInit.typeQuery, "query-type", "q", "", "質問集のタイプを指定します（task, likes）",
	)

	return cmdInit.Command
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// Init メソッドは "init" コマンドの本体です。
// ユーザーのホームディレクトリに、デフォルトのタスクと設定ファイルを作成します。
func (c *Command) Init(cmd *cobra.Command, args []string) error {
	/* 設定ディレクトリの取得 */
	pathDirHome, err := OsUserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get user home directory path")
	}

	pathDirConf := filepath.Join(pathDirHome, ".config", "qiitask")

	/* ファイルのパス取得 */
	pathFileConf := filepath.Join(pathDirConf, config.NameConf)
	pathFileTask := filepath.Join(pathDirConf, todo.NameFile)

	// 設定ファイルの確認
	if util.IsFile(pathFileConf) {
		msg := fmt.Sprintf(
			"既存の設定ファイルがすでに存在します。\n    %v\n強制的に初期化しますか？",
			c.AppInfo.Config.FileUsed(),
		)

		if isYes, err := c.CUI.Confirm(msg); err != nil {
			return errors.Wrap(err, "error during confirmation")
		} else if !isYes {
			return errors.New("初期化をキャンセルしました")
		}
	}

	// 既存タスクの確認
	if util.IsFile(pathFileTask) {
		msg := fmt.Sprintf(
			"既存のタスクがすでに存在します。\n    %v\n強制的に初期化しますか？",
			c.AppInfo.Tasks.Global.FileUsed(),
		)

		if isYes, err := c.CUI.Confirm(msg); err != nil {
			return errors.Wrap(err, "error during confirmation")
		} else if !isYes {
			return errors.New("初期化をキャンセルしました")
		}
	}

	return c.ForceInit(cmd, pathDirConf)
}

func (c *Command) ForceInit(cmd *cobra.Command, pathDirConf string) error {
	if err := OsMkdirAll(pathDirConf, 0o755); err != nil {
		return errors.Wrap(err, "failed to create app config directory")
	}

	/* ファイルのパス取得 */
	pathFileConf := filepath.Join(pathDirConf, config.NameConf)
	pathFileTask := filepath.Join(pathDirConf, todo.NameFile)

	if err := c.AppInfo.Config.SaveAs(pathFileConf); err != nil {
		return errors.Wrap(err, "failed to save config file")
	}

	if err := c.AppInfo.Tasks.Global.WriteToPath(pathFileTask); err != nil {
		return errors.Wrap(err, "failed to save task file")
	}

	cmd.Println(fmt.Sprintf(
		"アプリの初期化を行いました。以下のディレクトリにファイルを作成しました。\n    %v\n",
		pathDirConf,
	))

	return nil
}
