package main

import (
	"github.com/fd/go-cli/cli"
)

import (
	_ "github.com/fd/forklift/root"
	_ "github.com/fd/forklift/setup"
)

func main() { cli.Main(nil, nil) }
