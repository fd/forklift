package installer

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/fd/forklift/deploypack/helpers"
)

func Run(ref string, update bool) error {
	dir, err := helpers.Path(ref)
	if err != nil {
		return err
	}

	last_modified, err := deploypack_stat(dir)
	if err != nil {
		return err
	}

	// already up to date
	if update == false && os.Getenv("FORKLIFT_UPDATE_DEPLOYPACKS") != "true" {
		if last_modified.After(time.Now().AddDate(0, 0, -7)) {
			return nil
		}
	}

	if last_modified.IsZero() {
		return deploypack_install(ref, dir)
	} else {
		return deploypack_update(ref, dir)
	}
}

func deploypack_install(ref, dir string) error {
	fmt.Printf("Installing deploypack: %s\n", ref)

	err := os.MkdirAll(path.Dir(dir), 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", "--depth", "1", ref, dir)
	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Stdin = nil

	err = cmd.Run()
	if err != nil {
		return err
	}

	err = os.Chtimes(dir, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func deploypack_update(ref, dir string) error {
	fmt.Printf("Updating deploypack: %s\n", ref)

	cmd := exec.Command("git", "pull", ref, "master")
	cmd.Dir = dir
	cmd.Stderr = nil
	cmd.Stdout = nil
	cmd.Stdin = nil

	err := cmd.Run()
	if err != nil {
		return err
	}

	err = os.Chtimes(dir, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func deploypack_stat(dir string) (last_modified time.Time, err error) {
	var (
		fi os.FileInfo
	)

	fi, err = os.Stat(dir)

	if os.IsNotExist(err) {
		err = nil
		return
	}

	if err != nil {
		return
	}

	last_modified = fi.ModTime()
	return
}
