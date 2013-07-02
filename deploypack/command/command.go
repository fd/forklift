package command

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bitbucket.org/mrhenry/forklift/deploypack/runner"
)

type Handler interface {
	ProcessConfiguration(ctx *Context) error
}

func Run(h Handler) {
	ctx := &Context{
		stdout: os.Stdout,
	}

	os.Stdout = os.Stderr

	err := h.ProcessConfiguration(ctx)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
		return
	}
}

type Context struct {
	stdout io.Writer
}

func (c *Context) LoadConfig(cnf interface{}) error {
	return json.NewDecoder(os.Stdin).Decode(cnf)
}

func (c *Context) DumpConfig(cnf interface{}) error {
	return json.NewEncoder(c.stdout).Encode(cnf)
}

func (c *Context) RunDeploypack(ref string, in, out interface{}) error {
	return runner.Run(ref, in, out)
}

func (c *Context) Target() string {
	return os.Getenv("FORKLIFT_TARGET")
}

func (c *Context) Dir() string {
	return os.Getenv("FORKLIFT_DIR")
}
