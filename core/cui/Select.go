package cui

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  Methods
// ----------------------------------------------------------------------------

// AskOne は survey.AskOne のラッパーですがタイマー付きです。
// 呼び出されてから、UI.Timeout 秒経過した場合は強制的に ansDefault の値が返さ
// れます。UI.Timeout = 0 の場合はタイマーは発動しません。
//
// タイムアウト時に Ctrl+c を送信させたい場合は以下を ansDefault に設定します。
//
//     ansDefault = string([]rune{terminal.KeyInterrupt})
func (ui *UI) AskOne(prompt survey.Prompt, ansDefault string) (string, error) {
	var selected string

	if err := SurveyAskOne(prompt, &selected); err != nil {
		msgError := "unknown error occurred during query list selection"

		if err == terminal.InterruptErr {
			msgError = "選択がキャンセル（ctrl+c）されました"
		}

		return "", errors.Wrap(err, msgError)
	}

	return selected, nil
}

// Select はユーザー選択の UI です。戻り値は selection のうち選択された値です。
//
// cui.UI.Timeout が 1 以上にセットされていた場合、UI.Timeout 秒経過すると
// ansDefault が強制的に返されます。
func (ui *UI) Select(msg string, selection []string, ansDefault string, helpMsg string) (string, error) {
	if ui.ForceError {
		return "", errors.New("forced error")
	}

	if ui.ForceString != "" {
		return ui.ForceString, nil
	}

	prompt := &survey.Select{
		Message: msg,
		Options: selection,
		Help:    helpMsg,
	}

	return ui.AskOne(prompt, ansDefault)
}
