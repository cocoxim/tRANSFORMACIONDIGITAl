
package swag

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/tools/go/loader"
)

// PackagesDefinitions map[package import path]*PackageDefinitions.
type PackagesDefinitions struct {
	files             map[*ast.File]*AstFileInfo
	packages          map[string]*PackageDefinitions
	uniqueDefinitions map[string]*TypeSpecDef
	parseDependency   bool
	debug             Debugger
}

// NewPackagesDefinitions create object PackagesDefinitions.
func NewPackagesDefinitions() *PackagesDefinitions {
	return &PackagesDefinitions{
		files:             make(map[*ast.File]*AstFileInfo),
		packages:          make(map[string]*PackageDefinitions),
		uniqueDefinitions: make(map[string]*TypeSpecDef),
	}
}

// ParseFile parse a source file.
func (pkgDefs *PackagesDefinitions) ParseFile(packageDir, path string, src interface{}, flag ParseFlag) error {
	// positions are relative to FileSet
	fileSet := token.NewFileSet()
	astFile, err := goparser.ParseFile(fileSet, path, src, goparser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file %s, error:%+v", path, err)
	}
	return pkgDefs.collectAstFile(fileSet, packageDir, path, astFile, flag)
}

// collectAstFile collect ast.file.
func (pkgDefs *PackagesDefinitions) collectAstFile(fileSet *token.FileSet, packageDir, path string, astFile *ast.File, flag ParseFlag) error {
	if pkgDefs.files == nil {
		pkgDefs.files = make(map[*ast.File]*AstFileInfo)
	}

	if pkgDefs.packages == nil {
		pkgDefs.packages = make(map[string]*PackageDefinitions)
	}

	// return without storing the file if we lack a packageDir
	if packageDir == "" {
		return nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dependency, ok := pkgDefs.packages[packageDir]
	if ok {
		// return without storing the file if it already exists
		_, exists := dependency.Files[path]
		if exists {
			return nil
		}

		dependency.Files[path] = astFile
	} else {
		pkgDefs.packages[packageDir] = NewPackageDefinitions(astFile.Name.Name, packageDir).AddFile(path, astFile)
	}

	pkgDefs.files[astFile] = &AstFileInfo{
		FileSet:     fileSet,
		File:        astFile,
		Path:        path,
		PackagePath: packageDir,
		ParseFlag:   flag,
	}

	return nil
}

// RangeFiles for range the collection of ast.File in alphabetic order.
func (pkgDefs *PackagesDefinitions) RangeFiles(handle func(info *AstFileInfo) error) error {
	sortedFiles := make([]*AstFileInfo, 0, len(pkgDefs.files))
	for _, info := range pkgDefs.files {
		// ignore package path prefix with 'vendor' or $GOROOT,
		// because the router info of api will not be included these files.
		if strings.HasPrefix(info.PackagePath, "vendor") || strings.HasPrefix(info.Path, runtime.GOROOT()) {
			continue
		}
		sortedFiles = append(sortedFiles, info)
	}

	sort.Slice(sortedFiles, func(i, j int) bool {
		return strings.Compare(sortedFiles[i].Path, sortedFiles[j].Path) < 0
	})

	for _, info := range sortedFiles {
		err := handle(info)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseTypes parse types
// @Return parsed definitions.
func (pkgDefs *PackagesDefinitions) ParseTypes() (map[*TypeSpecDef]*Schema, error) {
	parsedSchemas := make(map[*TypeSpecDef]*Schema)
	for astFile, info := range pkgDefs.files {
		pkgDefs.parseTypesFromFile(astFile, info.PackagePath, parsedSchemas)
		pkgDefs.parseFunctionScopedTypesFromFile(astFile, info.PackagePath, parsedSchemas)
	}
	pkgDefs.removeAllNotUniqueTypes()
	pkgDefs.evaluateAllConstVariables()
	pkgDefs.collectConstEnums(parsedSchemas)
	return parsedSchemas, nil
}

func (pkgDefs *PackagesDefinitions) parseTypesFromFile(astFile *ast.File, packagePath string, parsedSchemas map[*TypeSpecDef]*Schema) {
	for _, astDeclaration := range astFile.Decls {
		generalDeclaration, ok := astDeclaration.(*ast.GenDecl)
		if !ok {
			continue
		}
		if generalDeclaration.Tok == token.TYPE {
			for _, astSpec := range generalDeclaration.Specs {
				if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
					typeSpecDef := &TypeSpecDef{
						PkgPath:  packagePath,
						File:     astFile,
						TypeSpec: typeSpec,
					}

					if idt, ok := typeSpec.Type.(*ast.Ident); ok && IsGolangPrimitiveType(idt.Name) && parsedSchemas != nil {
						parsedSchemas[typeSpecDef] = &Schema{
							PkgPath: typeSpecDef.PkgPath,
							Name:    astFile.Name.Name,
							Schema:  PrimitiveSchema(TransToValidSchemeType(idt.Name)),
						}
					}

					if pkgDefs.uniqueDefinitions == nil {
						pkgDefs.uniqueDefinitions = make(map[string]*TypeSpecDef)
					}

					fullName := typeSpecDef.TypeName()

					anotherTypeDef, ok := pkgDefs.uniqueDefinitions[fullName]
					if ok {
						if anotherTypeDef == nil {
							typeSpecDef.NotUnique = true
							fullName = typeSpecDef.TypeName()
							pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
						} else if typeSpecDef.PkgPath != anotherTypeDef.PkgPath {
							pkgDefs.uniqueDefinitions[fullName] = nil
							anotherTypeDef.NotUnique = true
							pkgDefs.uniqueDefinitions[anotherTypeDef.TypeName()] = anotherTypeDef
							typeSpecDef.NotUnique = true
							fullName = typeSpecDef.TypeName()
							pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
						}
					} else {
						pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
					}

					if pkgDefs.packages[typeSpecDef.PkgPath] == nil {
						pkgDefs.packages[typeSpecDef.PkgPath] = NewPackageDefinitions(astFile.Name.Name, typeSpecDef.PkgPath).AddTypeSpec(typeSpecDef.Name(), typeSpecDef)
					} else if _, ok = pkgDefs.packages[typeSpecDef.PkgPath].TypeDefinitions[typeSpecDef.Name()]; !ok {
						pkgDefs.packages[typeSpecDef.PkgPath].AddTypeSpec(typeSpecDef.Name(), typeSpecDef)
					}
				}
			}
		} else if generalDeclaration.Tok == token.CONST {
			// collect consts
			pkgDefs.collectConstVariables(astFile, packagePath, generalDeclaration)
		}
	}
}

func (pkgDefs *PackagesDefinitions) parseFunctionScopedTypesFromFile(astFile *ast.File, packagePath string, parsedSchemas map[*TypeSpecDef]*Schema) {
	for _, astDeclaration := range astFile.Decls {
		funcDeclaration, ok := astDeclaration.(*ast.FuncDecl)
		if ok && funcDeclaration.Body != nil {
			for _, stmt := range funcDeclaration.Body.List {
				if declStmt, ok := (stmt).(*ast.DeclStmt); ok {
					if genDecl, ok := (declStmt.Decl).(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
						for _, astSpec := range genDecl.Specs {
							if typeSpec, ok := astSpec.(*ast.TypeSpec); ok {
								typeSpecDef := &TypeSpecDef{
									PkgPath:    packagePath,
									File:       astFile,
									TypeSpec:   typeSpec,
									ParentSpec: astDeclaration,
								}

								if idt, ok := typeSpec.Type.(*ast.Ident); ok && IsGolangPrimitiveType(idt.Name) && parsedSchemas != nil {
									parsedSchemas[typeSpecDef] = &Schema{
										PkgPath: typeSpecDef.PkgPath,
										Name:    astFile.Name.Name,
										Schema:  PrimitiveSchema(TransToValidSchemeType(idt.Name)),
									}
								}

								if pkgDefs.uniqueDefinitions == nil {
									pkgDefs.uniqueDefinitions = make(map[string]*TypeSpecDef)
								}

								fullName := typeSpecDef.TypeName()

								anotherTypeDef, ok := pkgDefs.uniqueDefinitions[fullName]
								if ok {
									if anotherTypeDef == nil {
										typeSpecDef.NotUnique = true
										fullName = typeSpecDef.TypeName()
										pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
									} else if typeSpecDef.PkgPath != anotherTypeDef.PkgPath {
										pkgDefs.uniqueDefinitions[fullName] = nil
										anotherTypeDef.NotUnique = true
										pkgDefs.uniqueDefinitions[anotherTypeDef.TypeName()] = anotherTypeDef
										typeSpecDef.NotUnique = true
										fullName = typeSpecDef.TypeName()
										pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
									}
								} else {
									pkgDefs.uniqueDefinitions[fullName] = typeSpecDef
								}

								if pkgDefs.packages[typeSpecDef.PkgPath] == nil {
									pkgDefs.packages[typeSpecDef.PkgPath] = NewPackageDefinitions(astFile.Name.Name, typeSpecDef.PkgPath).AddTypeSpec(fullName, typeSpecDef)
								} else if _, ok = pkgDefs.packages[typeSpecDef.PkgPath].TypeDefinitions[fullName]; !ok {
									pkgDefs.packages[typeSpecDef.PkgPath].AddTypeSpec(fullName, typeSpecDef)
								}
							}
						}

					}
				}
			}
		}
	}
}

func (pkgDefs *PackagesDefinitions) collectConstVariables(astFile *ast.File, packagePath string, generalDeclaration *ast.GenDecl) {
	pkg, ok := pkgDefs.packages[packagePath]
	if !ok {
		pkg = NewPackageDefinitions(astFile.Name.Name, packagePath)
		pkgDefs.packages[packagePath] = pkg
	}

	var lastValueSpec *ast.ValueSpec
	for _, astSpec := range generalDeclaration.Specs {
		valueSpec, ok := astSpec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if len(valueSpec.Names) == 1 && len(valueSpec.Values) == 1 {
			lastValueSpec = valueSpec
		} else if len(valueSpec.Names) == 1 && len(valueSpec.Values) == 0 && valueSpec.Type == nil && lastValueSpec != nil {
			valueSpec.Type = lastValueSpec.Type
			valueSpec.Values = lastValueSpec.Values
		}
		pkg.AddConst(astFile, valueSpec)