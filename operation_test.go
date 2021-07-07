package swag

import (
	"encoding/json"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestParseEmptyComment(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)
	err := operation.ParseComment("//", nil)

	assert.NoError(t, err)
}

func TestParseTagsComment(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)
	err := operation.ParseComment(`/@Tags pet, store,user`, nil)
	assert.NoError(t, err)
	assert.Equal(t, operation.Tags, []string{"pet", "store", "user"})
}

func TestParseAcceptComment(t *testing.T) {
	t.Parallel()

	comment := `/@Accept json,xml,plain,html,mpfd,x-www-form-urlencoded,json-api,json-stream,octet-stream,png,jpeg,gif,application/xhtml+xml,application/health+json`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Equal(t,
		operation.Consumes,
		[]string{"application/json",
			"text/xml",
			"text/plain",
			"text/html",
			"multipart/form-data",
			"application/x-www-form-urlencoded",
			"application/vnd.api+json",
			"application/x-json-stream",
			"application/octet-stream",
			"image/png",
			"image/jpeg",
			"image/gif",
			"application/xhtml+xml",
			"application/health+json"})
}

func TestParseAcceptCommentErr(t *testing.T) {
	t.Parallel()

	comment := `/@Accept unknown`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseProduceComment(t *testing.T) {
	t.Parallel()

	expected := `{
    "produces": [
        "application/json",
        "text/xml",
        "text/plain",
        "text/html",
        "multipart/form-data",
        "application/x-www-form-urlencoded",
        "application/vnd.api+json",
        "application/x-json-stream",
		"application/octet-stream",
		"image/png",
		"image/jpeg",
		"image/gif",
		"application/health+json"
    ]
}`
	comment := `/