package main

import (
	"os"
	"testing"

	"github.com/KEINOS/go-utiles/util"
)

// ----------------------------------------------------------------------------
//  Helper Functions for Main Test
// ----------------------------------------------------------------------------

// MockOsExit mocks the OsExit behavior. Which is a copy of os.Exit().
func MockOsExit(t *testing.T, capturedExitCode *int) (deferFunction func()) {
	t.Helper()

	// Save current function and restore at the end.
	oldOsExit := OsExit
	oldUtilOsExit := util.OsExit

	// Mock OsExit which captures the exit code
	OsExit = func(code int) {
		*capturedExitCode = code
	}

	// Mock UtilOsExit which captures the exit code
	util.OsExit = func(code int) {
		*capturedExitCode = code
	}

	return func() {
		OsExit = oldOsExit
		util.OsExit = oldUtilOsExit
	}
}

// MockOsArgs temporary replaces the os.Args values during test.
func MockOsArgs(t *testing.T, args []string) (deferFunction func()) {
	t.Helper()

	// Backup value
	oldOsArgs := os.Args

	// Mock the args value
	os.Args = append([]string{t.Name()}, args...)

	return func() {
		os.Args = oldOsArgs
	}
}
