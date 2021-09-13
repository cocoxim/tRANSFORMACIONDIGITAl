
package swag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	goparser "go/parser"
	"go/token"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/KyleBanks/depth"
	"github.com/go-openapi/spec"
)

const (
	// CamelCase indicates using CamelCase strategy for struct field.
	CamelCase = "camelcase"

	// PascalCase indicates using PascalCase strategy for struct field.
	PascalCase = "pascalcase"

	// SnakeCase indicates using SnakeCase strategy for struct field.
	SnakeCase = "snakecase"

	idAttr                  = "@id"
	acceptAttr              = "@accept"
	produceAttr             = "@produce"
	paramAttr               = "@param"
	successAttr             = "@success"
	failureAttr             = "@failure"
	responseAttr            = "@response"
	headerAttr              = "@header"
	tagsAttr                = "@tags"
	routerAttr              = "@router"
	summaryAttr             = "@summary"
	deprecatedAttr          = "@deprecated"
	securityAttr            = "@security"
	titleAttr               = "@title"
	conNameAttr             = "@contact.name"
	conURLAttr              = "@contact.url"
	conEmailAttr            = "@contact.email"
	licNameAttr             = "@license.name"
	licURLAttr              = "@license.url"
	versionAttr             = "@version"
	descriptionAttr         = "@description"
	descriptionMarkdownAttr = "@description.markdown"
	secBasicAttr            = "@securitydefinitions.basic"
	secAPIKeyAttr           = "@securitydefinitions.apikey"
	secApplicationAttr      = "@securitydefinitions.oauth2.application"
	secImplicitAttr         = "@securitydefinitions.oauth2.implicit"
	secPasswordAttr         = "@securitydefinitions.oauth2.password"
	secAccessCodeAttr       = "@securitydefinitions.oauth2.accesscode"
	tosAttr                 = "@termsofservice"
	extDocsDescAttr         = "@externaldocs.description"
	extDocsURLAttr          = "@externaldocs.url"
	xCodeSamplesAttr        = "@x-codesamples"
	scopeAttrPrefix         = "@scope."
)

// ParseFlag determine what to parse
type ParseFlag int

const (
	// ParseNone parse nothing
	ParseNone ParseFlag = 0x00
	// ParseOperations parse operations
	ParseOperations = 0x01
	// ParseModels parse models
	ParseModels = 0x02
	// ParseAll parse operations and models
	ParseAll = ParseOperations | ParseModels
)

var (
	// ErrRecursiveParseStruct recursively parsing struct.
	ErrRecursiveParseStruct = errors.New("recursively parsing struct")

	// ErrFuncTypeField field type is func.
	ErrFuncTypeField = errors.New("field type is func")

	// ErrFailedConvertPrimitiveType Failed to convert for swag to interpretable type.
	ErrFailedConvertPrimitiveType = errors.New("swag property: failed convert primitive type")

	// ErrSkippedField .swaggo specifies field should be skipped.
	ErrSkippedField = errors.New("field is skipped by global overrides")
)

var allMethod = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodPut:     {},
	http.MethodPost:    {},
	http.MethodDelete:  {},
	http.MethodOptions: {},
	http.MethodHead:    {},
	http.MethodPatch:   {},
}

// Parser implements a parser for Go source files.
type Parser struct {
	// swagger represents the root document object for the API specification
	swagger *spec.Swagger

	// packages store entities of APIs, definitions, file, package path etc.  and their relations
	packages *PackagesDefinitions

	// parsedSchemas store schemas which have been parsed from ast.TypeSpec
	parsedSchemas map[*TypeSpecDef]*Schema

	// outputSchemas store schemas which will be export to swagger
	outputSchemas map[*TypeSpecDef]*Schema

	// PropNamingStrategy naming strategy
	PropNamingStrategy string

	// ParseVendor parse vendor folder
	ParseVendor bool

	// ParseDependencies whether swag should be parse outside dependency folder
	ParseDependency bool

	// ParseInternal whether swag should parse internal packages
	ParseInternal bool

	// Strict whether swag should error or warn when it detects cases which are most likely user errors
	Strict bool

	// RequiredByDefault set validation required for all fields by default
	RequiredByDefault bool

	// structStack stores full names of the structures that were already parsed or are being parsed now
	structStack []*TypeSpecDef

	// markdownFileDir holds the path to the folder, where markdown files are stored
	markdownFileDir string

	// codeExampleFilesDir holds path to the folder, where code example files are stored
	codeExampleFilesDir string

	// collectionFormatInQuery set the default collectionFormat otherwise then 'csv' for array in query params
	collectionFormatInQuery string

	// excludes excludes dirs and files in SearchDir
	excludes map[string]struct{}

	// tells parser to include only specific extension
	parseExtension string

	// debugging output goes here
	debug Debugger

	// fieldParserFactory create FieldParser
	fieldParserFactory FieldParserFactory

	// Overrides allows global replacements of types. A blank replacement will be skipped.
	Overrides map[string]string

	// parseGoList whether swag use go list to parse dependency
	parseGoList bool

	// tags to filter the APIs after
	tags map[string]struct{}
}

// FieldParserFactory create FieldParser.
type FieldParserFactory func(ps *Parser, field *ast.Field) FieldParser

// FieldParser parse struct field.
type FieldParser interface {
	ShouldSkip() bool
	FieldName() (string, error)
	CustomSchema() (*spec.Schema, error)
	ComplementSchema(schema *spec.Schema) error
	IsRequired() (bool, error)
}

// Debugger is the interface that wraps the basic Printf method.
type Debugger interface {
	Printf(format string, v ...interface{})
}

// New creates a new Parser with default properties.
func New(options ...func(*Parser)) *Parser {
	parser := &Parser{
		swagger: &spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Contact: &spec.ContactInfo{},
						License: nil,
					},
					VendorExtensible: spec.VendorExtensible{
						Extensions: spec.Extensions{},
					},
				},