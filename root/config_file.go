package root

import (
	"os"

	"github.com/fd/forklift/apps"
)

func (cmd *Root) LoadTarget() (*apps.Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	cnf, err := apps.Load(wd, cmd.Target, cmd.DryRun, cmd.UpdateDeploypacks)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}
