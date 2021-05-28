
package swag

import (
	"bytes"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

// Check of @Param @Success @Failure @Response @Header
var specialTagForSplit = map[string]bool{
	paramAttr:    true,
	successAttr:  true,
	failureAttr:  true,
	responseAttr: true,
	headerAttr:   true,
}

var skipChar = map[byte]byte{
	'"': '"',
	'(': ')',
	'{': '}',
	'[': ']',
}

// Formatter implements a formatter for Go source files.
type Formatter struct {
	// debugging output goes here
	debug Debugger
}

// NewFormatter create a new formatter instance.
func NewFormatter() *Formatter {
	formatter := &Formatter{
		debug: log.New(os.Stdout, "", log.LstdFlags),
	}
	return formatter
}

// Format formats swag comments in contents. It uses fileName to report errors
// that happen during parsing of contents.
func (f *Formatter) Format(fileName string, contents []byte) ([]byte, error) {
	fileSet := token.NewFileSet()
	ast, err := goparser.ParseFile(fileSet, fileName, contents, goparser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Formatting changes are described as an edit list of byte range
	// replacements. We make these content-level edits directly rather than
	// changing the AST nodes and writing those out (via [go/printer] or
	// [go/format]) so that we only change the formatting of Swag attribute
	// comments. This won't touch the formatting of any other comments, or of
	// functions, etc.
	maxEdits := 0
	for _, comment := range ast.Comments {
		maxEdits += len(comment.List)
	}
	edits := make(edits, 0, maxEdits)

	for _, comment := range ast.Comments {
		formatFuncDoc(fileSet, comment.List, &edits)
	}

	return edits.apply(contents), nil
}

type edit struct {
	begin       int
	end         int
	replacement []byte
}

type edits []edit

func (edits edits) apply(contents []byte) []byte {
	// Apply the edits with the highest offset first, so that earlier edits
	// don't affect the offsets of later edits.
	sort.Slice(edits, func(i, j int) bool {
		return edits[i].begin > edits[j].begin