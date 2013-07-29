package apps

import (
	"fmt"
)

func (app *App) Deploy() error {
	app.Pause()
	defer app.Unpause()

	var (
		err error
	)

	err = app.sync_collaborators()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.sync_features()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.sync_domains()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.sync_addons()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.sync_config()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.push_repository()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.run_post_push_commands()
	if err != nil {
		return err
	}
	fmt.Println("")

	err = app.tag_repository()
	if err != nil {
		return err
	}

	return nil
}
