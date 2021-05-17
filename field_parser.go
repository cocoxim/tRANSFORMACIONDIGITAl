package swag

import (
	"fmt"
	"go/ast"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"github.com/go-openapi/spec"
)

var _ FieldParser = &tagBaseFieldParser{p: nil, field: nil, tag: ""}

const (
	requiredLabel    = "required"
	optionalLabel    = "optional"
	swaggerTypeTag   = "swaggertype"
	swaggerIgnoreTag = "swaggerignore"
)

type tagBaseFieldParser struct {
	p     *Parser
	field *ast.Field
	tag   reflect.StructTag
}

func newTagBaseFieldParser(p *Parser, field *ast.Field) FieldParser {
	fieldParser := tagBaseFieldParser{
		p:     p,
		field: field,
		tag:   "",
	}
	if fieldParser.field.Tag != nil {
		fieldParser.tag = reflect.StructTag(strings.ReplaceAll(field.Tag.Value, "`", ""))
	}

	return &fieldParser
}

func (ps *tagBaseFieldParser) ShouldSkip() bool {
	// Skip non-exported fields.
	if ps.field.Names != nil && !ast.IsExported(ps.field.Names[0].Name) {
		return true
	}

	if ps.field.Tag == nil {
		return false
	}

	ignoreTag := ps.tag.Get(swaggerIgnoreTag)
	if strings.EqualFold(ignoreTag, "true") {
		return true
	}

	// json:"tag,hoge"
	name := strings.TrimSpace(strings.Split(ps.tag.Get(jsonTag), ",")[0])
	if name == "-" {
		return true
	}

	return false
}

func (ps *tagBaseFieldParser) FieldName() (string, error) {
	var name string
	if ps.field.Tag != nil {
		// json:"tag,hoge"
		name = strings.TrimSpace(strings.Split(ps.tag.Get(jsonTag), ",")[0])

		if name != "" {
			return name, nil
		}
	}

	if ps.field.Names == nil {
		return "", nil
	}

	switch ps.p.PropNamingStrategy {
	case SnakeCase:
		return toSnakeCase(ps.field.Names[0].Name), nil
	case PascalCase:
		return ps.field.Names[0].Name, nil
	default:
		return toLowerCamelCase(ps.field.Names[0].Name), nil
	}
}

func toSnakeCase(in string) string {
	var (
		runes  = []rune(in)
		length = len(runes)
		out    []rune
	)

	for idx := 0; idx < length; idx++ {
		if idx > 0 && unicode.IsUpper(runes[idx]) &&
			((idx+1 < length && unicode.IsLower(runes[idx+1])) || unicode.IsLower(runes[idx-1])) {
			out = append(out, '_')
		}

		out = append(out, unicode.ToLower(runes[idx]))
	}

	return string(out)
}

func toLowerCamelCase(in string) string {
	var flag bool

	out := make([]rune, len(in))

	runes := []rune(in)
	for i, curr := range runes {
		if (i == 0 && unicode.IsUpper(curr)) || (flag && unicode.IsUpper(curr)) {
			out[i] = unicode.ToLower(curr)
			flag = true

			continue
		}

		out[i] = curr
		flag = false
	}

	return string(out)
}

func (ps *tagBaseFieldParser) CustomSchema() (*spec.Schema, error) {
	if ps.field.Tag == nil {
		return nil, nil
	}

	typeTag := ps.tag.Get(swaggerTypeTag)
	if typeTag != "" {
		return BuildCustomSchema(strings.Split(typeTag, ","))
	}

	return nil, nil
}

type structField struct {
	schemaType   string
	arrayType    string
	formatType   string
	maximum      *float64
	minimum      *float64
	multipleOf   *float64
	maxLength    *int64
	minLength    *int64
	maxItems     *int64
	minItems     *int64
	exampleValue interface{}
	enums        []interface{}
	enumVarNames []interface{}
	unique       bool
}

// splitNotWrapped slices s into all substrings separated by sep if sep is not
// wrapped by brackets and returns a slice of the substrings between those separators.
func splitNotWrapped(s string, sep rune) []string {
	openCloseMap := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}

	var (
		result    = make([]string, 0)
		current   = strings.Builder{}
		openCount = 0
		openChar  rune
	)

	for _, char := range s {
		switch {
		case openChar == 0 && openCloseMap[char] != 0:
			openChar = char

			openCount++

			current.WriteRune(char)
		case char == openChar:
			openCo