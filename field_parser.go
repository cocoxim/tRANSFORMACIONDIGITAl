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
		o