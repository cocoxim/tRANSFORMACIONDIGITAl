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

func TestFormat_DefaultExcludes(t *testing.T)