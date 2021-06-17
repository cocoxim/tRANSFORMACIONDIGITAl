//go:build go1.18
// +build go1.18

package swag

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLogger struct {
	Messages []string
}

func (t *testLogger) Printf(format string, v ...interface{}) {
	t.Messages = append(t.Messages, fmt.Sprintf(format, v...))
}

func TestParseGenericsBasic(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_basic"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	p.Overrides = map[string]string{
		"types.Field[string]":               "string",
		"types.DoubleField[string,string]":  "[]string",
		"types.TrippleField[string,string]": "[][]string",
	}

	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsArrays(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_arrays"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsNested(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_nested"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsMultiLevelNesting(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_multi_level_nesting"
	expected, err := os.ReadFi