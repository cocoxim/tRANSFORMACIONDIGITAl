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
			if len(genericParam) <= 2 || genericParam[:2] != "[]" {
				break
			}
			genericParam = genericParam[2:]
			arrayDepth++
		}

		typeDef := pkgDefs.FindTypeSpec(genericParam, file)
		if typeDef != nil {
			genericParam = typeDef.TypeName()
			if _, ok := pkgDefs.uniqueDefinitions[genericParam]; !ok {
				pkgDefs.uniqueDefinitions[genericParam] = typeDef
			}
		}

		genericParamTypeDefs[original.TypeSpec.TypeParams.List[i].Names[0].Name] = &genericTypeSpec{
			ArrayDepth: arrayDepth,
			TypeSpec:   typeDef,
			Name:       genericParam,
		}
	}

	name = fmt.Sprintf("%s%s-", string(IgnoreNameOverridePrefix), original.TypeName())
	var nameParts []string
	for _, def := range original.TypeSpec.TypeParams.List {
		if specDef, ok := genericParamTypeDefs[def.Names[0].Name]; ok {
			var prefix = ""
			if specDef.ArrayDepth == 1 {
				prefix = "array_"
			} else if specDef.ArrayDepth > 1 {
				prefix = fmt.Sprintf("array%d_", specDef.ArrayDepth)
			}
			nameParts = append(nameParts, prefix+specDef.TypeName())
		}
	}

	name += strings.Replace(strings.Join(nameParts, "-"), ".", "_", -1)

	if typeSpec, ok := pkgDefs.uniqueDefinitions[name]; ok {
		return typeSpec
	}

	parametrizedTypeSpec := &TypeSpecDef{
		File:    original.File,
		PkgPath: original.PkgPath,
		TypeSpec: &ast.TypeSpec{
			Name: &ast.Ident{
				Name:    name,
				NamePos: original.TypeSpec.Name.NamePos,
				Obj:     original.TypeSpec.Name.Obj,
			},
			Doc:    original.TypeSpec.Doc,
			Assign: original.TypeSpec.Assign,
		},
	}
	pkgDefs.uniqueDefinitions[name] = parametrizedTypeSpec

	parametrizedTypeSpec.TypeSpec.Type = pkgDefs.resolveGenericType(original.File, original.TypeSpec.Type, genericParamTypeDefs)

	return parametrizedTypeSpec
}

// splitGenericsTypeName splits a generic struct name in his parts
func splitGenericsTypeName(fullGenericForm string) (string, []string) {
	//remove all spaces character
	f