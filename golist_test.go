package swag

import (
	"context"
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListPackages(t *testing.T) {

	cases := []struct {
		name      string
		args      []string
		searchDir string
		except    error
	}{
		{
			name:      "errorArgs",
			args:      []string{"-abc"},
			searchDir: "testdata/golist",
			except:    fmt.Errorf("exit status 2"),
		},
		{
			name:      "normal",
			args:      []string{"-deps"},
			searchDir: "testdata/golist",
			except:    nil,
		},
		{
			name:      "list error",
			args:      []string{"-deps"},
			searchDir: "testdata/golist_not_exist",
			except:    errors.New("searchDir not exist"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := listPackages(context.TODO(), c.searchDir, nil, c.args...)
			if c.except != nil {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetAllGoFileInfoFromDepsByList(t *testing.T) {
	p := New(ParseUsingGoList(true))
	pwd, err := os.Getwd()
	assert.NoError(t, err)
	cases := []struct {
		name           string
		buildPackage   *build.Package
		ignoreInternal bool
		except         error
