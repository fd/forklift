package apps

import (
	"fmt"
	"os"
	"os/exec"
)

func (app *App) run_post_push_commands() error {
	for _, command := range app.PostPushCommands {
		err := app.run_post_push_command(command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) run_post_push_command(command string) error {
	if app.config.DryRun {
		fmt.Printf("Run: %s\n - skipped (dry run)\n", command)
		return nil
	}

	cmd := exec.Command("heroku", "run", command)
	cmd.Env = append(
		os.Environ(),
		[]string{
			"HEROKU_APP=" + app.AppName,
		}...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "/"
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error=%s", err)
		return err
	}

	return nil
}
