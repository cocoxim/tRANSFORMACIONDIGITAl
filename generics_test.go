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
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsProperty(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_property"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	os.WriteFile(searchDir+"/expected.json", b, fs.ModePerm)
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsNames(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_names"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParseGenericsPackageAlias(t *testing.T) {
	t.Parallel()

	searchDir := "testdata/generics_package_alias/internal"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New(SetParseDependency(true))
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
}

func TestParametrizeStruct(t *testing.T) {
	pd := PackagesDefinitions{
		packages:          make(map[string]*PackageDefinitions),
		uniqueDefinitions: make(map[string]*TypeSpecDef),
	}
	// valid
	typeSpec := pd.parametrizeGenericType(
		&ast.File{Name: &ast.Ident{Name: "test2"}},
		&TypeSpecDef{
			File: &ast.File{Name: &ast.Ident{Name: "test"}},
			TypeSpec: &ast.TypeSpec{
				Name:       &ast.Ident{Name: "Field"},
				TypeParams: &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "T"}}}, {Names: []*ast.Ident{{Name: "T2"}}}}},
				Type:       &ast.StructType{Struct: 100, Fields: &ast.FieldList{Opening: 101, Closing: 102}},
			}}, "test.Field[string, []string]")
	assert.NotNil(t, typeSpec)
	assert.Equal(t, "$test.Field-string-array_string", typeSpec.Name())
	assert.Equal(t, "test.Field-string-array_string", typeSpec.TypeName())

	// definition contains one type params, but two type params are provided
	typeSpec = pd.parametrizeGenericType(
		&ast.File{Name: &ast.Ident{Name: "test2"}},
		&TypeSpecDef{
			TypeSpec: &ast.TypeSpec{
				Name:       &ast.Ident{Name: "Field"},
				TypeParams: &ast.FieldList{List: []*ast.Field{{Names: []*ast.Ident{{Name: "T"}}}}},
				Type:       &ast.StructType{Struct: 100, Fields: &ast.FieldList{Opening: 101, Closing: 102}},
			}}, "test.Field[string, string]")
	assert.Nil(t, typeSpec)

	// defi