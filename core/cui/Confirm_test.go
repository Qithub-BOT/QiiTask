package cui_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Qithub-BOT/QiiTask/core/cui"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirm(t *testing.T) {
	// Backup and defer recover
	oldSurveyAskOne := cui.SurveyAskOne
	defer func() {
		cui.SurveyAskOne = oldSurveyAskOne
	}()

	expect := true

	// Mock user selection as "Yes"(true)
	cui.SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		return surveyCore.WriteAnswer(response, "", expect)
	}

	obj := cui.New()

	response, err := obj.Confirm("are you test?")

	require.NoError(t, err)
	assert.True(t, response)
}

func TestConfirm_forced_error(t *testing.T) {
	obj := cui.New()

	// Force error
	obj.ForceError = true

	result, err := obj.Confirm("test select")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "forced error")
	assert.False(t, result, "on error returned value should be false")
}

func TestConfirm_forced_bool(t *testing.T) {
	obj := cui.New()

	{
		// Force true
		obj.ForceTrue = true
		obj.ForceFalse = false

		result, err := obj.Confirm("test select")

		require.NoError(t, err)
		assert.True(t, result, "property ForceTrue=true should return true")
	}
	{
		// Force false
		obj.ForceTrue = false
		obj.ForceFalse = true

		result, err := obj.Confirm("test select")

		require.NoError(t, err)
		assert.False(t, result, "property ForceFalse=true should return true")
	}
}

func TestConfirm_interrupt_by_user(t *testing.T) {
	// Backup and defer recover
	oldSurveyAskOne := cui.SurveyAskOne
	defer func() {
		cui.SurveyAskOne = oldSurveyAskOne
	}()

	// Mock user cancel (ctl+c) by returning error
	cui.SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		return terminal.InterruptErr
	}

	obj := cui.New()

	result, err := obj.Confirm("interrupt me if you can")

	require.Error(t, err)
	require.False(t, result, "on error it shuld return false")
	assert.Contains(t, err.Error(), "選択がキャンセルされました")
}
