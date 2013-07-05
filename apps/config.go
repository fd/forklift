package apps

import (
	"bytes"
	"fmt"
	"os"
	"path"
)

type Config struct {
	App        App
	Env        Env
	Deploypack string
	Unused     map[string]interface{}
}

func Load(target string) (*Config, error) {

	cnf, err := DecodeTOML(r)
	if err != nil {
		return err
	}

	for cnf.Deploypack != "" {
		var (
			in  bytes.Buffer
			out bytes.Buffer
		)

		err := EncodeJSON(&in, cnf)
		if err != nil {
			return err
		}

		// run

		cnf, err = DecodeJSON(&out)
		if err != nil {
			return err
		}
	}

}

func lookup_config_filename(target string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Unable to find configuration file for %s (%s)", target, err)
	}

	dir := wd
	subpath := path.Join(".forklift", target+".toml")

	for dir != "/" {
		filename := path.Join(dir, subpath)

		fi, err := os.Stat(filename)
		if err == nil && fi.Mode().IsRegular() {
			return filename, nil
		}

		dir = path.Dir(dir)
	}

	return "", fmt.Errorf("Unable to find configuration file for %s in %s", target, wd)
}
