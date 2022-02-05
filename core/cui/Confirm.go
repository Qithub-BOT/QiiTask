package cui

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/pkg/errors"
)

// Confirm は Yes/No 系の確認をします。Yes の場合は true が返されます。
// また、問い合せ中に ctrl+c （SIGINT）が送られてきた場合はエラーが返されます。
func (ui *UI) Confirm(msg string) (isYes bool, err error) {
	if ui.ForceError {
		return false, errors.New("forced error")
	}

	if ui.ForceFalse {
		return false, nil
	}

	if ui.ForceTrue {
		return true, nil
	}

	var userChoice bool

	prompt := &survey.Confirm{
		Message: msg,
	}

	if err := SurveyAskOne(prompt, &userChoice); err != nil {
		msgError := "unknown error occurred during confirmation"

		if err == terminal.InterruptErr {
			msgError = "選択がキャンセルされました"
		}

		return false, errors.Wrap(err, msgError)
	}

	return userChoice, nil
}
