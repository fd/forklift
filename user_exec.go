package main

import (
	"os"
	"os/user"
	"path/filepath"
	"syscall"
)

func user_exec() {
	home, err := get_home_dir()
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

var home_dir string

func get_home_dir() (string, error) {
	if home_dir == "" {
		u, err := user.Current()
		if err != nil {
			return "", err
		}

		home_dir = u.HomeDir
	}
	return home_dir, nil
}
