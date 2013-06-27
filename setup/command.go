package setup

import (
	"fmt"

	"github.com/fd/forklift/root"
	"github.com/fd/go-cli/cli"
)

func init() {
	cli.Register(Setup{})
}

// - update collaborators
// - update domains
// - update syslog drains
// - update addons
// - update config
type Setup struct {
	root.Root
	cli.Arg0 `name:"setup"`

	cli.Manual `
    Usage:   forklift setup <options>
    Summary: Update application configurations
  `
}

func (cmd *Setup) Main() error {
	cmd.Pause()
	defer cmd.Unpause()

	var (
		err error
	)

	err = cmd.LoadConfig()
	if err != nil {
		return err
	}

	err = cmd.sync_collaborators()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = cmd.sync_domains()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = cmd.sync_config()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = cmd.sync_addons()
	if err != nil {
		return err
	}

	return nil
}
