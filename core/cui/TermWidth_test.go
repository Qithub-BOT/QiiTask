package cui_test

import (
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/stretchr/testify/assert"
)

func TestTermWidth(t *testing.T) {
	ui := cui.New()

	result := ui.TermWidth()

	assert.Zero(t, result)
}
