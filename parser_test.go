package swag

import (
	"bytes"
	"encoding/json"
	"errors"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

const defaultParseDepth = 100

const mainAPIFile = "main.go"

func TestNew(t *testing.T) {
	t.Run("SetMarkdownFileDirectory", func(t *testing.T) {
		t.Parallel()

		expected := "docs/markdown"
		p := New(SetMarkdownFileDirectory(expected))
		assert.Equal(t, expected, p.markdownFileDir)
	})

	t.Run("SetCodeExamplesDirectory", func(t *testing.T) {
		t.Parallel()

		expected := "docs/examples"
		p := New(SetCodeExamplesDirectory(expected))
		assert.Equal(t, expected, p.codeExampleFilesDir)
	})

	t.Run("SetStrict", func(t *testing.T) {
		t.Parallel()

		p := New()
		assert.Equal(t, false, p.Strict)

		p = New(SetStrict(true))
		assert.Equal(t, true, p.Strict)
	})

	t.Run("SetDebugger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

		p := New(SetDebugger(logger))
		assert.Equal(t, logger, p.debug)
	})

	t.Run("SetFieldParserFactory", func(t *testing.T) {
		t.Parallel()

		p := New(SetFieldParserFactory(nil))
		assert.Nil(t, p.fieldParserFactory)
	})
}

func TestSetOverrides(t *testing.T) {
	t.Parallel()

	overrides := map[string]string{
		"foo": "bar",
	}

	p := New(SetOverrides(overrides))
	assert.Equal(t, overrides, p.Overrides)
}

func TestOverrides_getTypeSchema(t *testing.T) {
	t.Parallel()

	overrides := map[string]string{
		"sql.NullString": "string",
	}

	p := New(SetOverrides(overrides))

	t.Run("Override sql.NullString by string", func(t *testing.T) {
		t.Parallel()

		s, err := p.getTypeSchema("sql.NullString", nil, false)
		if assert.NoError(t, err) {
			assert.Truef(t, s.Type.Contains("string"), "type sql.NullString should be overridden by string")
		}
	})

	t.Run("Missing Override for sql.NullInt64", func(t *testing.T) {
		t.Parallel()

		_, err := p.getTypeSchema("sql.NullInt64", nil, false)
		if assert.Error(t, err) {
			assert.Equal(t, "cannot find type definition: sql.NullInt64", err.Error())
		}
	})
}

func TestParser_ParseDefinition(t *testing.T) {
	p := New()

	// Parsing existing type
	definition := &TypeSpecDef{
		PkgPath: "github.com/swagger/swag",
		File: &ast.Fil