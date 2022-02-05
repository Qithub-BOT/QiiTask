package cui_test

import (
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	obj1 := cui.New()
	obj2 := cui.New()

	assert.NotSame(t, obj1, obj2, "two pointers should not reference the same object")
}
