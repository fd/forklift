package root

import (
	"io/ioutil"

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

	config, err := load_config("./.forklift.toml")
	if err != nil {
		return err
	}

	cmd.Config = config
	return nil
}

func load_config(path string) (*Config, error) {
	var (
		data       []byte
		tree       *toml.TomlTree
		config_map map[string]interface{}
		config     *Config
		deploypack string
		err        error
	)

	data, err = ioutil.ReadFile(path)
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
