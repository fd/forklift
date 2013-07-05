package apps

import (
	"fmt"
)

func (app *App) SetMaintenance(flag bool) error {
	data := struct {
		Maintenance bool `json:"maintenance"`
	}{flag}

	err := app.HttpV3("PATCH", data, nil, "/apps/%s", app.AppName)
	if err != nil {
		if flag {
			fmt.Printf("Maintenance: still off (error: %s)\n", err)
		} else {
			fmt.Printf("Maintenance: still on (error: %s)\n", err)
		}
		return err
	}

	if flag {
		fmt.Printf("Maintenance: on\n")
	} else {
		fmt.Printf("Maintenance: off\n")
	}

	return nil
}
