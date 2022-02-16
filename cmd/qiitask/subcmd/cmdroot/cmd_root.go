package cmdroot

import (
	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdinit"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdlist"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsay"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsort"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/spf13/cobra"
)

// ----------------------------------------------------------------------------
//  Commnad Struct
// ----------------------------------------------------------------------------

// Command is the struct to hold cobra.Command and it's flag options.
type Command struct {
	*cobra.Command
}

// ----------------------------------------------------------------------------
//  Public Functions
// ----------------------------------------------------------------------------

// New は各コマンドの親コマンド（root）の新規オブジェクト（のポインタ）を返します。
//
// 引数の appInfo オブジェクトは、各コマンドにも渡されます。
func New(appInfo *appinfo.AppInfo) *cobra.Command {
	if appInfo.Version == "" {
		appInfo.Version = "(unknown)"
	}

	cmdRoot := &Command{
		&cobra.Command{
			Use:   "qiitask",
			Short: "ものごとの優先度をソートアルゴリズムを使ってキメるツール",
			Long: util.HereDoc(`
				About:
				  QiiTask は "todo.txt" メソッドで書かれたタスクを管理するツールです。
				  一番の特徴はタスクの優先度を決めるソート・モードを持っていることで
				  す。対話式に 2 択の質問に答えていくことで、ソートのアルゴリズムを利
				  用して優先度をキメることが期待できます。
			`),
			Example: util.HereDoc(`
				$ qiitask --version
				$ qiitask help
				$ qiitask list

				■ コマンドの中にはサブコマンドがあるものもあります。以下の hello は say
				のサブコマンドです。

				$ qiitask say hello

				■ コマンドのヘルプ表示

				$ qiitask list -h .......... list コマンドのヘルプを表示します。
				$ qiitask help list ........ list コマンドのヘルプを表示します。
				$ qiitask say hello -h ..... hello サブコマンドのヘルプを表示します。
				$ qiitask help say hello ... hello サブコマンドのヘルプを表示します。

				■ bash などのシェル補完スクリプト（completion）を出力することもできます。

				$ qiitask completion --help
				$ qiitask completion bash > qiitask_completion.sh`,
				"  ", // HereDoc のインデント
			),
			Version: appInfo.Version,
		},
	}

	// Define persistent flags for the app.
	cmdRoot.PersistentFlags().BoolVar(&appInfo.IsVerbose, "verbose", false, "displays debug info if any")

	// Add child commands to the "root" command.
	cmdRoot.AddCommand(
		cmdsay.New(appInfo),  // Add "say" command (with grand child command "hello")
		cmdlist.New(appInfo), // Add "list" command
		cmdsort.New(appInfo), // Add "sort" command
		cmdinit.New(appInfo), // Add "init" command
	)

	return cmdRoot.Command
}
