package apps

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// - git remote update origin
// - git push {heroku remote} {tip}:master --force
func (app *App) push_repository() error {
	var (
		cmd *exec.Cmd
		src string
		err error
	)

	cmd = exec.Command("git", "remote", "update", "origin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	cmd.Stdin = nil
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("")

	if app.Upstream != "" {
		cmd = exec.Command("git", "tag", "-l", "deploy-"+app.Upstream+"-*")
		data, err := cmd.Output()
		if err != nil {
			return err
		}

		tags := strings.Split(strings.TrimSpace(string(data)), "\n")
		sort.Strings(tags)
		if len(tags) == 0 {
			return fmt.Errorf("No upstream deploy found for target %s", app.Upstream)
		}
		tag := tags[len(tags)-1]
		if tag == "" {
			return fmt.Errorf("No upstream deploy found for target %s", app.Upstream)
		}
		src = tag
	}

	if src == "" {
		fmt.Printf("Pushing %s:\n", app.config.Target)
	} else {
		fmt.Printf("Pushing %s => %s:\n", src, app.config.Target)
	}
	if app.config.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		if src == "" {
			src = "origin/master"
		}

		cmd = exec.Command("git", "rev-parse", src+"^{commit}")
		sha_data, err := cmd.Output()
		if err != nil {
			return err
		}

		sha := strings.TrimSpace(string(sha_data))

		cmd = exec.Command("git", "push", "git@heroku.com:"+app.AppName+".git", sha+":refs/heads/master", "--force")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		cmd.Stdin = nil
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

// - git tag -a "deploy-{target}-{timestamp}" -F - {tip}
// - git push origin "deploy-{target}-{timestamp}"
func (app *App) tag_repository() error {
	var (
		now time.Time
		tag string
		msg string
		cmd *exec.Cmd
		err error
	)

	now = time.Now().UTC()
	tag = fmt.Sprintf("deploy-%s-%s", app.config.Target, now.Format("20060102150405"))
	msg = fmt.Sprintf("Deploy to %s at %s by %s", app.config.Target, now.Format(time.RFC3339), app.config.Env.CurrentUser.Email)

	fmt.Printf("Tagging commit as %s\n", tag)
	if app.config.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		cmd = exec.Command("git", "tag", "-a", tag, "-m", msg, "origin/master")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		cmd.Stdin = nil
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	fmt.Println("")

	fmt.Printf("Pushing tag to origin:\n")
	if app.config.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		cmd = exec.Command("git", "push", "origin", tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stdout
		cmd.Stdin = nil
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
