/*
Package cmdsay は "say" コマンドを定義します。
*/
package cmdsay

import (
	"strings"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsay/cmdhello"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
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
	isReverse bool // flag for "--reverse" option
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New は "say" コマンドの新規オブジェクト（のポインタ）を返します。
func New(appInfo *appinfo.AppInfo) *cobra.Command {
	// コマンドのインスタンス生成
	cmdSay := new(Command)

	// コマンドの割り当て
	cmdSay.Command = &cobra.Command{
		Use:   "say [message ...]",
		Short: "引数をオウム返しで表示します。",
		Long: util.HereDoc(`
				About:
				  'say' コマンドは引数をオウム返しするだけのコマンドです。コマン
				  ド作成時の参考にしてください。
			`),
		Example: util.HereDoc(`
				qiitask say foo bar           // foo and bar!
				qiitask say --reverse foo bar // !rab dna oof
			`, "  "),
	}

	// Set app info (conf and tasks)
	cmdSay.AppInfo = appInfo

	// RunE function
	cmdSay.Command.RunE = cmdSay.say

	// Define flags for `say` command.
	cmdSay.Flags().BoolVarP(
		&cmdSay.isReverse, "reverse", "r", false, "Reverses the output.",
	)

	// "say" コマンドのサブ・コマンドとして "hello" コマンドを追加
	cmdSay.AddCommand(cmdhello.New())

	return cmdSay.Command
}

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// say メソッドは "say" コマンドの本体です。
// 引数をオウム返しするだけですが、複数の引数がある場合は "and" でつないで返しま
// す。`--reverse` オプションが指定されていた場合は、文字列を前後反転させたもの
// を返します。
func (c *Command) say(cmd *cobra.Command, args []string) error {
	msgToGreet := "せいっ！" // デフォルトの返事

	// 引数からメッセージ内容を作成
	if len(args) > 0 {
		names := strings.Join(args, " and ") // 結合
		msgToGreet = names + "!"
	}

	// `--reverse` オプションが指定されていた場合の処理
	if c.isReverse {
		msgToGreet = reverseString(msgToGreet)
	}

	cmd.Println(msgToGreet) // fmt.Println と同等ですが、cmd 経由で出力します。

	return nil
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

// reverseString 関数は input の値を前後反転させた値を返します。
func reverseString(input string) string {
	var msgTmp string

	for _, v := range input {
		msgTmp = string(v) + msgTmp
	}

	return msgTmp
}
