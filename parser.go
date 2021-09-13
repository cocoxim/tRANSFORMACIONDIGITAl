
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
				Paths: &spec.Paths{
					Paths: make(map[string]spec.PathItem),
					VendorExtensible: spec.VendorExtensible{
						Extensions: nil,
					},
				},
				Definitions:         make(map[string]spec.Schema),
				SecurityDefinitions: make(map[string]*spec.SecurityScheme),
			},
			VendorExtensible: spec.VendorExtensible{
				Extensions: nil,
			},
		},
		packages:           NewPackagesDefinitions(),
		debug:              log.New(os.Stdout, "", log.LstdFlags),
		parsedSchemas:      make(map[*TypeSpecDef]*Schema),
		outputSchemas:      make(map[*TypeSpecDef]*Schema),
		excludes:           make(map[string]struct{}),
		tags:               make(map[string]struct{}),
		fieldParserFactory: newTagBaseFieldParser,
		Overrides:          make(map[string]string),
	}

	for _, option := range options {
		option(parser)
	}

	parser.packages.debug = parser.debug

	return parser
}

// SetParseDependency sets whether to parse the dependent packages.
func SetParseDependency(parseDependency bool) func(*Parser) {
	return func(p *Parser) {
		p.ParseDependency = parseDependency
		if p.packages != nil {
			p.packages.parseDependency = parseDependency
		}
	}
}

// SetMarkdownFileDirectory sets the directory to search for markdown files.
func SetMarkdownFileDirectory(directoryPath string) func(*Parser) {
	return func(p *Parser) {
		p.markdownFileDir = directoryPath
	}
}

// SetCodeExamplesDirectory sets the directory to search for code example files.
func SetCodeExamplesDirectory(directoryPath string) func(*Parser) {
	return func(p *Parser) {
		p.codeExampleFilesDir = directoryPath
	}
}

// SetExcludedDirsAndFiles sets directories and files to be excluded when searching.
func SetExcludedDirsAndFiles(excludes string) func(*Parser) {
	return func(p *Parser) {
		for _, f := range strings.Split(excludes, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				f = filepath.Clean(f)
				p.excludes[f] = struct{}{}
			}
		}
	}
}

// SetTags sets the tags to be included
func SetTags(include string) func(*Parser) {
	return func(p *Parser) {
		for _, f := range strings.Split(include, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				p.tags[f] = struct{}{}
			}
		}
	}
}

// SetParseExtension parses only those operations which match given extension
func SetParseExtension(parseExtension string) func(*Parser) {
	return func(p *Parser) {
		p.parseExtension = parseExtension
	}
}

// SetStrict sets whether swag should error or warn when it detects cases which are most likely user errors.
func SetStrict(strict bool) func(*Parser) {
	return func(p *Parser) {
		p.Strict = strict
	}
}

// SetDebugger allows the use of user-defined implementations.
func SetDebugger(logger Debugger) func(parser *Parser) {
	return func(p *Parser) {
		if logger != nil {
			p.debug = logger
		}
	}
}

// SetFieldParserFactory allows the use of user-defined implementations.
func SetFieldParserFactory(factory FieldParserFactory) func(parser *Parser) {
	return func(p *Parser) {
		p.fieldParserFactory = factory
	}
}

// SetOverrides allows the use of user-defined global type overrides.
func SetOverrides(overrides map[string]string) func(parser *Parser) {
	return func(p *Parser) {
		for k, v := range overrides {
			p.Overrides[k] = v
		}
	}
}

// ParseUsingGoList sets whether swag use go list to parse dependency
func ParseUsingGoList(enabled bool) func(parser *Parser) {
	return func(p *Parser) {
		p.parseGoList = enabled
	}
}

// ParseAPI parses general api info for given searchDir and mainAPIFile.
func (parser *Parser) ParseAPI(searchDir string, mainAPIFile string, parseDepth int) error {
	return parser.ParseAPIMultiSearchDir([]string{searchDir}, mainAPIFile, parseDepth)
}

// ParseAPIMultiSearchDir is like ParseAPI but for multiple search dirs.
func (parser *Parser) ParseAPIMultiSearchDir(searchDirs []string, mainAPIFile string, parseDepth int) error {
	for _, searchDir := range searchDirs {
		parser.debug.Printf("Generate general API Info, search dir:%s", searchDir)

		packageDir, err := getPkgName(searchDir)
		if err != nil {
			parser.debug.Printf("warning: failed to get package name in dir: %s, error: %s", searchDir, err.Error())
		}

		err = parser.getAllGoFileInfo(packageDir, searchDir)
		if err != nil {
			return err
		}
	}

	absMainAPIFilePath, err := filepath.Abs(filepath.Join(searchDirs[0], mainAPIFile))
	if err != nil {
		return err
	}

	// Use 'go list' command instead of depth.Resolve()
	if parser.ParseDependency {
		if parser.parseGoList {
			pkgs, err := listPackages(context.Background(), filepath.Dir(absMainAPIFilePath), nil, "-deps")
			if err != nil {
				return fmt.Errorf("pkg %s cannot find all dependencies, %s", filepath.Dir(absMainAPIFilePath), err)
			}

			length := len(pkgs)
			for i := 0; i < length; i++ {
				err := parser.getAllGoFileInfoFromDepsByList(pkgs[i])
				if err != nil {
					return err
				}
			}
		} else {
			var t depth.Tree
			t.ResolveInternal = true
			t.MaxDepth = parseDepth

			pkgName, err := getPkgName(filepath.Dir(absMainAPIFilePath))
			if err != nil {
				return err
			}

			err = t.Resolve(pkgName)
			if err != nil {
				return fmt.Errorf("pkg %s cannot find all dependencies, %s", pkgName, err)
			}
			for i := 0; i < len(t.Root.Deps); i++ {
				err := parser.getAllGoFileInfoFromDeps(&t.Root.Deps[i])
				if err != nil {
					return err
				}
			}
		}
	}

	err = parser.ParseGeneralAPIInfo(absMainAPIFilePath)
	if err != nil {
		return err
	}

	parser.parsedSchemas, err = parser.packages.ParseTypes()
	if err != nil {
		return err
	}

	err = parser.packages.RangeFiles(parser.ParseRouterAPIInfo)
	if err != nil {
		return err
	}

	return parser.checkOperationIDUniqueness()
}

func getPkgName(searchDir string) (string, error) {
	cmd := exec.Command("go", "list", "-f={{.ImportPath}}")
	cmd.Dir = searchDir

	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list command, %s, stdout:%s, stderr:%s", err, stdout.String(), stderr.String())
	}

	outStr, _ := stdout.String(), stderr.String()

	if outStr[0] == '_' { // will shown like _/{GOPATH}/src/{YOUR_PACKAGE} when NOT enable GO MODULE.
		outStr = strings.TrimPrefix(outStr, "_"+build.Default.GOPATH+"/src/")
	}

	f := strings.Split(outStr, "\n")

	outStr = f[0]

	return outStr, nil
}

// ParseGeneralAPIInfo parses general api info for given mainAPIFile path.
func (parser *Parser) ParseGeneralAPIInfo(mainAPIFile string) error {
	fileTree, err := goparser.ParseFile(token.NewFileSet(), mainAPIFile, nil, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("cannot parse source files %s: %s", mainAPIFile, err)
	}

	parser.swagger.Swagger = "2.0"

	for _, comment := range fileTree.Comments {
		comments := strings.Split(comment.Text(), "\n")
		if !isGeneralAPIComment(comments) {
			continue
		}

		err = parseGeneralAPIInfo(parser, comments)
		if err != nil {
			return err
		}
	}

	return nil
}

func parseGeneralAPIInfo(parser *Parser, comments []string) error {
	previousAttribute := ""

	// parsing classic meta data model
	for line := 0; line < len(comments); line++ {
		commentLine := comments[line]
		commentLine = strings.TrimSpace(commentLine)
		if len(commentLine) == 0 {
			continue
		}
		fields := FieldsByAnySpace(commentLine, 2)

		attribute := fields[0]
		var value string
		if len(fields) > 1 {
			value = fields[1]
		}

		switch attr := strings.ToLower(attribute); attr {
		case versionAttr, titleAttr, tosAttr, licNameAttr, licURLAttr, conNameAttr, conURLAttr, conEmailAttr:
			setSwaggerInfo(parser.swagger, attr, value)
		case descriptionAttr:
			if previousAttribute == attribute {
				parser.swagger.Info.Description += "\n" + value

				continue
			}

			setSwaggerInfo(parser.swagger, attr, value)
		case descriptionMarkdownAttr:
			commentInfo, err := getMarkdownForTag("api", parser.markdownFileDir)
			if err != nil {
				return err
			}

			setSwaggerInfo(parser.swagger, descriptionAttr, string(commentInfo))

		case "@host":
			parser.swagger.Host = value
		case "@basepath":
			parser.swagger.BasePath = value

		case acceptAttr:
			err := parser.ParseAcceptComment(value)
			if err != nil {
				return err
			}
		case produceAttr:
			err := parser.ParseProduceComment(value)
			if err != nil {
				return err
			}
		case "@schemes":
			parser.swagger.Schemes = strings.Split(value, " ")
		case "@tag.name":
			parser.swagger.Tags = append(parser.swagger.Tags, spec.Tag{
				TagProps: spec.TagProps{
					Name: value,
				},
			})
		case "@tag.description":
			tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
			tag.TagProps.Description = value
			replaceLastTag(parser.swagger.Tags, tag)
		case "@tag.description.markdown":
			tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]

			commentInfo, err := getMarkdownForTag(tag.TagProps.Name, parser.markdownFileDir)
			if err != nil {
				return err
			}

			tag.TagProps.Description = string(commentInfo)
			replaceLastTag(parser.swagger.Tags, tag)
		case "@tag.docs.url":
			tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
			tag.TagProps.ExternalDocs = &spec.ExternalDocumentation{
				URL:         value,
				Description: "",
			}

			replaceLastTag(parser.swagger.Tags, tag)
		case "@tag.docs.description":
			tag := parser.swagger.Tags[len(parser.swagger.Tags)-1]
			if tag.TagProps.ExternalDocs == nil {
				return fmt.Errorf("%s needs to come after a @tags.docs.url", attribute)
			}

			tag.TagProps.ExternalDocs.Description = value
			replaceLastTag(parser.swagger.Tags, tag)

		case secBasicAttr, secAPIKeyAttr, secApplicationAttr, secImplicitAttr, secPasswordAttr, secAccessCodeAttr:
			scheme, err := parseSecAttributes(attribute, comments, &line)
			if err != nil {
				return err
			}

			parser.swagger.SecurityDefinitions[value] = scheme

		case "@query.collection.format":
			parser.collectionFormatInQuery = TransToValidCollectionFormat(value)

		case extDocsDescAttr, extDocsURLAttr:
			if parser.swagger.ExternalDocs == nil {
				parser.swagger.ExternalDocs = new(spec.ExternalDocumentation)
			}
			switch attr {
			case extDocsDescAttr:
				parser.swagger.ExternalDocs.Description = value
			case extDocsURLAttr:
				parser.swagger.ExternalDocs.URL = value
			}

		default:
			if strings.HasPrefix(attribute, "@x-") {
				extensionName := attribute[1:]

				extExistsInSecurityDef := false
				// for each security definition
				for _, v := range parser.swagger.SecurityDefinitions {
					// check if extension exists
					_, extExistsInSecurityDef = v.VendorExtensible.Extensions.GetString(extensionName)
					// if it exists in at least one, then we stop iterating
					if extExistsInSecurityDef {
						break
					}
				}

				// if it is present on security def, don't add it again
				if extExistsInSecurityDef {
					break
				}

				if len(value) == 0 {
					return fmt.Errorf("annotation %s need a value", attribute)
				}

				var valueJSON interface{}
				err := json.Unmarshal([]byte(value), &valueJSON)
				if err != nil {
					return fmt.Errorf("annotation %s need a valid json value", attribute)
				}

				if strings.Contains(extensionName, "logo") {
					parser.swagger.Info.Extensions.Add(extensionName, valueJSON)
				} else {
					if parser.swagger.Extensions == nil {
						parser.swagger.Extensions = make(map[string]interface{})
					}

					parser.swagger.Extensions[attribute[1:]] = valueJSON
				}
			}
		}

		previousAttribute = attribute
	}

	return nil
}

func setSwaggerInfo(swagger *spec.Swagger, attribute, value string) {
	switch attribute {
	case versionAttr:
		swagger.Info.Version = value
	case titleAttr:
		swagger.Info.Title = value
	case tosAttr:
		swagger.Info.TermsOfService = value
	case descriptionAttr:
		swagger.Info.Description = value
	case conNameAttr:
		swagger.Info.Contact.Name = value
	case conEmailAttr:
		swagger.Info.Contact.Email = value
	case conURLAttr:
		swagger.Info.Contact.URL = value
	case licNameAttr:
		swagger.Info.License = initIfEmpty(swagger.Info.License)
		swagger.Info.License.Name = value
	case licURLAttr:
		swagger.Info.License = initIfEmpty(swagger.Info.License)
		swagger.Info.License.URL = value
	}
}

func parseSecAttributes(context string, lines []string, index *int) (*spec.SecurityScheme, error) {
	const (
		in               = "@in"
		name             = "@name"
		descriptionAttr  = "@description"
		tokenURL         = "@tokenurl"
		authorizationURL = "@authorizationurl"
	)

	var search []string

	attribute := strings.ToLower(FieldsByAnySpace(lines[*index], 2)[0])
	switch attribute {
	case secBasicAttr:
		return spec.BasicAuth(), nil
	case secAPIKeyAttr:
		search = []string{in, name}
	case secApplicationAttr, secPasswordAttr:
		search = []string{tokenURL}
	case secImplicitAttr:
		search = []string{authorizationURL}
	case secAccessCodeAttr:
		search = []string{tokenURL, authorizationURL}
	}

	// For the first line we get the attributes in the context parameter, so we skip to the next one
	*index++

	attrMap, scopes := make(map[string]string), make(map[string]string)
	extensions, description := make(map[string]interface{}), ""

	for ; *index < len(lines); *index++ {
		v := strings.TrimSpace(lines[*index])
		if len(v) == 0 {
			continue
		}

		fields := FieldsByAnySpace(v, 2)
		securityAttr := strings.ToLower(fields[0])
		var value string
		if len(fields) > 1 {
			value = fields[1]
		}

		for _, findterm := range search {
			if securityAttr == findterm {
				attrMap[securityAttr] = value

				break
			}
		}

		isExists, err := isExistsScope(securityAttr)
		if err != nil {
			return nil, err
		}

		if isExists {
			scopes[securityAttr[len(scopeAttrPrefix):]] = v[len(securityAttr):]
		}

		if strings.HasPrefix(securityAttr, "@x-") {
			// Add the custom attribute without the @
			extensions[securityAttr[1:]] = value
		}

		// Not mandatory field
		if securityAttr == descriptionAttr {
			description = value
		}

		// next securityDefinitions
		if strings.Index(securityAttr, "@securitydefinitions.") == 0 {
			// Go back to the previous line and break
			*index--

			break
		}
	}

	if len(attrMap) != len(search) {
		return nil, fmt.Errorf("%s is %v required", context, search)
	}

	var scheme *spec.SecurityScheme

	switch attribute {
	case secAPIKeyAttr:
		scheme = spec.APIKeyAuth(attrMap[name], attrMap[in])
	case secApplicationAttr:
		scheme = spec.OAuth2Application(attrMap[tokenURL])
	case secImplicitAttr:
		scheme = spec.OAuth2Implicit(attrMap[authorizationURL])
	case secPasswordAttr:
		scheme = spec.OAuth2Password(attrMap[tokenURL])
	case secAccessCodeAttr:
		scheme = spec.OAuth2AccessToken(attrMap[authorizationURL], attrMap[tokenURL])
	}

	scheme.Description = description

	for extKey, extValue := range extensions {
		scheme.AddExtension(extKey, extValue)
	}

	for scope, scopeDescription := range scopes {
		scheme.AddScope(scope, scopeDescription)
	}

	return scheme, nil
}

func initIfEmpty(license *spec.License) *spec.License {
	if license == nil {
		return new(spec.License)
	}

	return license
}

// ParseAcceptComment parses comment for given `accept` comment string.
func (parser *Parser) ParseAcceptComment(commentLine string) error {
	return parseMimeTypeList(commentLine, &parser.swagger.Consumes, "%v accept type can't be accepted")
}

// ParseProduceComment parses comment for given `produce` comment string.
func (parser *Parser) ParseProduceComment(commentLine string) error {
	return parseMimeTypeList(commentLine, &parser.swagger.Produces, "%v produce type can't be accepted")
}

func isGeneralAPIComment(comments []string) bool {
	for _, commentLine := range comments {
		commentLine = strings.TrimSpace(commentLine)
		if len(commentLine) == 0 {
			continue
		}
		attribute := strings.ToLower(FieldsByAnySpace(commentLine, 2)[0])
		switch attribute {
		// The @summary, @router, @success, @failure annotation belongs to Operation
		case summaryAttr, routerAttr, successAttr, failureAttr, responseAttr:
			return false
		}
	}

	return true
}

func getMarkdownForTag(tagName string, dirPath string) ([]byte, error) {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

		if !strings.Contains(fileName, ".md") {
			continue
		}

		if strings.Contains(fileName, tagName) {
			fullPath := filepath.Join(dirPath, fileName)

			commentInfo, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, fmt.Errorf("Failed to read markdown file %s error: %s ", fullPath, err)
			}

			return commentInfo, nil
		}
	}

	return nil, fmt.Errorf("Unable to find markdown file for tag %s in the given directory", tagName)
}

func isExistsScope(scope string) (bool, error) {
	s := strings.Fields(scope)
	for _, v := range s {
		if strings.Contains(v, scopeAttrPrefix) {
			if strings.Contains(v, ",") {
				return false, fmt.Errorf("@scope can't use comma(,) get=" + v)
			}
		}
	}

	return strings.Contains(scope, scopeAttrPrefix), nil
}

func getTagsFromComment(comment string) (tags []string) {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	attribute := strings.Fields(commentLine)[0]
	lineRemainder, lowerAttribute := strings.TrimSpace(commentLine[len(attribute):]), strings.ToLower(attribute)

	if lowerAttribute == tagsAttr {
		for _, tag := range strings.Split(lineRemainder, ",") {
			tags = append(tags, strings.TrimSpace(tag))