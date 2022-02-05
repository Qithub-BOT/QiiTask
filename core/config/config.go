/*
Package config はアプリの設定（"config.json"）を管理するオブジェクトのパッケージです。

このオブジェクトには、ソートなどに使われる質問、質問の表示時間などの情報が含まれます。
*/
package config

import (
	"path/filepath"

	"github.com/Qithub-BOT/QiiTask/core/query"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ----------------------------------------------------------------------------
//  Constants
// ----------------------------------------------------------------------------

// NameConf はアプリ設定のファイル名です。
const NameConf = "config.json"

// ----------------------------------------------------------------------------
//  型の定義
// ----------------------------------------------------------------------------

// Config は viper.Viper 型（github.com/spf13/viper）を埋め込んだ、アプリの設定
// 情報を保持する拡張型です。
type Config struct {
	*viper.Viper
	// 拡張プロパティ
	pathDir      string // 設定ファイルの読み込み・保存先のディレクトリ
	pathFile     string // 設定ファイルが読み込まれいた場合のファイルのパス
	triedLoading bool   // Load() を実行した場合 true
}

// ----------------------------------------------------------------------------
//  Constructor
// ----------------------------------------------------------------------------

// New はデフォルトの設定済み新規 Config インスタンスを返します。
//
// 引数はファイルの検索先のパスとして指定する必要があります。
func New(pathDirCurr, pathDirHome string) (*Config, error) {
	config := new(Config)
	config.Viper = viper.New()

	// 検索対象のファイル名をセット
	config.SetConfigName(NameConf)
	config.SetConfigType("json") // この設定は必須

	// 設定ファイルの検索パスの追加
	config.AddConfigPath(
		filepath.Join(pathDirCurr, ".qiitask"), // カレント
	)
	config.AddConfigPath(
		filepath.Join(pathDirHome, ".config", "qiitask"), // ユーザホーム
	)

	// デフォルトの保存先ディレクトリをユーザ・ホームにセット。
	// この値は設定ファイルを見つけた場合に上書きされます。
	config.pathDir = filepath.Join(pathDirHome, ".config", "qiitask")

	config.SetDefault("standby_interval", 5)       // 入力待ち時間（秒）
	config.SetDefault("separator_interval", 5)     // リストの区切り位置（行）
	config.SetDefault("queries", new(query.Query)) // ソートに使う質問集

	if err := config.Load(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize app config")
	}

	return config, nil
}

// ----------------------------------------------------------------------------
//  拡張 Methods
// ----------------------------------------------------------------------------

// FileUsed は読み込まれている設定ファイルのパスを返します。
func (c *Config) FileUsed() string {
	return c.pathFile
}

// GetQueryDescription は質問一覧の説明文を返します。
func (c *Config) GetQueryDescription() string {
	result := ""

	if queryInterface, ok := c.Get("queries").(query.Query); ok {
		result = queryInterface.Description
	}

	return result
}

// GetQueryObjective は質問一覧のうち、客観的な質問を返します。
func (c *Config) GetQueryObjective() []string {
	result := []string{""}

	if queryInterface, ok := c.Get("queries").(query.Query); ok {
		result = queryInterface.Objective
	}

	return result
}

// GetQuerySubjective は質問一覧のうち、主観的な質問を返します。
func (c *Config) GetQuerySubjective() []string {
	result := []string{""}

	if queryInterface, ok := c.Get("queries").(query.Query); ok {
		result = queryInterface.Subjective
	}

	return result
}

func (c *Config) GetQuery() *query.Query {
	if queryInterface, ok := c.Get("queries").(query.Query); ok {
		return &queryInterface
	}

	return nil
}

// IsDefaultConf は設定ファイルが存在せずデフォルトの設定を使っている場合に true
// を返します。
func (c *Config) IsDefaultConf() bool {
	return c.pathFile == "" || c.GetQuery() != nil
}

// Load はアプリの設定ファイルを読み込みます。
//
// ファイルが存在するも、読み込みに失敗した場合に error を返します。ファイルが
// 存在しない場合はデフォルトの設定を読み込みます。
func (c *Config) Load() error {
	c.triedLoading = true

	if err := c.ReadInConfig(); err != nil {
		// Ignore error if the error was file-not-found and use default
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.pathFile = ""

			return nil
		}

		// Wrap error if err was other than file-not-found
		return errors.Wrap(err, "config file was found but another error was produced")
	}

	var queryStruct query.Query

	queryInterface := c.Get("queries")
	if err := mapstructure.Decode(queryInterface, &queryStruct); err != nil {
		return errors.Wrap(err, "failed to unmarshal. Config file was found but the structure is not compatible")
	}

	c.Set("queries", queryStruct) // "queries" 要素に query.Query 型で再割り当て
	c.pathFile = c.ConfigFileUsed()
	c.pathDir = filepath.Dir(c.ConfigFileUsed())

	return nil
}

// OverWrite は設定ファイルを上書きします。
func (c *Config) OverWrite() error {
	if c.pathFile == "" {
		return errors.New("config file path not set. Use SaveAs method instead")
	}

	return c.WriteConfig()
}

// SaveAs は設定ファイルを保存します。
func (c *Config) SaveAs(pathFile string) error {
	return c.WriteConfigAs(pathFile)
}
