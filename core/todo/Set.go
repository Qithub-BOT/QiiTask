package todo

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  型の定義
// ----------------------------------------------------------------------------

// Set はローカルおよびグローバルのタスクをセットで保持するための型です。
type Set struct {
	Local  *Todo // Local: カレントディレクトリにあるタスク
	Global *Todo // Global: ユーザのホームディレクトリにあるタスク
}

// ----------------------------------------------------------------------------
//  Global Variables
// ----------------------------------------------------------------------------

// FilePathAbs は filepath.Abs のコピーです。テスト時に filepath.Abs の動作をモック
// する為に変数に代入しています。
var FilePathAbs = filepath.Abs

// OsGetwd は os.Getwd のコピーです。テスト時に os.Getwd の動作をモックする為に変数
// に代入しています。
var OsGetwd = os.Getwd

// OsUserHomeDir は os.UserHomeDir のコピーです。テスト時に os.UserHomeDir の動作を
// モックする為に変数に代入しています。
var OsUserHomeDir = os.UserHomeDir

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// NewSet は新規 Set インスタンスを返します。Set はカレントディレクト
// リのタスク（Local タスク）およびユーザーのホームディレクトリ （Global タスク）
// を含みます。
//
// 引数はファイルの検索先のパスとして指定する必要があります。
func NewSet(pathDirCurr, pathDirHome string) (*Set, error) {
	var err error

	pathDirCurrAbs, err := normalizePathDirCurr(pathDirCurr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path")
	}

	pathDirHomeAbs, err := normalizePathDirHome(pathDirHome)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get absolute path")
	}

	list := new(Set)

	if list.Local, err = New(pathDirCurrAbs); err != nil {
		return nil, err
	}

	if list.Global, err = New(pathDirHomeAbs); err != nil {
		return nil, err
	}

	return list, nil
}

func normalizePathDirCurr(path string) (string, error) {
	var err error

	if path == "" || path == "." || path == "./" {
		path, err = OsGetwd()
		if err != nil {
			return "", errors.Wrap(err, "failed to get current working dir")
		}
	}

	return FilePathAbs(path)
}

func normalizePathDirHome(path string) (string, error) {
	if path == "" || path == "~" || path == "~/" {
		var err error

		path, err = OsUserHomeDir()
		if err != nil {
			return "", errors.Wrap(err, "failed to get user home directory")
		}
	}

	return FilePathAbs(path)
}
