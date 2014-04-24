package main

import (
	"gopkg.cc/fd/go-cli/v1/cli"
)

import (
	_ "github.com/fd/forklift/deploy"
	_ "github.com/fd/forklift/root"
	_ "github.com/fd/forklift/update"
)

func main() {
	if !in_dev_mode() {
		user_exec()
		auto_update()
	}

	cli.Main(nil, nil)
}
