package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fd/forklift/update"
	"github.com/fd/forklift/util/user"
	"github.com/fd/go-cli/cli"
)

const default_auto_update_interval = 24 * time.Hour

func auto_update() {
	if should_update() {
		updated := do_update()
		mark_last_check()
		if updated {
			user_exec()
		}
	}
}

func do_update() bool {
	home, err := user.Home()
	if err != nil {
		return false
	}

	path := filepath.Join(home, ".forklift", "bin", "forklift")
	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return false
	}

	cmd := update.Update{}
	cmd.Root.Arg0 = cli.Arg0(path)
	cmd.Main()
	return true
}

func should_update() bool {
	var (
		interval = auto_update_interval()
	)

	if interval < 0 {
		return false
	}

	if interval == 0 {
		return true
	}

	home, err := user.Home()
	if err != nil {
		return false
	}

	path := filepath.Join(home, ".forklift", "conf", "autoupdate.last")

	_, err = os.Stat(path)
	if err != nil {
		return true
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return true
	}

	last_check, err := time.Parse("2006-01-02", strings.TrimSpace(string(data)))
	if err != nil {
		return true
	}

	next_check := last_check.Add(interval)

	return next_check.Before(time.Now())
}

func mark_last_check() {
	home, err := user.Home()
	if err != nil {
		return
	}

	path := filepath.Join(home, ".forklift", "conf", "autoupdate.last")

	err = os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return
	}

	data := []byte(time.Now().Format("2006-01-02") + "\n")
	ioutil.WriteFile(path, data, 0600)
}

func auto_update_interval() time.Duration {
	home, err := user.Home()
	if err != nil {
		return default_auto_update_interval
	}

	path := filepath.Join(home, ".forklift", "conf", "autoupdate")

	_, err = os.Stat(path)
	if err != nil {
		return default_auto_update_interval
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return default_auto_update_interval
	}

	text := strings.TrimSpace(string(data))

	if text == "never" {
		return -1
	}

	if text == "always" {
		return 0
	}

	i, err := strconv.Atoi(text)
	if err != nil {
		return default_auto_update_interval
	}

	d := time.Duration(i) * 24 * time.Hour
	if d < 0 {
		d = default_auto_update_interval
	}

	return d
}
