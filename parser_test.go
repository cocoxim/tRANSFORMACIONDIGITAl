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

		expected := "docs/markdown"
		p := New(SetMarkdownFileDirectory(expected))
		assert.Equal(t, expected, p.markdownFileDir)
	})

	t.Run("SetCodeExamplesDirectory", func(t *testing.T) {
		t.Parallel()

		expected := "docs/examples"
		p := New(SetCodeExamplesDirectory(expected))
		assert.Equal(t, expected, p.codeExampleFilesDir)
	})

	t.Run("SetStrict", func(t *testing.T) {
		t.Parallel()

		p := New()
		assert.Equal(t, false, p.Strict)

		p = New(SetStrict(true))
		assert.Equal(t, true, p.Strict)
	})

	t.Run("SetDebugger", func(t *testing.T) {
		t.Parallel()

		logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

		p := New(SetDebugger(logger))
		assert.Equal(t, logger, p.debug)
	})

	t.Run("SetFieldParserFactory", func(t *testing.T) {
		t.Parallel()

		p := New(SetFieldParserFactory(nil))
		assert.Nil(t, p.fieldParserFactory)
	})
}

func TestSetOverrides(t *testing.T) {
	t.Parallel()

	overrides := map[string]string{
		"foo": "bar",
	}

	p := New(SetOverrides(overrides))
	assert.Equal(t, overrides, p.Overrides)
}

func TestOverrides_getTypeSchema(t *testing.T) {
	t.Parallel()

	overrides := map[string]string{
		"sql.NullString": "string",
	}

	p := New(SetOverrides(overrides))

	t.Run("Override sql.NullString by string", func(t *testing.T) {
		t.Parallel()

		s, err := p.getTypeSchema("sql.NullString", nil, false)
		if assert.NoError(t, err) {
			assert.Truef(t, s.Type.Contains("string"), "type sql.NullString should be overridden by string")
		}
	})

	t.Run("Missing Override for sql.NullInt64", func(t *testing.T) {
		t.Parallel()

		_, err := p.getTypeSchema("sql.NullInt64", nil, false)
		if assert.Error(t, err) {
			assert.Equal(t, "cannot find type definition: sql.NullInt64", err.Error())
		}
	})
}

func TestParser_ParseDefinition(t *testing.T) {
	p := New()

	// Parsing existing type
	definition := &TypeSpecDef{
		PkgPath: "github.com/swagger/swag",
		File: &ast.File{
			Name: &ast.Ident{
				Name: "swag",
			},
		},
		TypeSpec: &ast.TypeSpec{
			Name: &ast.Ident{
				Name: "Test",
			},
		},
	}

	expected := &Schema{}
	p.parsedSchemas[definition] = expected

	schema, err := p.ParseDefinition(definition)
	assert.NoError(t, err)
	assert.Equal(t, expected, schema)

	// Parsing *ast.FuncType
	definition = &TypeSpecDef{
		PkgPath: "github.com/swagger/swag/model",
		File: &ast.File{
			Name: &ast.Ident{
				Name: "model",
			},
		},
		TypeSpec: &ast.TypeSpec{
			Name: &ast.Ident{
				Name: "Test",
			},
			Type: &ast.FuncType{},
		},
	}
	_, err = p.ParseDefinition(definition)
	assert.Error(t, err)

	// Parsing *ast.FuncType with parent spec
	definition = &TypeSpecDef{
		PkgPath: "github.com/swagger/swag/model",
		File: &ast.File{
			Name: &ast.Ident{
				Name: "model",
			},
		},
		TypeSpec: &ast.TypeSpec{
			Name: &ast.Ident{
				Name: "Test",
			},
			Type: &ast.FuncType{},
		},
		ParentSpec: &ast.FuncDecl{
			Name: ast.NewIdent("TestFuncDecl"),
		},
	}
	_, err = p.ParseDefinition(definition)
	assert.Error(t, err)
	assert.Equal(t, "model.TestFuncDecl.Test", definition.TypeName())
}

func TestParser_ParseGeneralApiInfo(t *testing.T) {
	t.Parallel()

	expected := `{
    "schemes": [
        "http",
        "https"
    ],
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.\nIt has a lot of beautiful features.",
        "title": "Swagger Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0",
        "x-logo": {
            "altText": "Petstore logo",
            "backgroundColor": "#FFFFFF",
            "url": "https://redocly.github.io/redoc/petstore-logo.png"
        }
    },
    "host": "petstore.swagger.io",
    "basePath": "/v2",
    "paths": {},
    "securityDefinitions": {
        "ApiKeyAuth": {
            "description": "some description",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            },
            "x-tokenname": "id_token"
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            },
            "x-google-audiences": "some_audience.google.com"
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api"
    },
    "x-google-endpoints": [
        {
            "allowCors": true,
            "name": "name.endpoints.environment.cloud.goog"
        }
    ],
    "x-google-marks": "marks values"
}`
	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)

	p := New()

	err := p.ParseGeneralAPIInfo("testdata/main.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralApiInfoTemplated(t *testing.T) {
	t.Parallel()

	expected := `{
    "swagger": "2.0",
    "info": {
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        }
    },
    "paths": {},
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        },
        "BasicAuth": {
            "type": "basic"
        },
        "OAuth2AccessCode": {
            "type": "oauth2",
            "flow": "accessCode",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information"
            }
        },
        "OAuth2Application": {
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Implicit": {
            "type": "oauth2",
            "flow": "implicit",
            "authorizationUrl": "https://example.com/oauth/authorize",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "write": " Grants write access"
            }
        },
        "OAuth2Password": {
            "type": "oauth2",
            "flow": "password",
            "tokenUrl": "https://example.com/oauth/token",
            "scopes": {
                "admin": " Grants read and write access to administrative information",
                "read": " Grants read access",
                "write": " Grants write access"
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api"
    },
    "x-google-endpoints": [
        {
            "allowCors": true,
            "name": "name.endpoints.environment.cloud.goog"
        }
    ],
    "x-google-marks": "marks values"
}`
	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)

	p := New()

	err := p.ParseGeneralAPIInfo("testdata/templated.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralApiInfoExtensions(t *testing.T) {
	// should return an error because extension value is not a valid json
	t.Run("Test invalid extension value", func(t *testing.T) {
		t.Parallel()

		expected := "annotation @x-google-endpoints need a valid json value"
		gopath := os.Getenv("GOPATH")
		assert.NotNil(t, gopath)

		p := New()

		err := p.ParseGeneralAPIInfo("testdata/extensionsFail1.go")
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	})

	// should return an error because extension don't have a value
	t.Run("Test missing extension value", func(t *testing.T) {
		t.Parallel()

		expected := "annotation @x-google-endpoints need a value"
		gopath := os.Getenv("GOPATH")
		assert.NotNil(t, gopath)

		p := New()

		err := p.ParseGeneralAPIInfo("testdata/extensionsFail2.go")
		if assert.Error(t, err) {
			assert.Equal(t, expected, err.Error())
		}
	})
}

func TestParser_ParseGeneralApiInfoWithOpsInSameFile(t *testing.T) {
	t.Parallel()

	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "This is a sample server Petstore server.\nIt has a lot of beautiful features.",
        "title": "Swagger Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "1.0"
    },
    "paths": {}
}`

	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)

	p := New()

	err := p.ParseGeneralAPIInfo("testdata/single_file_api/main.go")
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))
}

func TestParser_ParseGeneralAPIInfoMarkdown(t *testing.T) {
	t.Parallel()

	p := New(SetMarkdownFileDirectory("testdata"))
	mainAPIFile := "testdata/markdown.go"
	err := p.ParseGeneralAPIInfo(mainAPIFile)
	assert.NoError(t, err)

	expected := `{
    "swagger": "2.0",
    "info": {
        "description": "Swagger Example API Markdown Description",
        "title": "Swagger Example API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "version": "1.0"
    },
    "paths": {},
    "tags": [
        {
            "description": "Users Tag Markdown Description",
            "name": "users"
        }
    ]
}`
	b, _ := json.MarshalIndent(p.swagger, "", "    ")
	assert.Equal(t, expected, string(b))

	p = New()

	err = p.ParseGeneralAPIInfo(mainAPIFile)
	assert.Error(t, err)
}

func TestParser_ParseGeneralApiInfoFailed(t *testing.T) {
	t.Parallel()

	gopath := os.Getenv("GOPATH")
	assert.NotNil(t, gopath)
	p := New()
	assert.Error(t, p.ParseGeneralAPIInfo("testdata/noexist.go"))
}

func TestParser_ParseAcceptComment(t *testing.T) {
	t.Parallel()

	expected := []string{
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
		"application/xhtml+xml",
		"application/health+json",
	}

	comment := `@Accept json,xml,plain,html,mpfd,x-www-form-urlencoded,json-api,json-stream,octet-stream,png,jpeg,gif,application/xhtml+xml,application/health+json`

	parser := New()
	assert.NoError(t, parseGeneralAPIInfo(parser, []string{comment}))
	assert.Equal(t, parser.swagger.Consumes, expected)

	assert.Error(t, parseGeneralAPIInfo(parser, []string{`@Accept cookies,candies`}))

	parser = New()
	assert.NoError(t, parser.ParseAcceptComment(comment[len(acceptAttr)+1:]))
	assert.Equal(t, parser.swagger.Consumes, expected)
}

func TestParser_ParseProduceComment(t *testing.T) {
	t.Parallel()

	expected := []string{
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
		"application/xhtml+xml",
		"application/health+json",
	}

	comment := `@Produce json,xml,plain,html,mpfd,x-www-form-urlencoded,json-api,json-stream,octet-stream,png,jpeg,gif,application/xhtml+xml,application/health+json`

	parser := New()
	assert.NoError(t, parseGeneralAPIInfo(parser, []string{comment}))
	assert.Equal(t, parser.swagger.Produces, expected)

	assert.Error(t, parseGeneralAPIInfo(parser, []string{`@Produce cookies,candies`}))

	parser = New()
	assert.NoError(t, parser.ParseProduceComment(comment[len(produceAttr)+1:]))
	assert.Equal(t, parser.swagger.Produces, expected)
}

func TestParser_ParseGeneralAPIInfoCollectionFormat(t *testing.T) {
	t.Parallel()

	parser := New()
	assert.NoError(t, parseGeneralAPIInfo(parser, []string{
		"@query.collection.format csv",
	}))
	assert.Equal(t, parser.collectionFormatInQuery, "csv")

	assert.NoError(t, parseGeneralAPIInfo(parser, []string{
		"@query.collection.format tsv",
	}))
	assert.Equal(t, parser.collectionFormatInQuery, "tsv")
}

func TestParser_ParseGeneralAPITagGroups(t *testing.T) {
	t.Parallel()

	parser := New()
	assert.NoError(t, parseGeneralAPIInfo(parser, []string{
		"@x-tagGroups [{\"name\":\"General\",\"tags\":[\"lanes\",\"video-recommendations\"]}]",
	}))

	expected := []interface{}{map[string]interface{}{"name": "General", "tags": []interface{}{"lanes", "video-recommendations"}}}
	assert.Equal(t, parser.swagger.Extensions["x-tagGroups"], expected)
}

func 