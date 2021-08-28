package swag

import (
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

// PackageDefinitions files and definition in a package.
type PackageDefinitions struct {
	// files in this package, map key is file's relative path starting package path
	Files map[string]*ast.File

	// definitions in this package, map key is typeName
	TypeDefinitions map[string]*TypeSpecDef

	// const variables in this package, map key is the name
	ConstTable map[string]*ConstVariable

	// const variables in order in this package
	OrderedConst []*ConstVariable

	// package name
	Name string

	// package path
	Path string
}

// ConstVariableGlobalEvaluator an interface used to evaluate enums across packages
type ConstVariableGlobalEvaluator interface {
	EvaluateConstValue(pkg *PackageDefinitions, cv *ConstVariable, recursiveStack map[string]struct{}) (interface{}, ast.Expr)
	EvaluateConstValueByName(file *ast.File, pkgPath, constVariableName string, recursiveStack map[string]struct{}) (interface{}, ast.Expr)
	FindTypeSpec(typeName string, file *ast.File) *TypeSpecDef
}

// NewPackageDefinitions new a PackageDefinitions object
func NewPackageDefinitions(name, pkgPath string) *PackageDefinitions {
	return &PackageDefinitions{
		Name:            name,
		Path:            pkgPath,
		Files:           make(map[string]*ast.File),
		TypeDefinitions: make(map[string]*TypeSpecDef),
		ConstTable:      make(map[string]*ConstVariable),
	}
}

// AddFile add a file
func (pkg *PackageDefinitions) AddFile(pkgPath string, file *ast.File) *PackageDefinitions {
	pkg.Files[pkgPath] = file
	return pkg
}

// AddTypeSpec add a type spec.
func (pkg *PackageDefinitions) AddTypeSpec(name string, typeSpec *TypeSpecDef) *PackageDefinitions {
	pkg.TypeDefinitions[name] = typeSpec
	return pkg
}

// AddConst add a const variable.
func (pkg *PackageDefinitions) AddConst(astFile *ast.File, valueSpec *