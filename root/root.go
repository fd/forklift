package root

import (
	"github.com/fd/go-cli/cli"
)

func init() {
	cli.Register(Root{})
}

type Root struct {
	cli.Root
	cli.Arg0

	Account           string `env:"HEROKU_EMAIL"`
	ApiKey            string `env:"HEROKU_API_KEY"`
	DryRun            bool   `flag:"--dry"`
	UpdateDeploypacks bool   `flag:"--update-deploypacks"`
	Target            string `flag:"-t" env:"TARGET"`

	cli.Manual `
    Usage:   forklift <cmd> <options>
    Summary: Manage heroku applications

    .Account:
      The email of the heroku account to use.

    .ApiKey:
      The API key for the current heroku account.

    .DryRun:
      Don't actually change anything.

    .UpdateDeploypacks:
      Update the deploypacks before using them.

    .Target:
      The name of the environment.
  `
}
