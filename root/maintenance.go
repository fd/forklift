package root

import (
	"fmt"
)

func (cmd *Root) SetMaintenance(flag bool) error {
	err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	data := struct {
		Maintenance bool `json:"maintenance"`
	}{flag}

	err = cmd.Http("PATCH", data, nil, "/apps/%s", cmd.Config.Name)
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
