package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetQuery(t *testing.T) {
	// Empty conf
	{
		pathDirTest := filepath.Join(
			util.GetPathDirRepo(),
			"testdata",
			"golden",
			"working_task_under_conf",
		)

		conf, err := config.New(pathDirTest, t.TempDir())
		require.NoError(t, err)

		// Test
		result := conf.GetQuery()

		// Assertions
		require.Nil(t, result, "if conf file is empty it should return an error")
	}

	// Golden conf
	{
		pathDirTest := filepath.Join(
			util.GetPathDirRepo(),
			"testdata",
			"golden",
			"working_with_task_and_conf",
		)

		conf, err := config.New(pathDirTest, t.TempDir())
		require.NoError(t, err)

		// Test
		result := conf.GetQuery()

		// Assertions
		require.NotNil(t, result)
		assert.Equal(t, "これはテスト用の質問集です", result.Description)
	}
}

func TestLoad(t *testing.T) {
	pathDirTest := filepath.Join(util.GetPathDirRepo(), "testdata", "golden", "working_with_task_and_conf")
	pathDirConf := filepath.Join(pathDirTest, ".qiitask")

	conf, err := config.New(pathDirTest, t.TempDir())
	require.NoError(t, err)

	assert.NotEmpty(t, conf.FileUsed(), "it should load an existing conf file at: %v", pathDirConf)
	assert.Equal(t, int(10), conf.GetInt("separator_interval"), "it should return the value of the conf file")
	assert.Equal(t, int(15), conf.GetInt("standby_interval"), "it should return the value of the conf file")

	// 客観的な質問の読み込みテスト
	queryObjective := conf.GetQueryObjective()
	require.NotEmpty(t, queryObjective, "failed to load objective query")
	assert.Equal(t, "どちらが重要ですか", queryObjective[0])

	// 主観的な質問の読み込みテスト
	querySubjective := conf.GetQuerySubjective()
	require.NotEmpty(t, querySubjective, "failed to load subjective query")
	assert.Equal(t, "いま手を付けたいのはどちらですか", querySubjective[0])

	// 質問集の説明文のテスト
	queryDescription := conf.GetQueryDescription()
	assert.Equal(t, "これはテスト用の質問集です", queryDescription)
}

func TestLoad_queries_wrong(t *testing.T) {
	pathDirTest := filepath.Join(util.GetPathDirRepo(), "testdata", "error", "queries_as_wrong")

	_, err := config.New(pathDirTest, t.TempDir())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to initialize app config")
	assert.Contains(t, err.Error(), "failed to unmarshal. Config file was found but the structure is not compatible")
}

func TestLoad_malformed_conf(t *testing.T) {
	pathDirTest := filepath.Join(util.GetPathDirRepo(), "testdata", "error", "malformed_conf")

	conf, err := config.New(pathDirTest, t.TempDir())
	require.Nil(t, conf, "on error the object should be nil")
	require.Error(t, err, "mal-formed conf should return an error")

	assert.Contains(t, err.Error(), "config file was found but another error was produced")
}

func TestNew_default_conf_values(t *testing.T) {
	// instantiate object with empty dir as a target
	conf, err := config.New(t.TempDir(), t.TempDir())

	require.NoError(t, err)
	require.True(t, conf.IsDefaultConf(),
		"on conf file not found, it should use the default conf")

	for _, test := range []struct {
		expectValue interface{}
		key         string
		errMsg      string
	}{
		{5, "standby_interval", "%#v key should contain the default value %#v"},
		{5, "separator_interval", "%#v key should contain the default value %#v"},
	} {
		expect := test.expectValue
		actual := conf.Get(test.key)

		assert.Equal(t, expect, actual, test.errMsg, test.key, actual)
	}
}

func TestNew_instantiation(t *testing.T) {
	obj1, err := config.New("", "")
	require.NoError(t, err)

	obj2, err := config.New("", "")
	require.NoError(t, err)

	// We do not use singleton style Viper usage.
	assert.NotSame(t, obj1, obj2, "returned pointers should not reference the same object")
}

func TestOverWrite(t *testing.T) {
	// Get the file path of the original conf
	pathDirConfOriginal := filepath.Join(util.GetPathDirRepo(), "testdata", "golden", "only_conf")
	pathFileConfOriginal := filepath.Join(pathDirConfOriginal, "config.json")
	require.FileExists(t, pathFileConfOriginal, "conf file for test not found")

	// Create dummy dir for testing
	pathDirTemp := t.TempDir()
	pathDirConf := filepath.Join(pathDirTemp, ".qiitask")

	require.NoError(t, os.Mkdir(pathDirConf, 0o766), "failed to create temp dir during test")

	// Copy conf to temp
	pathFileConf := filepath.Join(pathDirConf, "config.json")

	require.NoError(t, util.CopyFile(pathFileConfOriginal, pathFileConf),
		"failed to copy the original conf file to temp during test")

	conf, err := config.New(pathDirTemp, t.TempDir())
	require.NoError(t, err)
	require.Equal(t, pathFileConf, conf.FileUsed(), "it should use the existing conf")

	// Change conf value before overwrite
	conf.Set("separator_interval", int(1))
	conf.Set("standby_interval", int(2))

	// Overwrite
	require.NoError(t, conf.OverWrite(), "failed to overwrite file")

	// Assert content
	confByte, err := os.ReadFile(pathFileConf)
	require.NoError(t, err, "failed to read witten file")

	assert.Equal(t,
		"{\n  \"queries\": {},\n  \"separator_interval\": 1,\n  \"standby_interval\": 2\n}",
		string(confByte),
	)
}

func TestOverWrite_config_path_not_set(t *testing.T) {
	obj, err := config.New(t.TempDir(), t.TempDir())
	require.NoError(t, err)

	err = obj.OverWrite()

	require.Error(t, err, "overwriting an empty dir should return error")
	assert.Contains(t, err.Error(), "config file path not set. Use SaveAs method instead")
}

func TestSaveAs(t *testing.T) {
	pathDirTemp := t.TempDir()
	pathFileConf := filepath.Join(pathDirTemp, "config.json")

	obj, err := config.New(pathDirTemp, pathDirTemp)
	require.NoError(t, err)

	// Set dummy conf value
	obj.Set("separator_interval", 100)
	obj.Set("standby_interval", 200)

	// Save
	err = obj.SaveAs(pathFileConf)
	require.NoError(t, err)

	// Assert content
	confByte, err := os.ReadFile(pathFileConf)
	require.NoError(t, err, "failed to read witten file")

	assert.Equal(t,
		"{\n  \"queries\": {},\n  \"separator_interval\": 100,\n  \"standby_interval\": 200\n}",
		string(confByte),
	)
}
