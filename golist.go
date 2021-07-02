package swag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"
)

func listPackages(ctx context.Context, dir string, env []string, args ...string) (pkgs []*build.Package, finalErr error) {
	cmd := exec.CommandContext(ctx, "go", append([]string{"list", "-json", "-e"}, args...)...)
	cmd.Env = env
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf
	defer func() {
		if stderrBuf.Len() > 0 {
			finalErr = fmt.Errorf("%v\n%s", finalErr, stderrBuf.Bytes())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(stdout)
	for dec.More() {
		var pkg build.Package
		err = dec.Decode(&pkg)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, &pkg)
	}
	err = cmd.Wait()
	if err != n