package command

import (
	"fmt"
	"os"

	"github.com/fd/forklift/apps"
)

type Handler interface {
	ProcessConfiguration(cnf *apps.Config) (*apps.Config, error)
}

func Run(h Handler) {
	var (
		stdout = os.Stdout
		cnf    *apps.Config
		err    error
	)

	os.Stdout = os.Stderr

	cnf, err = apps.DecodeJSON(os.Stdin)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
		return
	}

	cnf.Deploypack = ""
	cnf.RootDir = os.Getenv("FORKLIFT_DIR")
	cnf.Target = os.Getenv("FORKLIFT_TARGET")
	cnf.DryRun = os.Getenv("FORKLIFT_DRYRUN") == "true"
	cnf.UpdateDeploypacks = os.Getenv("FORKLIFT_UPDATE_DEPLOYPACKS") == "true"

	cnf, err = h.ProcessConfiguration(cnf)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
		return
	}

	err = apps.EncodeJSON(stdout, cnf)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
		return
	}
}
