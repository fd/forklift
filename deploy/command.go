package deploy

import (
	"github.com/fd/forklift/root"
	"github.com/fd/forklift/static/gopkg.in/fd/go-cli.v1/cli"
)

func init() {
	cli.Register(Deploy{})
}

type Deploy struct {
	root.Root
	cli.Arg0 `name:"deploy"`

	cli.Manual `
    Usage:   forklift deploy <options>
    Summary: Update application configurations
  `
}

func (cmd *Deploy) Main() error {
	target, err := cmd.LoadTarget()
	if err != nil {
		return err
	}

	err = target.App.Deploy()
	if err != nil {
		return err
	}

	return nil
}
