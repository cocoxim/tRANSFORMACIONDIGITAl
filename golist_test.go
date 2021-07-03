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