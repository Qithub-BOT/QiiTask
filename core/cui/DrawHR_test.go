package cui_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/stretchr/testify/assert"
)

func ExampleUI_DrawHR() {
	ui := cui.New()

	ui.DrawHR() // Draw horizontal text line

	// Output:
	// --------------------------------------------------------------------------------
}

func TestDrawHR(t *testing.T) {
	ui := cui.New()
	buffer := &bytes.Buffer{}

	ui.MirrorIO = buffer

	ui.DrawHR()

	expect := strings.Repeat("-", cui.WidthTermDefault) + "\n"
	actual := buffer.String()

	assert.Equal(t, expect, actual)
}
