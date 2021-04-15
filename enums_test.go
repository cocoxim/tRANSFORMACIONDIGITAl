package swag

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGlobalEnums(t *testing.T) {
	searchDir := "testdata/enums"
	expected, err := os.ReadFile(filepath.Join(searchDir, "expected.json"))
	assert.NoError(t, err)

	p := New()
	err = p.ParseAPI(searchDir, mainAPIFile, defaultParseDepth)
	assert.NoError(t, err)
	b, err := json.MarshalIndent(p.swagger, "", "    ")
	assert.NoError(t, err)
	assert.Equal(t, string(expected), string(b))
	constsPath := "github.com/swaggo/swag/testdata/enums/consts"
	assert.Equal(t, 64, p.packages.packages[constsPath].ConstTable["uintSize"].Value)
	assert.Equal(t, int32(62), p.packages.packages[constsPath].ConstTable["maxBase"].Value)
	assert.Equal(t, 8, p.packages.packages[constsPath].ConstTable["shlByLen"].Value)
	as