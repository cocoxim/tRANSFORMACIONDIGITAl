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
	assert.NoError(t, New().Build(&Config{Search