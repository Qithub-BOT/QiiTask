package cui

import (
	"golang.org/x/term"
)

// TermWidth はターミナルの横幅の文字数（バイト数）を返します。取得できない場合
// は 0 を返します。
func (ui *UI) TermWidth() int {
	var result int

	result, _, err := term.GetSize(int(OsStdout.Fd()))
	if err != nil {
		result = 0 // 取得できなかった場合はデフォルト値をセット
	}

	return result
}
