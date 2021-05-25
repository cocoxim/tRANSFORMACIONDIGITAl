
package swag

import (
	"go/ast"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestDefaultFieldParser(t *testing.T) {
	t.Run("Example tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"string"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" example:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, "one", schema.Example)

		schema = spec.Schema{}
		schema.Type = []string{"string"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" example:""`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, "", schema.Example)

		schema = spec.Schema{}
		schema.Type = []string{"float"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" example:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)
	})

	t.Run("Format tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"string"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" format:"csv"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, "csv", schema.Format)
	})

	t.Run("Required tag", func(t *testing.T) {
		t.Parallel()

		got, err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" binding:"required"`,
			}},
		).IsRequired()
		assert.NoError(t, err)
		assert.Equal(t, true, got)

		got, err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" validate:"required"`,
			}},
		).IsRequired()
		assert.NoError(t, err)
		assert.Equal(t, true, got)
	})

	t.Run("Default required tag", func(t *testing.T) {
		t.Parallel()

		got, err := newTagBaseFieldParser(
			&Parser{
				RequiredByDefault: true,
			},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test"`,
			}},
		).IsRequired()
		assert.NoError(t, err)
		assert.True(t, got)
	})

	t.Run("Optional tag", func(t *testing.T) {
		t.Parallel()

		got, err := newTagBaseFieldParser(
			&Parser{
				RequiredByDefault: true,
			},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" binding:"optional"`,
			}},
		).IsRequired()
		assert.NoError(t, err)
		assert.False(t, got)

		got, err = newTagBaseFieldParser(
			&Parser{
				RequiredByDefault: true,
			},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" validate:"optional"`,
			}},
		).IsRequired()
		assert.NoError(t, err)
		assert.False(t, got)
	})

	t.Run("Extensions tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"int"}
		schema.Extensions = map[string]interface{}{}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" extensions:"x-nullable,x-abc=def,!x-omitempty,x-example=[0, 9],x-example2={çãíœ, (bar=(abc, def)), [0,9]}"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, true, schema.Extensions["x-nullable"])
		assert.Equal(t, "def", schema.Extensions["x-abc"])
		assert.Equal(t, false, schema.Extensions["x-omitempty"])
		assert.Equal(t, "[0, 9]", schema.Extensions["x-example"])
		assert.Equal(t, "{çãíœ, (bar=(abc, def)), [0,9]}", schema.Extensions["x-example2"])
	})

	t.Run("Enums tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"string"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" enums:"a,b,c"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{"a", "b", "c"}, schema.Enum)

		schema = spec.Schema{}
		schema.Type = []string{"float"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" enums:"a,b,c"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)
	})

	t.Run("EnumVarNames tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"int"}
		schema.Extensions = map[string]interface{}{}
		schema.Enum = []interface{}{}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" enums:"0,1,2" x-enum-varnames:"Daily,Weekly,Monthly"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{"Daily", "Weekly", "Monthly"}, schema.Extensions["x-enum-varnames"])

		schema = spec.Schema{}
		schema.Type = []string{"int"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" enums:"0,1,2,3" x-enum-varnames:"Daily,Weekly,Monthly"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)

		// Test for an array of enums
		schema = spec.Schema{}
		schema.Type = []string{"array"}
		schema.Items = &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"int"},
				},
			},
		}
		schema.Extensions = map[string]interface{}{}
		schema.Enum = []interface{}{}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" enums:"0,1,2" x-enum-varnames:"Daily,Weekly,Monthly"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, []interface{}{"Daily", "Weekly", "Monthly"}, schema.Items.Schema.Extensions["x-enum-varnames"])
		assert.Equal(t, spec.Extensions{}, schema.Extensions)
	})

	t.Run("Default tag", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"string"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" default:"pass"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		assert.Equal(t, "pass", schema.Default)

		schema = spec.Schema{}
		schema.Type = []string{"float"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" default:"pass"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)
	})

	t.Run("Numeric value", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"integer"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maximum:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		max := float64(1)
		assert.Equal(t, &max, schema.Maximum)

		schema = spec.Schema{}
		schema.Type = []string{"integer"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maximum:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)

		schema = spec.Schema{}
		schema.Type = []string{"number"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maximum:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		max = float64(1)
		assert.Equal(t, &max, schema.Maximum)

		schema = spec.Schema{}
		schema.Type = []string{"number"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maximum:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)

		schema = spec.Schema{}
		schema.Type = []string{"number"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" multipleOf:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		multipleOf := float64(1)
		assert.Equal(t, &multipleOf, schema.MultipleOf)

		schema = spec.Schema{}
		schema.Type = []string{"number"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" multipleOf:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)

		schema = spec.Schema{}
		schema.Type = []string{"integer"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" minimum:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		min := float64(1)
		assert.Equal(t, &min, schema.Minimum)

		schema = spec.Schema{}
		schema.Type = []string{"integer"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" minimum:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)
	})

	t.Run("String value", func(t *testing.T) {
		t.Parallel()

		schema := spec.Schema{}
		schema.Type = []string{"string"}
		err := newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maxLength:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		max := int64(1)
		assert.Equal(t, &max, schema.MaxLength)

		schema = spec.Schema{}
		schema.Type = []string{"string"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" maxLength:"one"`,
			}},
		).ComplementSchema(&schema)
		assert.Error(t, err)

		schema = spec.Schema{}
		schema.Type = []string{"string"}
		err = newTagBaseFieldParser(
			&Parser{},
			&ast.Field{Tag: &ast.BasicLit{
				Value: `json:"test" minLength:"1"`,
			}},
		).ComplementSchema(&schema)
		assert.NoError(t, err)
		min := int64(1)
		assert.Equal(t, &min, schema.MinLength)
