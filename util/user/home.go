package user

import (
	"fmt"
	"os"
	"os/user"
)

var home_dir string

func Home() (string, error) {
	if home_dir == "" {
		u, err := user.Current()
		if err == nil {
			home_dir = u.HomeDir
		}
	}
	if home_dir == "" {
		home_dir = os.Getenv("HOME")
	}
	if home_dir == "" {
		return "", fmt.Errorf("Unable to find home directory for %s", os.Getenv("USER"))
	}
	return home_dir, nil
}
