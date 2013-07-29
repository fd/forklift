package apps

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/fd/forklift/deploypack/runner"
)

type Config struct {
	Target            string
	RootDir           string
	DryRun            bool
	UpdateDeploypacks bool
	App               App
	Env               Env
	Deploypack        string
	Unused            map[string]interface{}
}

func Load(wd, target string, dryrun, update_buildpacks bool) (*Config, error) {
	root_dir, filename, err := lookup_config_filename(wd, target)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cnf, err := DecodeTOML(f)
	if err != nil {
		return nil, err
	}

	cnf.Target = target
	cnf.RootDir = root_dir
	cnf.DryRun = dryrun
	cnf.UpdateDeploypacks = update_buildpacks

	cnf, err = Expand(cnf)
	if err != nil {
		return nil, err
	}

	err = cnf.Env.load_heroku_credentials()
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

func Expand(cnf *Config) (*Config, error) {
	var (
		target            = cnf.Target
		root_dir          = cnf.RootDir
		dryrun            = cnf.DryRun
		update_buildpacks = cnf.UpdateDeploypacks
		err               error
	)

	for cnf.Deploypack != "" {
		var (
			in  bytes.Buffer
			out bytes.Buffer
		)

		err = EncodeJSON(&in, cnf)
		if err != nil {
			return nil, err
		}

		environ := []string{
			"FORKLIFT_DIR=" + cnf.RootDir,
			"FORKLIFT_TARGET=" + cnf.Target,
		}

		if cnf.DryRun {
			environ = append(
				environ,
				"FORKLIFT_DRYRUN=true",
			)
		}

		if cnf.UpdateDeploypacks {
			environ = append(
				environ,
				"FORKLIFT_UPDATE_DEPLOYPACKS=true",
			)
		}

		err = runner.Run(
			cnf.Deploypack,
			cnf.RootDir,
			environ,
			&in,
			&out,
			cnf.UpdateDeploypacks,
		)
		if err != nil {
			return nil, err
		}

		cnf, err = DecodeJSON(&out)
		if err != nil {
			return nil, err
		}

		cnf.Target = target
		cnf.RootDir = root_dir
		cnf.DryRun = dryrun
		cnf.UpdateDeploypacks = update_buildpacks
	}

	return cnf, nil
}

func lookup_config_filename(wd, target string) (string, string, error) {
	dir := wd
	subpath := path.Join(".forklift", target+".toml")

	for dir != "/" {
		filename := path.Join(dir, subpath)

		fi, err := os.Stat(filename)
		if err == nil && fi.Mode().IsRegular() {
			return dir, filename, nil
		}

		dir = path.Dir(dir)
	}

	return "", "", fmt.Errorf("Unable to find configuration file for %s in %s", target, wd)
}
