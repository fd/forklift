package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/fd/forklift/deploypack/helpers"
	"github.com/fd/forklift/deploypack/installer"
)

func Run(ref string, wd string, environ []string, in io.Reader, out io.Writer, update bool) error {
	err := installer.Run(ref, update)
	if err != nil {
		return err
	}

	dir, err := helpers.Path(ref)
	if err != nil {
		return err
	}

	bin, err := lookup_bin(ref, dir)
	if err != nil {
		return err
	}

	cmd := exec.Command(bin)

	cmd.Dir = wd

	cmd.Env = append(
		os.Environ(),
		environ...,
	)

	w, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	defer w.Close()

	r, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer r.Close()

	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	_, err = io.Copy(w, in)
	if err != nil {
		return err
	}

	w.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func lookup_bin(ref, dir string) (string, error) {
	var (
		bin string
	)

	bin = path.Join(dir, "bin", fmt.Sprintf("deploy-%s-%s", runtime.GOOS, runtime.GOARCH))
	if _, err := os.Stat(bin); err == nil {
		return bin, err
	}

	bin = path.Join(dir, "bin", "deploy")
	if _, err := os.Stat(bin); err == nil {
		return bin, err
	}

	return "", fmt.Errorf("no deploy command for %s (in %s)", ref, dir)
}
