/*
Package appinfo はアプリの設定情報およびタスクを管理するオブジェクトを定義します。
*/
package appinfo

import (
	"github.com/Qithub-BOT/QiiTask/core/config"
	"github.com/Qithub-BOT/QiiTask/core/todo"
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
//  型の定義
// ----------------------------------------------------------------------------

// AppInfo はアプリの設定および読み込んだタスクを保持する型です。
//
// 生成されたオブジェクトは、アプリの各コマンド（cmd にあるパッケージ）に共有情
// 報として渡されます。
type AppInfo struct {
	Config    *config.Config // アプリの設定ファイル情報
	Tasks     *todo.Set      // ローカルおよびグローバルのタスク情報
	Version   string         // アプリのバージョン情報
	IsVerbose bool           // アプリの詳細表示モード（情報がある場合）
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New は初期化された AppInfo オブジェクトを返します。
//
// 返されたオブジェクトは pathDirCurr, pathDirHome に存在したアプリの設定および
// タスクを読み込んだ状態で返されます。
func New(pathDirCurr, pathDirHome, version string) (*AppInfo, error) {
	conf, err := config.New(pathDirCurr, pathDirHome)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate config object")
	}

	tasks, err := todo.NewSet(pathDirCurr, pathDirHome)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate tasks object")
	}

	app := new(AppInfo)

	app.Config = conf
	app.Tasks = tasks
	app.Version = version
	app.IsVerbose = false

	return app, nil
}
