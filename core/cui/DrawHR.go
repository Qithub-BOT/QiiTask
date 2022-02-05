package cui

import (
	"fmt"
	"strings"
)

// DrawHR はテキストの罫線を描画します。
func (ui *UI) DrawHR() {
	width := ui.TermWidth()

	if width == 0 {
		width = WidthTermDefault
	}

	hr := strings.Repeat("-", width)

	fmt.Fprintf(ui.MirrorIO, "%v\n", hr)
	//fmt.Println(hr)
}
