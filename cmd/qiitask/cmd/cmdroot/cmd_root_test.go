package cmdroot_test

import (
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/qiitask/cmd/cmdroot"
	"github.com/Qithub-BOT/QiiTask/qiitask/cmd/cmdsay"
	"github.com/Qithub-BOT/QiiTask/qiitask/cmd/cmdsay/cmdhello"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	appInfo, err := appinfo.New("", "", "")
	require.NoError(t, err)

	obj1 := cmdroot.New(appInfo)
	obj2 := cmdroot.New(appInfo)

	assert.NotSame(t, obj1, obj2, "it should not reference the same object")
}

func TestNew_empty_version(t *testing.T) {
	appInfo, err := appinfo.New("", "", "")
	require.NoError(t, err)

	mother := cmdroot.New(appInfo)

	expect := "(unknown)"
	actual := mother.Version

	assert.Equal(t, expect, actual)
}

func TestNew_has_child(t *testing.T) {
	appInfo, err := appinfo.New("", "", "")
	require.NoError(t, err)

	mother := cmdroot.New(appInfo)

	// Assertion children
	require.True(t, mother.HasSubCommands(),
		"new root object should contain sub-command/s")

	// Assertion of child command
	expectChild := cmdsay.New(appInfo)
	actualChild, _, err := mother.Find([]string{"say"})

	require.NoError(t, err, "root object should contain 'say' command")
	require.IsType(t, expectChild, actualChild,
		"command 'root' should contain 'cmdsay.Command' type object as a child")

	// Assertion of grand child command
	expectGrandChild := cmdhello.New()
	actualGrandChild, _, err := actualChild.Find([]string{"hello"})

	require.NoError(t, err, "'say' object should contain 'hello' command")
	assert.IsType(t, expectGrandChild, actualGrandChild,
		"command 'root' should contain 'hello' as a grandchild")
}
