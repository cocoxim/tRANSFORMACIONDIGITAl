package swag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	SearchDir = "./testdata/format_test"
	Excludes  = "./testdata/format_test/web"
	MainFile  = "main.go"
)

func testFormat(t *testing.T, filename, contents, want string) {
	got, err := NewFormatter().Format(filename, []byte(contents))
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func Test_FormatMain(t *testing.T) {
	contents := `package main
	// @title Swagger Example API
	// @version 1.0
	// @