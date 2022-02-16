package cui

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// TableStyle はテーブル出力時のスタイル ID 用の enum です。
type TableStyle int

const (
	// AsDefaultTable はテーブルをデフォルトの設定で出力するスタイル ID です。
	AsDefaultTable TableStyle = iota
	// AsSimpleTable はテーブルを ASCII 文字で出力するスタイル ID です。
	AsSimpleTable
	// AsColoredTable はテーブルを色付きで出力するスタイル ID です。
	AsColoredTable
	// AsMarkdownTable はテーブルを Markdown で出力するスタイル ID です。
	AsMarkdownTable
	// AsCSVTable はテーブルを CSV で出力するスタイル ID です。
	AsCSVTable
	// AsHTMLTable はテーブルを HTML で出力するスタイル ID です。
	AsHTMLTable
)

// DrawTable はテーブルのテキスト描画を行います。デフォルトで標準出力に出力します。
// 別の I/O に出力したい場合は、オブジェクトの MirrorIO に Writer をセットする必要
// があります。
func (ui *UI) DrawTable(tblToWrite table.Writer, style TableStyle) {
	if ui.MirrorIO != nil {
		tblToWrite.SetOutputMirror(ui.MirrorIO)
	}

	switch style {
	case AsColoredTable:
		tblToWrite.SetStyle(ui.getStyleTableColored())
	default:
		tblToWrite.SetStyle(ui.getStyleTableDefault())
	}

	switch style {
	case AsCSVTable:
		_ = tblToWrite.RenderCSV()
	case AsMarkdownTable:
		_ = tblToWrite.RenderMarkdown()
	case AsHTMLTable:
		_ = tblToWrite.RenderHTML()
	default:
		_ = tblToWrite.Render()
	}
}

// getStyleTableColored は色付きテーブルのスタイルを返します。
func (ui *UI) getStyleTableColored() table.Style {
	style := table.StyleColoredBlueWhiteOnBlack

	style.Color.Header = text.Colors{text.BgBlue, text.FgWhite}
	style.Color.Row = text.Colors{text.BgBlack, text.FgHiWhite}
	style.Color.RowAlternate = text.Colors{text.BgHiBlack, text.FgWhite}

	return style
}

// getStyleTableDefault はデフォルトのテーブル描画スタイルを返します（現在 ASCII フォーマット）。
func (ui *UI) getStyleTableDefault() table.Style {
	return table.StyleDefault
}
