/*
Package todo は "todo.txt" 形式のタスクを管理するためのパッケージです。

このパッケージは todotxt.TaskList 型（github.com/1set/todotxt）を拡張しています。

    QiiTask 独自のルール

    todo.txt の "(A)", "(B)" の優先度タグは、QiiTask では以下の意味を持ちます。

        (A) .... 部分ソート済み（主観的な質問でソートされたタスク、Subjective）
	    (B) .... 全体ソート済み（客観的な質問でソートされたタスク、Objective）
	    なし ... 未ソートのタスク
*/
package todo
