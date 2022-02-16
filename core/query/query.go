/*
Package query はソート時に使われるデフォルトの質問集です。

"query.json" を Unmarshal したものを返すだけのパッケージですが、この JSON ファイ
ルはバイナリに埋め込まれます。

"query.json" を最新版に更新する場合は `go generate ./...` を実行してください。（要 `curl`）

また、この "query.json" ファイルは https://git.io/JMPFq からダウンロードして埋め
込まれるため、** "query.json" ファイルは直接編集しないでください**。以下のリポジ
トリに PR してください。

   https://github.com/KEINOS/QiiTaskQuery

*/
package query

import (
	_ "embed" // 下記 //go:embed を機能させるために _ で読み込みだけさせておく。
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
)

// Query は各質問一覧を定義する構造体です。
type Query struct {
	Description string   `json:"description,omitempty" mapstructure:"description"`
	Objective   []string `json:"客観的な質問,omitempty" mapstructure:"客観的な質問"`
	Subjective  []string `json:"主観的な質問,omitempty" mapstructure:"主観的な質問"`
}

// Samples は同梱されている全ての質問一覧を保持する構造体です。
type Samples map[string]Query

// EmbeddedQuery はバイナリに埋め込まれたデフォルトの質問データです。
//go:embed query.json
var EmbeddedQuery []byte

// Description は指定された key の質問集の説明文を返します。（query.json の
// "key.description" 要素の値）
//
// 該当する key が存在しない場合は空の文字列を返します。指定可能な key は List()
// 関数で取得してください。
func Description(key string) string {
	var samples Samples

	if err := json.Unmarshal(EmbeddedQuery, &samples); err != nil {
		return ""
	}

	if _, ok := samples[key]; !ok {
		return ""
	}

	return samples[key].Description
}

// List は選択可能な質問集のキー名を返します。
func List() []string {
	var (
		samples Samples
		list    = []string{}
	)

	if err := json.Unmarshal(EmbeddedQuery, &samples); err != nil {
		return list
	}

	for key := range samples {
		list = append(list, key)
	}

	// map の range は順番をランダム変えるので常に固定の順番にするためにソート
	sort.Slice(list, func(i, j int) bool {
		return list[i] > list[j]
	})

	return list
}

// New は指定された key 名の質問集を返します。指定可能な key は List() 関数で取
// 得してください。
func New(key string) (Query, error) {
	var samples Samples

	if err := json.Unmarshal(EmbeddedQuery, &samples); err != nil {
		return Query{}, errors.Wrap(err, "failed to parse default JSON data to Go object")
	}

	if _, ok := samples[key]; !ok {
		return Query{}, errors.Errorf("%v key not found in the default list", key)
	}

	return samples[key], nil
}
