package main

import (
	"os"
)

func in_dev_mode() bool {
	return os.Getenv("DEVMODE") == "true"
}
