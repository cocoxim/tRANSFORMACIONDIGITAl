//go:build go1.18
// +build go1.18

package swag

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"github.com/go-openapi/spec"
)

type genericTypeSpec struct {
	ArrayDepth int
	TypeSpec   *TypeSpecDef
	Name       string
}

func (t *genericTypeSpec) TypeName() string {
	if t.TypeSpec != nil {
		return t.TypeSpec.TypeName()
	}
	return t.Name
}

func (pkgDefs *PackagesDefinitions) parametrizeGenericType(file *ast.File, original *TypeSpecDef, fullGenericForm string) *TypeSpecDef {
	if original == nil || original.TypeSpec.TypeParams == nil || len(original.TypeSpec.TypeParams.List) == 0 {
		return original
	}

	name, genericParams := splitGenericsTypeName(fullGenericForm)
	if genericParams == nil {
		return nil
	}

	genericParamTypeDefs := map[string]*genericTypeSpec{}
	if len(genericParams) != len(original.TypeSpec.TypeParams.List) {
		return nil
	}

	for i, genericParam := range genericParams {
		arrayDepth := 0
		for {
			if len(genericParam) <= 2 || genericParam[:2