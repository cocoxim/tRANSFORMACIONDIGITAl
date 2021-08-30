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
	assert.NoError(t, pd.ParseFile(packageDir, "testdata/