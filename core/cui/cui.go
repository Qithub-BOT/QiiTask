/*
Package cui は CUI 操作（描画含む）に関する処理をまとめたものです。
*/
package cui

import (
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

// TimeoutDefault は入力待ちのデフォルトの待ち時間です。（秒）
const TimeoutDefault = 5

// WidthTermDefault はターミナルの画面幅が取得dけいない場合のデフォルトの画面幅です。
const WidthTermDefault = 80

// SurveyAskOne は survey.AskOne() のコピーです。このパッケージは survey.AskOne
// の替わりに SurveyAskOne を利用しています。
//
// テスト中に survey.AskOne の挙動を変えたい場合に別の関数を代入して利用します。
// このパッケージを使う際に、エラーを強制的に返すだけであればオブジェクトの
// ForceError プロパティ （フィールド）を true にセットしてください。
var SurveyAskOne = survey.AskOne

// OsStdout は os.Stdout のコピーです。
//
// os.Stdout を参照するメソッド（TermWidth など）は、代わりに OsStdout を参照し
// ています。テスト時にモックの必要（代替関数を代入して挙動を変える必要）がある
// 場合に利用ください。
var OsStdout = os.Stdout

// OsStdin は os.Stdin のコピーです。
var OsStdin = os.Stdin

// ----------------------------------------------------------------------------
//  Type UI
// ----------------------------------------------------------------------------

// UI は CUI 操作（描画含む）に関するメソッドをまとめたものです。
type UI struct {
	MirrorIO    io.Writer // MirrorIO をセットすると、その IO に出力します。（デフォルト: cui.OsStdout）
	ForceString string    // ForceString を空（""）以外にすると強制的にその値を返します。
	Timeout     int       // Timeout はデフォルト値を返すまでの待機時間（秒）です。（デフォルト: 5）
	ForceInt    int       // ForceInt を 0 以外にすると強制的にその値を返します。
	ForceError  bool      // ForceError を true にすると強制的にエラーを返します。
	ForceTrue   bool      // ForceTrue を true にすると強制的に true を返します。（ForceFalse と併用はできません）
	ForceFalse  bool      // ForceFalse を true にすると強制的に false を返します。（ForceTrue と併用はできません）
}

// New は UI の新規インスタンスを返します。
func New() *UI {
	ui := new(UI)
	ui.MirrorIO = os.Stdout
	ui.Timeout = TimeoutDefault

	return ui
}
