<!-- markdownlint-disable MD001 MD041 MD031 MD033 -->
[![go1.16+](https://img.shields.io/badge/Go-1.16,%2017,%20latest-blue?logo=go)](https://github.com/Qithub-BOT/QiiTask/actions/workflows/go-versions.yml "Supported versions")
[![Go Reference](https://pkg.go.dev/badge/github.com/Qithub-BOT/QiiTask.svg)](https://pkg.go.dev/github.com/Qithub-BOT/QiiTask#section-documentation "Read generated documentation of the app")
[![GitHub Codespaces](https://img.shields.io/badge/Codespaces-compatible-blue?logo=github)](https://github.dev/Qithub-BOT/QiiTask "オンラインで VSCode を起動する")

# QiiTask <sub><sup><sup>alpha</sup></sup></sub>

QiiTask は、**ものごとの優先度をエンジニアらしくソート・アルゴリズムでキメるツール**です。

仕組み自体は単純で、文字列の配列をソートする際の比較処理（`a < b`）を、対話式でユーザーに決めさせるだけです。

```shellsession
$ cat todo.txt

```

## Install

- Via [Homebrew](https://brew.sh/index_ja) ([macOS](https://docs.brew.sh/Installation), [Linux](https://docs.brew.sh/Homebrew-on-Linux), [Windows with WSL2](https://docs.brew.sh/Homebrew-on-Linux))
    ```bash
    brew install qithub-bot/apps/qiitask
    ```

- Via Go Install
    ```bash
    go install "github.com/Qithub-BOT/QiiTask/qiitask@latest"
    ```
- Manual Install (Download binary from releases)
  - [Releases](https://github.com/Qithub-BOT/QiiTask/releases/latest) ページから該当する OS/CPU に合ったアーカイブをダウンロードして、パスの通ったディレクトリ に設置してください。（要実行権限）

## Statuses

[![Test on macOS/Win/Linux](https://github.com/Qithub-BOT/QiiTask/actions/workflows/platform-test.yaml/badge.svg)](https://github.com/Qithub-BOT/QiiTask/actions/workflows/platform-test.yaml)
[![go1.15+](https://github.com/Qithub-BOT/QiiTask/actions/workflows/version-tests.yaml/badge.svg)](https://github.com/Qithub-BOT/QiiTask/actions/workflows/version-tests.yaml)
[![golangci-lint](https://github.com/Qithub-BOT/QiiTask/actions/workflows/golangci-lint.yaml/badge.svg)](https://github.com/Qithub-BOT/QiiTask/actions/workflows/golangci-lint.yaml)
[![codecov](https://codecov.io/gh/Qithub-BOT/qiitask/branch/main/graph/badge.svg?token=kJJSFFNwE3)](https://codecov.io/gh/Qithub-BOT/qiitask "View details on CodeCov.IO")
[![Go Report Card](https://goreportcard.com/badge/github.com/Qithub-BOT/QiiTask)](https://goreportcard.com/report/github.com/Qithub-BOT/QiiTask "View on Go Report Card")
[![CodeQL](https://github.com/Qithub-BOT/QiiTask/actions/workflows/codeQL-analysis.yaml/badge.svg)](https://github.com/Qithub-BOT/QiiTask/actions/workflows/codeQL-analysis.yaml "Vulnerability Scan")

## Note

- Currently WIP

## License

- MIT
- (c) Copyright, [QiiTask Contributors](https://github.com/Qithub-BOT/QiiTask/graphs/contributors).
- (c) Copyright, [QiiTask Query Contributors](https://github.com/KEINOS/QiiTaskQuery/graphs/contributors).
