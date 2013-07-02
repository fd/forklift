package main

import (
	"github.com/fd/go-cli/cli"
)

import (
	_ "bitbucket.org/mrhenry/forklift/deploy"
	_ "bitbucket.org/mrhenry/forklift/root"
)

func main() { cli.Main(nil, nil) }
