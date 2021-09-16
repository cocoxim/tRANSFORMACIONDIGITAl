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

		ex