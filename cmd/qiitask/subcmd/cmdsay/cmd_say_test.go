package cmdsay_test

import (
	"testing"

	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdroot"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsay"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Test Functions
// ----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	obj1 := cmdsay.New(nil)
	obj2 := cmdsay.New(nil)

	assert.NotSame(t, obj1, obj2, "it should not reference the same object")
}

func TestNew_has_child(t *testing.T) {
	mother := cmdsay.New(nil)
	children := mother.Commands()
	expectChild := cmdsay.New(nil)

	require.True(t, mother.HasSubCommands(), "command 'say' should contain a sub-command")
	assert.IsType(t, expectChild, children[0], "command 'say' should contain 'world'")
}

func Test_say(t *testing.T) {
	for _, test := range []struct {
		expect string
		args   []string
	}{
		{
			"せいっ！", []string{},
		},
		{
			"！っいせ", []string{"--reverse"},
		},
		{
			"About:", []string{"--help"},
		},
		{
			"foo and bar!", []string{"foo", "bar"},
		},
		{
			"!olleH", []string{"hello", "--reverse"},
		},
		{
			"Hello, world!", []string{"hello", "world"},
		},
	} {
		args := append([]string{"say"}, test.args...)

		appInfo, err := appinfo.New("", "", "")
		if err != nil {
			t.Fatal(err)
		}

		mother := cmdroot.New(appInfo)
		mother.SetArgs(args)

		out := capturer.CaptureOutput(func() {
			err := mother.Execute()

			require.NoError(t, err,
				"Input args: %v\nExpect out: %v\n", test.args, test.expect)
		})

		assert.Contains(t, out, test.expect,
			"Expect: %v\nActual: %v\n", test.expect, out)
	}
}
