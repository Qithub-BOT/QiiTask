/*
Package main.
*/
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/KEINOS/go-utiles/util"
	"github.com/Qithub-BOT/QiiTask/core/appinfo"
	"github.com/Qithub-BOT/QiiTask/qiitask/cmd/cmdroot"
)

// FAILURE is an alias of exit status code to ease read.
const FAILURE = 1

// OsExit is a copy of `os.Exit` to ease mocking during test.
//     Ref: https://stackoverflow.com/a/40801733/8367711
//
// We found useful that `os.Exit` should be called only in the main package.
// And functions or methods in other packages should return an error rather than
// exiting the app. Doing so, it's a lot easier to test and trace.
var OsExit = os.Exit

// OsGetWd は os.Getwd のコピーです。テストでモックしやすいように変数に代入して
// 使います。
var OsGetWd = os.Getwd

// OsUserHomeDir は os.UserHomeDir のコピーです。テストでモックしやすいように変
// 数に代入して使います。
var OsUserHomeDir = os.UserHomeDir

// Version info via `-ldflags`.
// These values should be set via `-ldflags` option during build.
//
// Sample command to build:
//     $ VER_APP="$(git describe --tag)" // v1.0.0-alpha-g65a8e0c
//     $ go build -ldflags="-X 'main.version=${VER_APP}'" -o foo ./qiitask/
//     $ ./foo --version
//     qiitask version 1.0.0-alpha (g65a8e0c)
var (
	// The application version to display.
	version string
	// The short hash of the commit. This value will be used as build ver as well.
	commit string
)

// DebugReadBuildInfo is a copy of debug.ReadBuildInfo to ease testing.
var DebugReadBuildInfo = debug.ReadBuildInfo

// ----------------------------------------------------------------------------
//  Main
// ----------------------------------------------------------------------------

func main() {
	// Get directory path for searching files.
	pathDirCurr, err := OsGetWd()
	util.ExitOnErr(err)

	pathDirHome, err := OsUserHomeDir()
	util.ExitOnErr(err)

	// Get app config object
	appInfo, err := appinfo.New(pathDirCurr, pathDirHome, GetVersion())
	util.ExitOnErr(err)

	// Initialize the app then executes the designated command.
	app := cmdroot.New(appInfo)

	if err := app.Execute(); err != nil {
		if appInfo.IsVerbose {
			fmt.Fprintln(os.Stderr, "【詳細情報】")
			util.ExitOnErr(err)
		}

		OsExit(FAILURE)
	}
}

// ----------------------------------------------------------------------------
//  Functions
// ----------------------------------------------------------------------------

// GetVersion returns the app version without the "v" prefix.
func GetVersion() string {
	result := ""

	// Version via build flag
	if version != "" {
		result = version
	}

	// Version via commit tag on `go install`
	if result == "" {
		if buildInfo, ok := DebugReadBuildInfo(); ok {
			result = buildInfo.Main.Version
		}
	}

	return parseVersion(result)
}

func parseVersion(input string) string {
	if input == "" {
		return "(unknown)"
	}

	parsed, err := util.ParseVersion(input)
	if err != nil {
		return input
	}

	major := parsed["major"]
	minor := parsed["minor"]
	patch := parsed["patch"]
	prerelease := parsed["prerelease"]
	build := parsed["build"]

	if commit != "" {
		build = commit
	}

	if prerelease != "" {
		prerelease = "-" + prerelease
	}

	if build != "" {
		build = fmt.Sprintf(" (%v)", build)
	}

	return fmt.Sprintf("%v.%v.%v%v%v", major, minor, patch, prerelease, build)
}
