# コマンド名の命名ルール

なるべく以下の文法で使えるように努めてください。

```bash
# <> ... 必須, [] ... オプション
qiitask <動詞> <名詞> [--形容詞] [引数]
qiitask <DO> <WHAT> [--HOW] [HOWMUCH]
qiitask <コマンド> [引数] [--フラグ]
```

- 例
    - `qiitask list task --done`
    - `qiitask delete task 3`
    - `qiitask set task --done 4`
    - `qiitask add task --top Adventカレンダーのネタを考える`
    - `qiitask add task --global ミルクを買う`
    - `qiitask hello Qiitadon! --reverse`
