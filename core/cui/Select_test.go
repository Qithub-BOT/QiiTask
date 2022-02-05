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

func TestSelect(t *testing.T) {
	// Backup and defer recover
	oldSurveyAskOne := cui.SurveyAskOne
	defer func() {
		cui.SurveyAskOne = oldSurveyAskOne
	}()

	expect := "foo"

	// Mock user selection as "foo"
	cui.SurveyAskOne = func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
		return surveyCore.WriteAnswer(response, "", expect)
	}

	obj := cui.New()

	options := []string{
		expect,
		"bar",
		"buzz",
	}

	actual, err := obj.Select("test select", options, "dummy", "no help message")

	require.NoError(t, err)
	assert.Equal(t, expect, actual, "it should return the user selected string")
}

func TestSelect_forced_error(t *testing.T) {
	obj := cui.New()

	// Force error
	obj.ForceError = true

	options := []string{
		"dummy",
		"selection",
	}

	result, err := obj.Select("test select", options, "dummy", "no help message")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "forced error")
	assert.Empty(t, result, "on error return value should be empty")
}

func TestSelect_forced_string(t *testing.T) {
	obj := cui.New()

	expect := "this is a forced answer"

	// Force error
	obj.ForceString = expect

	options := []string{
		"dummy",
		"selection",
	}

	result, err := obj.Select("test select", options, "dummy", "no help message")

	require.NoError(t, err)
	assert.Equal(t, expect, result)
}

func TestSelect_interrupt_by_user(t *testing.T) {
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

	options := []string{
		"dummy",
		"selection",
	}

	result, err := obj.Select("test select", options, "dummy", "no help message")

	require.Error(t, err)
	require.Empty(t, result, "on error it shuld return empty")
	assert.Contains(t, err.Error(), "選択がキャンセル（ctrl+c）されました")
}
