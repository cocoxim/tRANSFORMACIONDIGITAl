package format

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat_Format(t *testing.T) {
	fx := setup(t)
	assert.NoError(t, New().Build(&Config{SearchDir: fx.basedir}))
	assert.True(t, fx.isFormatted("main.go"))
	assert.True(t, fx.isFormatted("api/api.go"))
}

func TestFormat_ExcludeDir(t *testing.T) {
	fx := setup(t)
	assert.NoError(t, New().Build(&Config{
		SearchDir: fx.basedir,
		Excludes:  filepath.Join(fx.basedir, "api"),
	}))
	assert.False(t, fx.isFormatted("api/api.go"))
}

func TestFormat_ExcludeFile(t *testing.T) {
	fx := setup(t)
	assert.NoError(t, New().Build(&Config{
		SearchDir: fx.basedir,
		Excludes:  filepath.Join(fx.basedir, "main.go"),
	}))
	assert.False(t, fx.isFormatted("main.go"))
}

func TestFormat_DefaultExcludes(t *testing.T) {
	fx := setup(t)
	assert.NoError(t, New().Build(&Config{SearchDir: fx.basedir}))
	assert.False(t, fx.isFormatted("api/api_test.go"))
	assert.False(t, fx.isFormatted("docs/docs.go"))
}

func TestFormat_ParseError(t *testing.T) {
	fx := setup(t)
	os.WriteFile(filepath.Join(fx.basedir, "parse_error.go"), []byte(`package main
		func invalid() {`), 0644)
	assert.Error(t, New().Build(&Config{SearchDir: fx.basedir}))
}

func TestFormat_ReadError(t *testing.T) {
	fx := setup(t)
	os.Chmod(filepath.Join(fx.basedir, "main.go"), 0)
	assert.Error(t, New().Build(&Config{SearchDir: fx.basedir}))
}

func TestFormat_WriteError(t *testing.T) {
	fx := set