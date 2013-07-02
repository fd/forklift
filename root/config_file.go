package root

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	toml "github.com/pelletier/go-toml"

	"bitbucket.org/mrhenry/forklift/deploypack/helpers"
	"bitbucket.org/mrhenry/forklift/deploypack/runner"
)

type Config struct {
	Name          string
	Addons        []string
	Collaborators []string
	Domains       []string
	Environment   map[string]string
}

func (cmd *Root) LoadConfig() error {
	if cmd.Config != nil {
		return nil
	}

	filename, err := cmd.lookup_config_filename()
	if err != nil {
		return err
	}

	{
		dir := path.Dir(path.Dir(filename))
		os.Setenv("FORKLIFT_DIR", dir)
		os.Setenv("FORKLIFT_TARGET", cmd.Target)
		if cmd.DryRun {
			os.Setenv("FORKLIFT_DRYRUN", "true")
		}
		os.Chdir(dir)
	}

	config, err := load_config(filename)
	if err != nil {
		return err
	}

	cmd.Config = config
	return nil
}

func (cmd *Root) lookup_config_filename() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Unable to find configuration file for %s (%s)", cmd.Target, err)
	}

	dir := wd
	subpath := path.Join(".forklift", cmd.Target+".toml")

	for dir != "/" {
		filename := path.Join(dir, subpath)

		fi, err := os.Stat(filename)
		if err == nil && fi.Mode().IsRegular() {
			return filename, nil
		}

		dir = path.Dir(dir)
	}

	return "", fmt.Errorf("Unable to find configuration file for %s in %s", cmd.Target, wd)
}

func load_config(filename string) (*Config, error) {
	var (
		data       []byte
		tree       *toml.TomlTree
		config_map map[string]interface{}
		config     *Config
		deploypack string
		err        error
	)

	data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	{
		tree, err = toml.Load(string(data))
		if err != nil {
			return nil, err
		}

		config_map = map[string]interface{}(*tree)
		deploypack, err = helpers.ExtractDeploypack(config_map)
		if err != nil {
			return nil, err
		}
	}

	err = runner.Run(deploypack, config_map, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
