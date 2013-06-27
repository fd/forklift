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

	Account string `flag:"--account" env:"HEROKU_EMAIL"`
	ApiKey  string `env:"HEROKU_API_KEY"`

	Target string `flag:"-t" env:"TARGET"`

	cli.Manual `
    Usage:   forklift <cmd> <options>
    Summary: Manage heroku applications

    .Account:
      The email of the heroku account to use.

    .ApiKey:
      The API key for the current heroku account.

    .Target:
      The name of the environment.
  `

	Config    *Config
	Formation []Formation
}
