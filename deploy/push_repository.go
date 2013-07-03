package deploy

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
func (c *Deploy) push_repository() error {
	var (
		cmd *exec.Cmd
		src string
		err error
	)

	cmd = exec.Command("git", "remote", "update", "origin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil
	err = cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("")

	if c.Config.Upstream != "" {
		cmd = exec.Command("git", "tag", "-l", "deploy-"+c.Config.Upstream+"-*")
		data, err := cmd.Output()
		if err != nil {
			return err
		}

		tags := strings.Split(strings.TrimSpace(string(data)), "\n")
		sort.Strings(tags)
		if len(tags) == 0 {
			return fmt.Errorf("No upstream deploy found for target %s", c.Config.Upstream)
		}
		tag := tags[len(tags)-1]
		if tag == "" {
			return fmt.Errorf("No upstream deploy found for target %s", c.Config.Upstream)
		}
		src = tag
	}

	if src == "" {
		fmt.Printf("Pushing %s:\n", c.Target)
	} else {
		fmt.Printf("Pushing %s => %s:\n", src, c.Target)
	}
	if c.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		if src == "" {
			src = "origin/master"
		}
		cmd = exec.Command("git", "push", "git@heroku.com:"+c.Config.Name+".git", src+":master", "--force")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
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
func (c *Deploy) tag_repository() error {
	var (
		now time.Time
		tag string
		msg string
		cmd *exec.Cmd
		err error
	)

	now = time.Now().UTC()
	tag = fmt.Sprintf("deploy-%s-%s", c.Target, now.Format("20060102150405"))
	msg = fmt.Sprintf("Deploy to %s at %s by %s", c.Target, now.Format(time.RFC3339), c.Account)

	fmt.Printf("Tagging commit as %s\n", tag)
	if c.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		cmd = exec.Command("git", "tag", "-a", tag, "-m", msg, "origin/master")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = nil
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	fmt.Println("")

	fmt.Printf("Pushing tag to origin:\n")
	if c.DryRun {
		fmt.Printf(" - skipped (dry run)\n")
	} else {
		cmd = exec.Command("git", "push", "origin", tag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = nil
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}