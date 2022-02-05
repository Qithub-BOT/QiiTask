/*
Package cmdhello は "hello" コマンドを定義します。このコマンドは "say" コマンド
のサブ・コマンド（child）です。
*/
package cmdhello

import (
	"fmt"
	"strings"

	"github.com/KEINOS/go-utiles/util"
	"github.com/spf13/cobra"
)

// ----------------------------------------------------------------------------
//  Commnad Struct
// ----------------------------------------------------------------------------

// Command is the struct to hold cobra.Command and it's flag options.
type Command struct {
	*cobra.Command
	isReverse bool // flag value for "--reverse" option
}

// ----------------------------------------------------------------------------
//  Public Functions
// ----------------------------------------------------------------------------

// New は "hello" コマンドの新規オブジェクト（のポインタ）を返します。
func New() *cobra.Command {
	// コマンドのオブジェクト生成
	cmdHello := new(Command)

	// cobra.Command オブジェクトの割り当て
	cmdHello.Command = &cobra.Command{
		Use:   "hello",
		Short: "挨拶を表示します。",
		Long: util.HereDoc(`
				About:
				  'hello' コマンドは挨拶文を返します。引数がある場合は、それに合
				  わせた挨拶を返します。新規コマンドを作成する際の参考にしてくだ
				  さい。
			`),
		Example: util.HereDoc(`
			qiitask say hello                 // Hello!
			qiitask say hello world           // Hello, world!
			qiitask say hello world --reverse // !dlrow ,olleH
			`,
			"  ", // HereDoc のインデント
		),
	}

	// cobra.Command の RunE メソッドに、sayHelloWorld メソッドを代入
	cmdHello.RunE = cmdHello.sayHello

	// Define flags for `world` command.
	cmdHello.Flags().BoolVarP(
		&cmdHello.isReverse, "reverse", "r", cmdHello.isReverse, "Reverses the output.",
	)

	return cmdHello.Command
}

// ----------------------------------------------------------------------------
//  Private Functions
// ----------------------------------------------------------------------------

// reverseString reverses/flip the input string.
func reverseString(input string) string {
	var msgTmp string

	for _, v := range input {
		msgTmp = string(v) + msgTmp
	}

	return msgTmp
}

// ----------------------------------------------------------------------------

// sayHello は "hello" コマンドの本体です。
// デフォルトで "Hello, world!" を表示しますが、引数がある場合は、その引数に挨拶
// します。また "--reverse" オプションが指定されていた場合は、出力の前後を反転さ
// せたものを表示します。
func (c *Command) sayHello(cmd *cobra.Command, args []string) error {
	msgToGreet := "Hello!"

	if len(args) > 0 {
		names := strings.Join(args, " and ") // 結合
		msgToGreet = "Hello, " + names + "!"
	}

	if c.isReverse {
		msgToGreet = reverseString(msgToGreet)
	}

	fmt.Println(msgToGreet)

	return nil
}
