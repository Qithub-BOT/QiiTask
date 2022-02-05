package query_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/Qithub-BOT/QiiTask/core/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleDescription() {
	fmt.Println("Description of 'task' query:", query.Description("task"))
	fmt.Println("Description of 'likes' query:", query.Description("likes"))

	// Output:
	// Description of 'task' query: ものごとの整理向けの質問（作業タスクなど）
	// Description of 'likes' query: 好きなものの整理向けの質問（コレクションなど）
}

func ExampleList() {
	list := query.List()

	for i := 0; i < len(list); i++ {
		value := list[i]

		if value == "task" {
			fmt.Println("the list 'task' was found in the embedded data")

			continue
		}

		if value == "likes" {
			fmt.Println("the list 'likes' was found in the embedded data")

			continue
		}
	}

	// Output:
	// the list 'task' was found in the embedded data
	// the list 'likes' was found in the embedded data
}

func ExampleNew() {
	obj, err := query.New("task")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("質問の説明:", obj.Description)
	fmt.Println("客観的な質問[0]:", obj.Objective[0])
	fmt.Println("主観的な質問[0]:", obj.Subjective[0])

	// Output:
	// 質問の説明: ものごとの整理向けの質問（作業タスクなど）
	// 客観的な質問[0]: どちらが重要ですか
	// 主観的な質問[0]: いま手を付けたいのはどちらですか
}

func Test_failed_to_unmarshal(t *testing.T) {
	oldEmbeddedQuery := query.EmbeddedQuery
	defer func() {
		query.EmbeddedQuery = oldEmbeddedQuery
	}()

	// Mock embedded query
	query.EmbeddedQuery = []byte("foo:\nhoge='fuga'")

	// Test for New function
	{
		_, err := query.New("task")

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse default JSON data to Go object")
	}

	// Test for List function
	{
		result := query.List()

		// Assert
		require.Empty(t, result, "if the embedded data is broken it should return empty object")
		assert.IsType(t, []string{}, result, "the returned object should be a slice of string")
	}

	// Test for Description function
	{
		result := query.Description("task")

		// Assert
		require.Empty(t, result, "if the embedded data is broken it should return empty object")
	}
}

func Test_unknown_key(t *testing.T) {
	// Test for New function
	{
		_, err := query.New("unknown")

		// Assert
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown key not found in the default list")
	}

	// Test for Description function
	{
		result := query.Description("unknown")

		// Assert
		require.Empty(t, result, "unknown key should return an empty string")
	}
}
