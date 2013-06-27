package root

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

type Config struct {
	Name          string            `yaml:"name"`
	Addons        []string          `yaml:"addons,flow,omitempty"`
	Collaborators []string          `yaml:"collaborators,flow,omitempty"`
	Domains       []string          `yaml:"domains,flow,omitempty"`
	Config        map[string]string `yaml:"config,flow,omitempty"`
}

func (cmd *Root) LoadConfig() error {
	if cmd.Config != nil {
		return nil
	}

	config, err := load_config("./.forklift.yaml")
	if err != nil {
		return err
	}

	cmd.Config = config
	return nil
}

func load_config(path string) (*Config, error) {
	var (
		config *Config
	)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = goyaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
