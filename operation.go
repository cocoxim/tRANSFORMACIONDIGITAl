
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