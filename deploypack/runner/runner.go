package runner

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"

	"github.com/fd/forklift/deploypack/helpers"
	"github.com/fd/forklift/deploypack/installer"
)

func Run(ref string, in, out interface{}) error {
	var (
		data       []byte
		config_map map[string]interface{}
		err        error
	)

	data, err = json.Marshal(in)
	if err != nil {
		return err
	}

	for ref != "" {
		data, err = run(ref, in)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &config_map)
		if err != nil {
			return err
		}

		ref, err = helpers.ExtractDeploypack(config_map)
		if err != nil {
			return err
		}

		in = config_map
	}

	return json.Unmarshal(data, out)
}

func run(ref string, in interface{}) ([]byte, error) {
	err := installer.Run(ref)
	if err != nil {
		return nil, err
	}

	dir, err := helpers.Path(ref)
	if err != nil {
		return nil, err
	}

	bin, err := lookup_bin(ref, dir)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin)

	w, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	r, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(w).Encode(in)
	if err != nil {
		return nil, err
	}

	w.Close()

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, r)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
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
