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
	comment := `/@Produce json,xml,plain,html,mpfd,x-www-form-urlencoded,json-api,json-stream,octet-stream,png,jpeg,gif,application/health+json`
	operation := new(Operation)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")
	b, _ := json.MarshalIndent(operation, "", "    ")
	assert.JSONEq(t, expected, string(b))
}

func TestParseProduceCommentErr(t *testing.T) {
	t.Parallel()

	operation := new(Operation)
	err := operation.ParseComment("/@Produce foo", nil)
	assert.Error(t, err)
}

func TestParseRouterComment(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{wishlist_id} [get]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.RouterProperties, 1)
	assert.Equal(t, "/customer/get-wishlist/{wishlist_id}", operation.RouterProperties[0].Path)
	assert.Equal(t, "GET", operation.RouterProperties[0].HTTPMethod)

	comment = `/@Router /customer/get-wishlist/{wishlist_id} [unknown]`
	operation = NewOperation(nil)
	err = operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseRouterMultipleComments(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{wishlist_id} [get]`
	anotherComment := `/@Router /customer/get-the-wishlist/{wishlist_id} [post]`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	err = operation.ParseComment(anotherComment, nil)
	assert.NoError(t, err)

	assert.Len(t, operation.RouterProperties, 2)
	assert.Equal(t, "/customer/get-wishlist/{wishlist_id}", operation.RouterProperties[0].Path)
	assert.Equal(t, "GET", operation.RouterProperties[0].HTTPMethod)
	assert.Equal(t, "/customer/get-the-wishlist/{wishlist_id}", operation.RouterProperties[1].Path)
	assert.Equal(t, "POST", operation.RouterProperties[1].HTTPMethod)
}

func TestParseRouterOnlySlash(t *testing.T) {
	t.Parallel()

	comment := `// @Router / [get]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.RouterProperties, 1)
	assert.Equal(t, "/", operation.RouterProperties[0].Path)
	assert.Equal(t, "GET", operation.RouterProperties[0].HTTPMethod)
}

func TestParseRouterCommentWithPlusSign(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{proxy+} [post]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.RouterProperties, 1)
	assert.Equal(t, "/customer/get-wishlist/{proxy+}", operation.RouterProperties[0].Path)
	assert.Equal(t, "POST", operation.RouterProperties[0].HTTPMethod)
}

func TestParseRouterCommentWithDollarSign(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{wishlist_id}$move [post]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.RouterProperties, 1)
	assert.Equal(t, "/customer/get-wishlist/{wishlist_id}$move", operation.RouterProperties[0].Path)
	assert.Equal(t, "POST", operation.RouterProperties[0].HTTPMethod)
}

func TestParseRouterCommentNoDollarSignAtPathStartErr(t *testing.T) {
	t.Parallel()

	comment := `/@Router $customer/get-wishlist/{wishlist_id}$move [post]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseRouterCommentWithColonSign(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{wishlist_id}:move [post]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.RouterProperties, 1)
	assert.Equal(t, "/customer/get-wishlist/{wishlist_id}:move", operation.RouterProperties[0].Path)
	assert.Equal(t, "POST", operation.RouterProperties[0].HTTPMethod)
}

func TestParseRouterCommentNoColonSignAtPathStartErr(t *testing.T) {
	t.Parallel()

	comment := `/@Router :customer/get-wishlist/{wishlist_id}:move [post]`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseRouterCommentMethodSeparationErr(t *testing.T) {
	t.Parallel()

	comment := `/@Router /api/{id}|,*[get`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseRouterCommentMethodMissingErr(t *testing.T) {
	t.Parallel()

	comment := `/@Router /customer/get-wishlist/{wishlist_id}`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestOperation_ParseResponseWithDefault(t *testing.T) {
	t.Parallel()

	comment := `@Success default {object} nil "An empty response"`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	assert.Equal(t, "An empty response", operation.Responses.Default.Description)

	comment = `@Success 200,default {string} Response "A response"`
	operation = NewOperation(nil)

	err = operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	assert.Equal(t, "A response", operation.Responses.Default.Description)
	assert.Equal(t, "A response", operation.Responses.StatusCodeResponses[200].Description)
}

func TestParseResponseSuccessCommentWithEmptyResponse(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} nil "An empty response"`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `An empty response`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")
	expected := `{
    "responses": {
        "200": {
            "description": "An empty response"
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseFailureCommentWithEmptyResponse(t *testing.T) {
	t.Parallel()

	comment := `@Failure 500 {object} nil`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(operation, "", "    ")
	expected := `{
    "responses": {
        "500": {
            "description": "Internal Server Error"
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithObjectType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.OrderRow "Error message, if code != 200`
	operation := NewOperation(nil)
	operation.parser.addTestType("model.OrderRow")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "$ref": "#/definitions/model.OrderRow"
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedPrimitiveType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data=string,data2=int} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data": {
                                "type": "string"
                            },
                            "data2": {
                                "type": "integer"
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedPrimitiveArrayType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data=[]string,data2=[]int} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            },
                            "data2": {
                                "type": "array",
                                "items": {
                                    "type": "integer"
                                }
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedObjectType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data=model.Payload,data2=model.Payload2} "Error message, if code != 200`
	operation := NewOperation(nil)
	operation.parser.addTestType("model.CommonHeader")
	operation.parser.addTestType("model.Payload")
	operation.parser.addTestType("model.Payload2")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data": {
                                "$ref": "#/definitions/model.Payload"
                            },
                            "data2": {
                                "$ref": "#/definitions/model.Payload2"
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedArrayObjectType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data=[]model.Payload,data2=[]model.Payload2} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")
	operation.parser.addTestType("model.Payload")
	operation.parser.addTestType("model.Payload2")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/model.Payload"
                                }
                            },
                            "data2": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/model.Payload2"
                                }
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedFields(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data1=int,data2=[]int,data3=model.Payload,data4=[]model.Payload} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")
	operation.parser.addTestType("model.Payload")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data1": {
                                "type": "integer"
                            },
                            "data2": {
                                "type": "array",
                                "items": {
                                    "type": "integer"
                                }
                            },
                            "data3": {
                                "$ref": "#/definitions/model.Payload"
                            },
                            "data4": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/model.Payload"
                                }
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithDeepNestedFields(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.CommonHeader{data1=int,data2=[]int,data3=model.Payload{data1=int,data2=model.DeepPayload},data4=[]model.Payload{data1=[]int,data2=[]model.DeepPayload}} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")
	operation.parser.addTestType("model.Payload")
	operation.parser.addTestType("model.DeepPayload")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")
	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data1": {
                                "type": "integer"
                            },
                            "data2": {
                                "type": "array",
                                "items": {
                                    "type": "integer"
                                }
                            },
                            "data3": {
                                "allOf": [
                                    {
                                        "$ref": "#/definitions/model.Payload"
                                    },
                                    {
                                        "type": "object",
                                        "properties": {
                                            "data1": {
                                                "type": "integer"
                                            },
                                            "data2": {
                                                "$ref": "#/definitions/model.DeepPayload"
                                            }
                                        }
                                    }
                                ]
                            },
                            "data4": {
                                "type": "array",
                                "items": {
                                    "allOf": [
                                        {
                                            "$ref": "#/definitions/model.Payload"
                                        },
                                        {
                                            "type": "object",
                                            "properties": {
                                                "data1": {
                                                    "type": "array",
                                                    "items": {
                                                        "type": "integer"
                                                    }
                                                },
                                                "data2": {
                                                    "type": "array",
                                                    "items": {
                                                        "$ref": "#/definitions/model.DeepPayload"
                                                    }
                                                }
                                            }
                                        }
                                    ]
                                }
                            }
                        }
                    }
                ]
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithNestedArrayMapFields(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} []map[string]model.CommonHeader{data1=[]map[string]model.Payload,data2=map[string][]int} "Error message, if code != 200`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")
	operation.parser.addTestType("model.Payload")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")
	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "type": "array",
                "items": {
                    "type": "object",
                    "additionalProperties": {
                        "allOf": [
                            {
                                "$ref": "#/definitions/model.CommonHeader"
                            },
                            {
                                "type": "object",
                                "properties": {
                                    "data1": {
                                        "type": "array",
                                        "items": {
                                            "type": "object",
                                            "additionalProperties": {
                                                "$ref": "#/definitions/model.Payload"
                                            }
                                        }
                                    },
                                    "data2": {
                                        "type": "object",
                                        "additionalProperties": {
                                            "type": "array",
                                            "items": {
                                                "type": "integer"
                                            }
                                        }
                                    }
                                }
                            }
                        ]
                    }
                }
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithObjectTypeInSameFile(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} testOwner "Error message, if code != 200"`
	operation := NewOperation(nil)

	operation.parser.addTestType("swag.testOwner")

	fset := token.NewFileSet()
	astFile, err := goparser.ParseFile(fset, "operation_test.go", `package swag
	type testOwner struct {

	}
	`, goparser.ParseComments)
	assert.NoError(t, err)

	err = operation.ParseComment(comment, astFile)
	assert.NoError(t, err)

	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "$ref": "#/definitions/swag.testOwner"
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithObjectTypeAnonymousField(t *testing.T) {
	//TODO: test Anonymous
}

func TestParseResponseCommentWithObjectTypeErr(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {object} model.OrderRow "Error message, if code != 200"`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.notexist")

	err := operation.ParseComment(comment, nil)
	assert.Error(t, err)
}

func TestParseResponseCommentWithArrayType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {array} model.OrderRow "Error message, if code != 200`
	operation := NewOperation(nil)
	operation.parser.addTestType("model.OrderRow")
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	response := operation.Responses.StatusCodeResponses[200]
	assert.Equal(t, `Error message, if code != 200`, response.Description)
	assert.Equal(t, spec.StringOrArray{"array"}, response.Schema.Type)

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "Error message, if code != 200",
            "schema": {
                "type": "array",
                "items": {
                    "$ref": "#/definitions/model.OrderRow"
                }
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithBasicType(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 {string} string "it's ok'"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")
	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "it's ok'",
            "schema": {
                "type": "string"
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithBasicTypeAndCodes(t *testing.T) {
	t.Parallel()

	comment := `@Success 200,201,default {string} string "it's ok"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")
	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "it's ok",
            "schema": {
                "type": "string"
            }
        },
        "201": {
            "description": "it's ok",
            "schema": {
                "type": "string"
            }
        },
        "default": {
            "description": "it's ok",
            "schema": {
                "type": "string"
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseEmptyResponseComment(t *testing.T) {
	t.Parallel()

	comment := `@Success 200 "it is ok"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "it is ok"
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseEmptyResponseCommentWithCodes(t *testing.T) {
	t.Parallel()

	comment := `@Success 200,201,default "it is ok"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "it is ok"
        },
        "201": {
            "description": "it is ok"
        },
        "default": {
            "description": "it is ok"
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentWithHeader(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)
	err := operation.ParseComment(`@Success 200 "it's ok"`, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	err = operation.ParseComment(`@Header 200 {string} Token "qwerty"`, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, err := json.MarshalIndent(operation, "", "    ")
	assert.NoError(t, err)

	expected := `{
    "responses": {
        "200": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))

	err = operation.ParseComment(`@Header 200 "Mallformed"`, nil)
	assert.Error(t, err, "ParseComment should not fail")

	err = operation.ParseComment(`@Header 200,asdsd {string} Token "qwerty"`, nil)
	assert.Error(t, err, "ParseComment should not fail")
}

func TestParseResponseCommentWithHeaderForCodes(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)

	comment := `@Success 200,201,default "it's ok"`
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	comment = `@Header 200,201,default {string} Token "qwerty"`
	err = operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	comment = `@Header all {string} Token2 "qwerty"`
	err = operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, err := json.MarshalIndent(operation, "", "    ")
	assert.NoError(t, err)

	expected := `{
    "responses": {
        "200": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                },
                "Token2": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        },
        "201": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                },
                "Token2": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        },
        "default": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                },
                "Token2": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))

	comment = `@Header 200 "Mallformed"`
	err = operation.ParseComment(comment, nil)
	assert.Error(t, err, "ParseComment should not fail")
}

func TestParseResponseCommentWithHeaderOnlyAll(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)

	comment := `@Success 200,201,default "it's ok"`
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	comment = `@Header all {string} Token "qwerty"`
	err = operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, err := json.MarshalIndent(operation, "", "    ")
	assert.NoError(t, err)

	expected := `{
    "responses": {
        "200": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        },
        "201": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        },
        "default": {
            "description": "it's ok",
            "headers": {
                "Token": {
                    "type": "string",
                    "description": "qwerty"
                }
            }
        }
    }
}`
	assert.Equal(t, expected, string(b))

	comment = `@Header 200 "Mallformed"`
	err = operation.ParseComment(comment, nil)
	assert.Error(t, err, "ParseComment should not fail")
}

func TestParseEmptyResponseOnlyCode(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)
	err := operation.ParseComment(`@Success 200`, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "OK"
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseEmptyResponseOnlyCodes(t *testing.T) {
	t.Parallel()

	comment := `@Success 200,201,default`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err, "ParseComment should not fail")

	b, _ := json.MarshalIndent(operation, "", "    ")

	expected := `{
    "responses": {
        "200": {
            "description": "OK"
        },
        "201": {
            "description": "Created"
        },
        "default": {
            "description": ""
        }
    }
}`
	assert.Equal(t, expected, string(b))
}

func TestParseResponseCommentParamMissing(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)

	paramLenErrComment := `@Success notIntCode`
	paramLenErr := operation.ParseComment(paramLenErrComment, nil)
	assert.EqualError(t, paramLenErr, `can not parse response comment "notIntCode"`)

	paramLenErrComment = `@Success notIntCode {string} string "it ok"`
	paramLenErr = operation.ParseComment(paramLenErrComment, nil)
	assert.EqualError(t, paramLenErr, `can not parse response comment "notIntCode {string} string "it ok""`)

	paramLenErrComment = `@Success notIntCode "it ok"`
	paramLenErr = operation.ParseComment(paramLenErrComment, nil)
	assert.EqualError(t, paramLenErr, `can not parse response comment "notIntCode "it ok""`)
}

func TestOperation_ParseParamComment(t *testing.T) {
	t.Parallel()

	t.Run("integer", func(t *testing.T) {
		t.Parallel()
		for _, paramType := range []string{"header", "path", "query", "formData"} {
			t.Run(paramType, func(t *testing.T) {
				o := NewOperation(nil)
				err := o.ParseComment(`@Param some_id `+paramType+` int true "Some ID"`, nil)

				assert.NoError(t, err)
				assert.Equal(t, o.Parameters, []spec.Parameter{{
					SimpleSchema: spec.SimpleSchema{
						Type: "integer",
					},
					ParamProps: spec.ParamProps{
						Name:        "some_id",
						Description: "Some ID",
						In:          paramType,
						Required:    true,
					},
				}})
			})
		}
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		for _, paramType := range []string{"header", "path", "query", "formData"} {
			t.Run(paramType, func(t *testing.T) {
				o := NewOperation(nil)
				err := o.ParseComment(`@Param some_string `+paramType+` string true "Some String"`, nil)

				assert.NoError(t, err)
				assert.Equal(t, o.Parameters, []spec.Parameter{{
					SimpleSchema: spec.SimpleSchema{
						Type: "string",
					},
					ParamProps: spec.ParamProps{
						Name:        "some_string",
						Description: "Some String",
						In:          paramType,
						Required:    true,
					},
				}})
			})
		}
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		for _, paramType := range []string{"header", "path", "query", "formData"} {
			t.Run(paramType, func(t *testing.T) {
				assert.Error(t, NewOperation(nil).ParseComment(`@Param some_object `+paramType+` main.Object true "Some Object"`, nil))
			})
		}
	})

}

// Test ParseParamComment Query Params
func TestParseParamCommentBodyArray(t *testing.T) {
	t.Parallel()

	comment := `@Param names body []string true "Users List"`
	o := NewOperation(nil)
	err := o.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Equal(t, o.Parameters, []spec.Parameter{{
		ParamProps: spec.ParamProps{
			Name:        "names",
			Description: "Users List",
			In:          "body",
			Required:    true,
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Type: []string{"string"},
							},
						},
					},
				},
			},
		},
	}})
}

// Test ParseParamComment Params
func TestParseParamCommentArray(t *testing.T) {
	paramTypes := []string{"header", "path", "query"}

	for _, paramType := range paramTypes {
		t.Run(paramType, func(t *testing.T) {
			operation := NewOperation(nil)
			err := operation.ParseComment(`@Param names `+paramType+` []string true "Users List"`, nil)

			assert.NoError(t, err)

			b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
			expected := `[
    {
        "type": "array",
        "items": {
            "type": "string"
        },
        "description": "Users List",
        "name": "names",
        "in": "` + paramType + `",
        "required": true
    }
]`
			assert.Equal(t, expected, string(b))

			err = operation.ParseComment(`@Param names `+paramType+` []model.User true "Users List"`, nil)
			assert.Error(t, err)
		})
	}
}

// Test TestParseParamCommentDefaultValue Query Params
func TestParseParamCommentDefaultValue(t *testing.T) {
	t.Parallel()

	operation := NewOperation(nil)
	err := operation.ParseComment(`@Param names query string true "Users List" default(test)`, nil)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "type": "string",
        "default": "test",
        "description": "Users List",
        "name": "names",
        "in": "query",
        "required": true
    }
]`
	assert.Equal(t, expected, string(b))
}

// Test ParseParamComment Query Params
func TestParseParamCommentQueryArrayFormat(t *testing.T) {
	t.Parallel()

	comment := `@Param names query []string true "Users List" collectionFormat(multi)`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "type": "array",
        "items": {
            "type": "string"
        },
        "collectionFormat": "multi",
        "description": "Users List",
        "name": "names",
        "in": "query",
        "required": true
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByID(t *testing.T) {
	t.Parallel()

	comment := `@Param unsafe_id[lte] query int true "Unsafe query param"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "type": "integer",
        "description": "Unsafe query param",
        "name": "unsafe_id[lte]",
        "in": "query",
        "required": true
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByQueryType(t *testing.T) {
	t.Parallel()

	comment := `@Param some_id query int true "Some ID"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "type": "integer",
        "description": "Some ID",
        "name": "some_id",
        "in": "query",
        "required": true
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyType(t *testing.T) {
	t.Parallel()

	comment := `@Param some_id body model.OrderRow true "Some ID"`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.OrderRow")
	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "description": "Some ID",
        "name": "some_id",
        "in": "body",
        "required": true,
        "schema": {
            "$ref": "#/definitions/model.OrderRow"
        }
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyTextPlain(t *testing.T) {
	t.Parallel()

	comment := `@Param text body string true "Text to process"`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "description": "Text to process",
        "name": "text",
        "in": "body",
        "required": true,
        "schema": {
            "type": "string"
        }
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyTypeWithDeepNestedFields(t *testing.T) {
	t.Parallel()

	comment := `@Param body body model.CommonHeader{data=string,data2=int} true "test deep"`
	operation := NewOperation(nil)

	operation.parser.addTestType("model.CommonHeader")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.Parameters, 1)
	assert.Equal(t, "test deep", operation.Parameters[0].Description)
	assert.True(t, operation.Parameters[0].Required)

	b, err := json.MarshalIndent(operation.Parameters, "", "    ")
	assert.NoError(t, err)
	expected := `[
    {
        "description": "test deep",
        "name": "body",
        "in": "body",
        "required": true,
        "schema": {
            "allOf": [
                {
                    "$ref": "#/definitions/model.CommonHeader"
                },
                {
                    "type": "object",
                    "properties": {
                        "data": {
                            "type": "string"
                        },
                        "data2": {
                            "type": "integer"
                        }
                    }
                }
            ]
        }
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyTypeArrayOfPrimitiveGo(t *testing.T) {
	t.Parallel()

	comment := `@Param some_id body []int true "Some ID"`
	operation := NewOperation(nil)
	err := operation.ParseComment(comment, nil)

	assert.NoError(t, err)
	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "description": "Some ID",
        "name": "some_id",
        "in": "body",
        "required": true,
        "schema": {
            "type": "array",
            "items": {
                "type": "integer"
            }
        }
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyTypeArrayOfPrimitiveGoWithDeepNestedFields(t *testing.T) {
	t.Parallel()

	comment := `@Param body body []model.CommonHeader{data=string,data2=int} true "test deep"`
	operation := NewOperation(nil)
	operation.parser.addTestType("model.CommonHeader")

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)
	assert.Len(t, operation.Parameters, 1)
	assert.Equal(t, "test deep", operation.Parameters[0].Description)
	assert.True(t, operation.Parameters[0].Required)

	b, err := json.MarshalIndent(operation.Parameters, "", "    ")
	assert.NoError(t, err)
	expected := `[
    {
        "description": "test deep",
        "name": "body",
        "in": "body",
        "required": true,
        "schema": {
            "type": "array",
            "items": {
                "allOf": [
                    {
                        "$ref": "#/definitions/model.CommonHeader"
                    },
                    {
                        "type": "object",
                        "properties": {
                            "data": {
                                "type": "string"
                            },
                            "data2": {
                                "type": "integer"
                            }
                        }
                    }
                ]
            }
        }
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByBodyTypeErr(t *testing.T) {
	t.Parallel()

	comment := `@Param some_id body model.OrderRow true "Some ID"`
	operation := NewOperation(nil)
	operation.parser.addTestType("model.notexist")
	err := operation.ParseComment(comment, nil)

	assert.Error(t, err)
}

func TestParseParamCommentByFormDataType(t *testing.T) {
	t.Parallel()

	comment := `@Param file formData file true "this is a test file"`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    {
        "type": "file",
        "description": "this is a test file",
        "name": "file",
        "in": "formData",
        "required": true
    }
]`
	assert.Equal(t, expected, string(b))
}

func TestParseParamCommentByFormDataTypeUint64(t *testing.T) {
	t.Parallel()

	comment := `@Param file formData uint64 true "this is a test file"`
	operation := NewOperation(nil)

	err := operation.ParseComment(comment, nil)
	assert.NoError(t, err)

	b, _ := json.MarshalIndent(operation.Parameters, "", "    ")
	expected := `[
    