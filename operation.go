
package swag

import (
	"encoding/json"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	"golang.org/x/tools/go/loader"
)

// RouteProperties describes HTTP properties of a single router comment.
type RouteProperties struct {
	HTTPMethod string
	Path       string
}

// Operation describes a single API operation on a path.
// For more information: https://github.com/swaggo/swag#api-operation
type Operation struct {
	parser              *Parser
	codeExampleFilesDir string
	spec.Operation
	RouterProperties []RouteProperties
}

var mimeTypeAliases = map[string]string{
	"json":                  "application/json",
	"xml":                   "text/xml",
	"plain":                 "text/plain",
	"html":                  "text/html",
	"mpfd":                  "multipart/form-data",
	"x-www-form-urlencoded": "application/x-www-form-urlencoded",
	"json-api":              "application/vnd.api+json",
	"json-stream":           "application/x-json-stream",
	"octet-stream":          "application/octet-stream",
	"png":                   "image/png",
	"jpeg":                  "image/jpeg",
	"gif":                   "image/gif",
}

var mimeTypePattern = regexp.MustCompile("^[^/]+/[^/]+$")

// NewOperation creates a new Operation with default properties.
// map[int]Response.
func NewOperation(parser *Parser, options ...func(*Operation)) *Operation {
	if parser == nil {
		parser = New()
	}

	result := &Operation{
		parser:           parser,
		RouterProperties: []RouteProperties{},
		Operation: spec.Operation{
			OperationProps: spec.OperationProps{
				ID:           "",
				Description:  "",
				Summary:      "",
				Security:     nil,
				ExternalDocs: nil,
				Deprecated:   false,
				Tags:         []string{},
				Consumes:     []string{},
				Produces:     []string{},
				Schemes:      []string{},
				Parameters:   []spec.Parameter{},
				Responses: &spec.Responses{
					VendorExtensible: spec.VendorExtensible{
						Extensions: spec.Extensions{},
					},
					ResponsesProps: spec.ResponsesProps{
						Default:             nil,
						StatusCodeResponses: make(map[int]spec.Response),
					},
				},
			},
			VendorExtensible: spec.VendorExtensible{
				Extensions: spec.Extensions{},
			},
		},
		codeExampleFilesDir: "",
	}

	for _, option := range options {
		option(result)
	}

	return result
}

// SetCodeExampleFilesDirectory sets the directory to search for codeExamples.
func SetCodeExampleFilesDirectory(directoryPath string) func(*Operation) {
	return func(o *Operation) {
		o.codeExampleFilesDir = directoryPath
	}
}

// ParseComment parses comment for given comment string and returns error if error occurs.
func (operation *Operation) ParseComment(comment string, astFile *ast.File) error {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	fields := FieldsByAnySpace(commentLine, 2)
	attribute := fields[0]
	lowerAttribute := strings.ToLower(attribute)
	var lineRemainder string
	if len(fields) > 1 {
		lineRemainder = fields[1]
	}
	switch lowerAttribute {
	case descriptionAttr:
		operation.ParseDescriptionComment(lineRemainder)
	case descriptionMarkdownAttr:
		commentInfo, err := getMarkdownForTag(lineRemainder, operation.parser.markdownFileDir)
		if err != nil {
			return err
		}

		operation.ParseDescriptionComment(string(commentInfo))
	case summaryAttr:
		operation.Summary = lineRemainder
	case idAttr:
		operation.ID = lineRemainder
	case tagsAttr:
		operation.ParseTagsComment(lineRemainder)
	case acceptAttr:
		return operation.ParseAcceptComment(lineRemainder)
	case produceAttr:
		return operation.ParseProduceComment(lineRemainder)
	case paramAttr:
		return operation.ParseParamComment(lineRemainder, astFile)
	case successAttr, failureAttr, responseAttr:
		return operation.ParseResponseComment(lineRemainder, astFile)
	case headerAttr:
		return operation.ParseResponseHeaderComment(lineRemainder, astFile)
	case routerAttr:
		return operation.ParseRouterComment(lineRemainder)
	case securityAttr:
		return operation.ParseSecurityComment(lineRemainder)
	case deprecatedAttr:
		operation.Deprecate()
	case xCodeSamplesAttr:
		return operation.ParseCodeSample(attribute, commentLine, lineRemainder)
	default:
		return operation.ParseMetadata(attribute, lowerAttribute, lineRemainder)
	}

	return nil
}

// ParseCodeSample godoc.
func (operation *Operation) ParseCodeSample(attribute, _, lineRemainder string) error {
	if lineRemainder == "file" {
		data, err := getCodeExampleForSummary(operation.Summary, operation.codeExampleFilesDir)
		if err != nil {
			return err
		}

		var valueJSON interface{}

		err = json.Unmarshal(data, &valueJSON)
		if err != nil {
			return fmt.Errorf("annotation %s need a valid json value", attribute)
		}

		// don't use the method provided by spec lib, because it will call toLower() on attribute names, which is wrongly
		operation.Extensions[attribute[1:]] = valueJSON

		return nil
	}

	// Fallback into existing logic
	return operation.ParseMetadata(attribute, strings.ToLower(attribute), lineRemainder)
}

// ParseDescriptionComment godoc.
func (operation *Operation) ParseDescriptionComment(lineRemainder string) {
	if operation.Description == "" {
		operation.Description = lineRemainder

		return
	}

	operation.Description += "\n" + lineRemainder
}

// ParseMetadata godoc.
func (operation *Operation) ParseMetadata(attribute, lowerAttribute, lineRemainder string) error {
	// parsing specific meta data extensions
	if strings.HasPrefix(lowerAttribute, "@x-") {
		if len(lineRemainder) == 0 {
			return fmt.Errorf("annotation %s need a value", attribute)
		}

		var valueJSON interface{}

		err := json.Unmarshal([]byte(lineRemainder), &valueJSON)
		if err != nil {
			return fmt.Errorf("annotation %s need a valid json value", attribute)
		}

		// don't use the method provided by spec lib, because it will call toLower() on attribute names, which is wrongly
		operation.Extensions[attribute[1:]] = valueJSON
	}

	return nil
}

var paramPattern = regexp.MustCompile(`(\S+)\s+(\w+)\s+([\S. ]+?)\s+(\w+)\s+"([^"]+)"`)

func findInSlice(arr []string, target string) bool {
	for _, str := range arr {
		if str == target {
			return true
		}
	}

	return false
}

// ParseParamComment parses params return []string of param properties
// E.g. @Param	queryText		formData	      string	  true		        "The email for login"
//
//	[param name]    [paramType] [data type]  [is mandatory?]   [Comment]
//
// E.g. @Param   some_id     path    int     true        "Some ID".
func (operation *Operation) ParseParamComment(commentLine string, astFile *ast.File) error {
	matches := paramPattern.FindStringSubmatch(commentLine)
	if len(matches) != 6 {
		return fmt.Errorf("missing required param comment parameters \"%s\"", commentLine)
	}

	name := matches[1]
	paramType := matches[2]
	refType := TransToValidSchemeType(matches[3])

	// Detect refType
	objectType := OBJECT

	if strings.HasPrefix(refType, "[]") {
		objectType = ARRAY
		refType = strings.TrimPrefix(refType, "[]")
		refType = TransToValidSchemeType(refType)
	} else if IsPrimitiveType(refType) ||
		paramType == "formData" && refType == "file" {
		objectType = PRIMITIVE
	}

	var enums []interface{}
	if !IsPrimitiveType(refType) {
		schema, _ := operation.parser.getTypeSchema(refType, astFile, false)
		if schema != nil && len(schema.Type) == 1 && schema.Enum != nil {
			if objectType == OBJECT {
				objectType = PRIMITIVE
			}
			refType = TransToValidSchemeType(schema.Type[0])
			enums = schema.Enum
		}
	}

	requiredText := strings.ToLower(matches[4])
	required := requiredText == "true" || requiredText == requiredLabel
	description := matches[5]

	param := createParameter(paramType, description, name, objectType, refType, required, enums, operation.parser.collectionFormatInQuery)

	switch paramType {
	case "path", "header":
		switch objectType {
		case ARRAY:
			if !IsPrimitiveType(refType) {
				return fmt.Errorf("%s is not supported array type for %s", refType, paramType)
			}
		case OBJECT:
			return fmt.Errorf("%s is not supported type for %s", refType, paramType)
		}
	case "query", "formData":
		switch objectType {
		case ARRAY:
			if !IsPrimitiveType(refType) && !(refType == "file" && paramType == "formData") {
				return fmt.Errorf("%s is not supported array type for %s", refType, paramType)
			}
		case PRIMITIVE:
			break
		case OBJECT:
			schema, err := operation.parser.getTypeSchema(refType, astFile, false)
			if err != nil {
				return err
			}

			if len(schema.Properties) == 0 {
				return nil
			}

			items := schema.Properties.ToOrderedSchemaItems()

			for _, item := range items {
				name, prop := item.Name, item.Schema
				if len(prop.Type) == 0 {
					continue
				}

				switch {
				case prop.Type[0] == ARRAY && prop.Items.Schema != nil &&
					len(prop.Items.Schema.Type) > 0 && IsSimplePrimitiveType(prop.Items.Schema.Type[0]):

					param = createParameter(paramType, prop.Description, name, prop.Type[0], prop.Items.Schema.Type[0], findInSlice(schema.Required, name), enums, operation.parser.collectionFormatInQuery)

				case IsSimplePrimitiveType(prop.Type[0]):
					param = createParameter(paramType, prop.Description, name, PRIMITIVE, prop.Type[0], findInSlice(schema.Required, name), enums, operation.parser.collectionFormatInQuery)
				default:
					operation.parser.debug.Printf("skip field [%s] in %s is not supported type for %s", name, refType, paramType)

					continue
				}

				param.Nullable = prop.Nullable
				param.Format = prop.Format
				param.Default = prop.Default
				param.Example = prop.Example
				param.Extensions = prop.Extensions
				param.CommonValidations.Maximum = prop.Maximum
				param.CommonValidations.Minimum = prop.Minimum
				param.CommonValidations.ExclusiveMaximum = prop.ExclusiveMaximum
				param.CommonValidations.ExclusiveMinimum = prop.ExclusiveMinimum
				param.CommonValidations.MaxLength = prop.MaxLength
				param.CommonValidations.MinLength = prop.MinLength
				param.CommonValidations.Pattern = prop.Pattern
				param.CommonValidations.MaxItems = prop.MaxItems
				param.CommonValidations.MinItems = prop.MinItems
				param.CommonValidations.UniqueItems = prop.UniqueItems
				param.CommonValidations.MultipleOf = prop.MultipleOf
				param.CommonValidations.Enum = prop.Enum
				operation.Operation.Parameters = append(operation.Operation.Parameters, param)
			}

			return nil
		}
	case "body":
		if objectType == PRIMITIVE {
			param.Schema = PrimitiveSchema(refType)
		} else {
			schema, err := operation.parseAPIObjectSchema(commentLine, objectType, refType, astFile)
			if err != nil {
				return err
			}

			param.Schema = schema
		}
	default:
		return fmt.Errorf("%s is not supported paramType", paramType)
	}

	err := operation.parseParamAttribute(commentLine, objectType, refType, &param)
	if err != nil {
		return err
	}

	operation.Operation.Parameters = append(operation.Operation.Parameters, param)

	return nil
}

const (
	jsonTag             = "json"
	bindingTag          = "binding"
	defaultTag          = "default"
	enumsTag            = "enums"
	exampleTag          = "example"
	schemaExampleTag    = "schemaExample"
	formatTag           = "format"
	validateTag         = "validate"
	minimumTag          = "minimum"
	maximumTag          = "maximum"
	minLengthTag        = "minLength"
	maxLengthTag        = "maxLength"
	multipleOfTag       = "multipleOf"
	readOnlyTag         = "readonly"
	extensionsTag       = "extensions"
	collectionFormatTag = "collectionFormat"
)

var regexAttributes = map[string]*regexp.Regexp{
	// for Enums(A, B)
	enumsTag: regexp.MustCompile(`(?i)\s+enums\(.*\)`),
	// for maximum(0)
	maximumTag: regexp.MustCompile(`(?i)\s+maxinum|maximum\(.*\)`),
	// for minimum(0)
	minimumTag: regexp.MustCompile(`(?i)\s+mininum|minimum\(.*\)`),
	// for default(0)
	defaultTag: regexp.MustCompile(`(?i)\s+default\(.*\)`),
	// for minlength(0)
	minLengthTag: regexp.MustCompile(`(?i)\s+minlength\(.*\)`),
	// for maxlength(0)
	maxLengthTag: regexp.MustCompile(`(?i)\s+maxlength\(.*\)`),
	// for format(email)
	formatTag: regexp.MustCompile(`(?i)\s+format\(.*\)`),
	// for extensions(x-example=test)
	extensionsTag: regexp.MustCompile(`(?i)\s+extensions\(.*\)`),
	// for collectionFormat(csv)
	collectionFormatTag: regexp.MustCompile(`(?i)\s+collectionFormat\(.*\)`),
	// example(0)
	exampleTag: regexp.MustCompile(`(?i)\s+example\(.*\)`),
	// schemaExample(0)
	schemaExampleTag: regexp.MustCompile(`(?i)\s+schemaExample\(.*\)`),
}

func (operation *Operation) parseParamAttribute(comment, objectType, schemaType string, param *spec.Parameter) error {
	schemaType = TransToValidSchemeType(schemaType)

	for attrKey, re := range regexAttributes {
		attr, err := findAttr(re, comment)
		if err != nil {
			continue
		}

		switch attrKey {
		case enumsTag:
			err = setEnumParam(param, attr, objectType, schemaType)
		case minimumTag, maximumTag:
			err = setNumberParam(param, attrKey, schemaType, attr, comment)
		case defaultTag:
			err = setDefault(param, schemaType, attr)
		case minLengthTag, maxLengthTag:
			err = setStringParam(param, attrKey, schemaType, attr, comment)
		case formatTag:
			param.Format = attr
		case exampleTag:
			err = setExample(param, schemaType, attr)
		case schemaExampleTag:
			err = setSchemaExample(param, schemaType, attr)
		case extensionsTag:
			param.Extensions = setExtensionParam(attr)
		case collectionFormatTag:
			err = setCollectionFormatParam(param, attrKey, objectType, attr, comment)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func findAttr(re *regexp.Regexp, commentLine string) (string, error) {
	attr := re.FindString(commentLine)

	l, r := strings.Index(attr, "("), strings.Index(attr, ")")
	if l == -1 || r == -1 {
		return "", fmt.Errorf("can not find regex=%s, comment=%s", re.String(), commentLine)
	}

	return strings.TrimSpace(attr[l+1 : r]), nil
}

func setStringParam(param *spec.Parameter, name, schemaType, attr, commentLine string) error {
	if schemaType != STRING {
		return fmt.Errorf("%s is attribute to set to a number. comment=%s got=%s", name, commentLine, schemaType)
	}

	n, err := strconv.ParseInt(attr, 10, 64)
	if err != nil {
		return fmt.Errorf("%s is allow only a number got=%s", name, attr)
	}

	switch name {
	case minLengthTag:
		param.MinLength = &n
	case maxLengthTag:
		param.MaxLength = &n
	}

	return nil
}

func setNumberParam(param *spec.Parameter, name, schemaType, attr, commentLine string) error {
	switch schemaType {
	case INTEGER, NUMBER:
		n, err := strconv.ParseFloat(attr, 64)
		if err != nil {
			return fmt.Errorf("maximum is allow only a number. comment=%s got=%s", commentLine, attr)
		}

		switch name {
		case minimumTag:
			param.Minimum = &n
		case maximumTag:
			param.Maximum = &n
		}

		return nil
	default:
		return fmt.Errorf("%s is attribute to set to a number. comment=%s got=%s", name, commentLine, schemaType)
	}
}

func setEnumParam(param *spec.Parameter, attr, objectType, schemaType string) error {
	for _, e := range strings.Split(attr, ",") {
		e = strings.TrimSpace(e)

		value, err := defineType(schemaType, e)
		if err != nil {
			return err
		}

		switch objectType {
		case ARRAY:
			param.Items.Enum = append(param.Items.Enum, value)
		default:
			param.Enum = append(param.Enum, value)
		}
	}

	return nil
}

func setExtensionParam(attr string) spec.Extensions {
	extensions := spec.Extensions{}

	for _, val := range splitNotWrapped(attr, ',') {
		parts := strings.SplitN(val, "=", 2)
		if len(parts) == 2 {
			extensions.Add(parts[0], parts[1])

			continue
		}

		if len(parts[0]) > 0 && string(parts[0][0]) == "!" {
			extensions.Add(parts[0][1:], false)

			continue
		}

		extensions.Add(parts[0], true)
	}

	return extensions
}

func setCollectionFormatParam(param *spec.Parameter, name, schemaType, attr, commentLine string) error {
	if schemaType == ARRAY {
		param.CollectionFormat = TransToValidCollectionFormat(attr)

		return nil
	}

	return fmt.Errorf("%s is attribute to set to an array. comment=%s got=%s", name, commentLine, schemaType)
}

func setDefault(param *spec.Parameter, schemaType string, value string) error {
	val, err := defineType(schemaType, value)
	if err != nil {
		return nil // Don't set a default value if it's not valid
	}

	param.Default = val

	return nil
}

func setSchemaExample(param *spec.Parameter, schemaType string, value string) error {
	val, err := defineType(schemaType, value)
	if err != nil {
		return nil // Don't set a example value if it's not valid
	}
	// skip schema
	if param.Schema == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		//  replaces \r \n \t in example string values.
		param.Schema.Example = strings.NewReplacer(`\r`, "\r", `\n`, "\n", `\t`, "\t").Replace(v)
	default:
		param.Schema.Example = val
	}

	return nil
}

func setExample(param *spec.Parameter, schemaType string, value string) error {
	val, err := defineType(schemaType, value)
	if err != nil {
		return nil // Don't set a example value if it's not valid
	}

	param.Example = val

	return nil
}

// defineType enum value define the type (object and array unsupported).
func defineType(schemaType string, value string) (v interface{}, err error) {
	schemaType = TransToValidSchemeType(schemaType)

	switch schemaType {
	case STRING:
		return value, nil
	case NUMBER:
		v, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	case INTEGER:
		v, err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	case BOOLEAN:
		v, err = strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("enum value %s can't convert to %s err: %s", value, schemaType, err)
		}
	default: