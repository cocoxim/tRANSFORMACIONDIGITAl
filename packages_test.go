package swag

import (
	"go/ast"
	"go/token"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackagesDefinitions_ParseFile(t *testing.T) {
	pd := PackagesDefinitions{}
	packageDir := "github.com/swaggo/swag/testdata/simple"
	assert.NoError(t, pd.ParseFile(packageDir, "testdata/simple/main.go", nil, ParseAll))
	assert.Equal(t, 1, len(pd.packages))
	assert.Equal(t, 1, len(pd.files))
}

func TestPackagesDefinitions_collectAstFile(t *testing.T) {
	pd := PackagesDefinitions{}
	fileSet := token.NewFileSet()
	assert.NoError(t, pd.collectAstFile(fileSet, "", "", nil, ParseAll))

	firstFile := &ast.File{
		Name: &ast.Ident{Name: "main.go"},
	}

	packageDir := "github.com/swaggo/swag/testdata/simple"
	assert.NoError(t, pd.collectAstFile(fileSet, packageDir, "testdata/simple/"+firstFile.Name.String(), firstFile, ParseAll))
	assert.NotEmpty(t, pd.packages[packageDir])

	absPath, _ := filepath.Abs("testdata/simple/" + firstFile.Name.String())
	astFileInfo := &AstFileInfo{
		FileSet:     fileSet,
		File:        firstFile,
		Path:        absPath,
		PackagePath: packageDir,
		ParseFlag:   ParseAll,
	}
	assert.Equal(t, pd.files[firstFile], astFileInfo)

	// Override
	assert.NoError(t, pd.collectAstFile(fileSet, packageDir, "testdata/simple/"+firstFile.Name.String(), firstFile, ParseAll))
	assert.Equal(t, pd.files[firstFile], astFileInfo)

	// Another file
	secondFile := &ast.File{
		Name: &ast.Ident{Name: "api.go"},
	}
	assert.NoError(t, pd.collectAstFile(fileSet, packageDir, "testdata/simple/"+secondFile.Name.String(), secondFile, ParseAll))
}

func TestPackagesDefinitions_rangeFiles(t *testing.T) {
	pd := PackagesDefinitions{
		files: map[*ast.File]*AstFileInfo{
			{
				Name: &ast.Ident{Name: "main.go"},
			}: {
				File:        &ast.File{Name: &ast.Ident{Name: "main.go"}},
				Path:        "testdata/simple/main.go",
				PackagePath: "main",
			},
			{
				Name: &ast.Ident{Name: "api.go"},
			}: {
				File:        &ast.File{Name: &ast.Ident{Name: "api.go"}},
				Path:        "testdata/simple/api/api.go",
				PackagePath: "api",
			},
		},
	}

	i, expect := 0, []string{"testdata/simple/api/api.go", "testdata/simple/main.go"}
	_ = pd.RangeFiles(func(fileInfo *AstFileInfo) error {
		assert.Equal(t, expect[i], fileInfo.Path)
		i++
		return nil
	})
}

func TestPackagesDefinitions_ParseTypes(t *testing.T) {
	absPath, _ := filepath.Abs("")

	mainAST := ast.File{
		Name: &ast.Ident{Name: "main.go"},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{Name: "Test"},
						Type: &ast.Ident{
							Name: "string",
						},
					},
				},
			},
		},
	}

	pd := PackagesDefinitions{
		files: map[*ast.File]*AstFileInfo{
			&mainAST: {
				File:        &mainAST,
				Path:        filepath.Join(absPath, "testdata/simple/main.go"),
				PackagePath: "main",
			},
			{
				Name: &ast.Ident{Name: "api.go"},
			}: {
				File:        &ast.File{Name: &ast.Ident{Name: "api.go"}},
				Path:        filepath.Join(absPath, "testdata/simple/api/api.go"),
				PackagePath: "api",
			},
		},
		packages: make(map[string]*PackageDefinitions),
	}

	_, err := pd.ParseTypes()
	assert.NoError(t, err)
}

func TestPackagesDefinitions_parseFunctionScopedTypesFromFile(t *testing.T) {
	mainAST := &ast.File{
		Name: &ast.Ident{Name: "main.go"},
		Decls: []ast.Decl{
			&ast.FuncDecl{
				Name: ast.NewIdent("TestFuncDecl"),
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.TYPE,
								Specs: []ast.Spec{
									&ast.TypeSpec{
										Name: ast.NewIdent("response"),
										Type: ast.NewIdent("struct"),
									},
									&ast.TypeSpec{
										Name: ast.NewIdent("stringResponse"),
										Type: ast.NewIdent("string"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	pd := PackagesDefinitions{
		packages: make(map[string]*PackageDefinitions),
	}

	parsedSchema := 