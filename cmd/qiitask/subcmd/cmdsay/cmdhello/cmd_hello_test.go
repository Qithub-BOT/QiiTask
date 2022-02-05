package cmdhello_test

import (
	"testing"

	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdroot"
	"github.com/Qithub-BOT/QiiTask/cmd/qiitask/subcmd/cmdsay/cmdhello"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
//  Test Functions
// ----------------------------------------------------------------------------

func TestNew(t *testing.T) {
	obj1 := cmdhello.New()
	obj2 := cmdhello.New()

	assert.NotSame(t, obj1, obj2, "it should not reference the same object")
}

func Test_sayHello(t *testing.T) {
	appInfo, err := appinfo.New("", "", "")
	require.NoError(t, err)

	for _, test := range []struct {
		expect string
		args   []string
	}{
		{"Hello!", []string{}},
		{"!olleH", []string{"-r"}},
		{"!olleH", []string{"--reverse"}},
		{"Hello, foo and bar!", []string{"foo", "bar"}},
		{"!rab dna oof ,olleH", []string{"foo", "bar", "--reverse"}},
	} {
		args := append([]string{"say", "hello"}, test.args...)

		grandMother := cmdroot.New(appInfo)
		grandMother.SetArgs(args)

		out := capturer.CaptureOutput(func() {
			err := grandMother.Execute()

			require.NoError(t, err)
		})

		assert.Contains(t, out, test.expect)
	}
}
