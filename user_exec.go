package main

import (
	"os"
	"path/filepath"
	"syscall"

	"github.com/fd/forklift/util/user"
)

func user_exec() {
	home, err := user.Home()
	if err != nil {
		return
	}

	path := filepath.Join(home, ".forklift", "bin", "forklift")

	_, err = os.Stat(path)
	if err != nil {
		return
	}

	if os.Args[0] == path {
		return
	}

	os.Args[0] = path
	syscall.Exec(path, os.Args, os.Environ())
}
