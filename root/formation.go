package root

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type Formation struct {
	Type     string `json:"type"`
	Size     uint8  `json:"size"`
	Quantity uint8  `json:"quantity"`
}

func (cmd *Root) BreakFormation() error {
	fmt.Println("Break formation:")
	tabw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)

	err := cmd.LoadConfig()
	if err != nil {
		return err
	}

	data := []Formation{}

	err = cmd.Http("GET", nil, &data, "/apps/%s/formation", cmd.Config.Name)
	if err != nil {
		return err
	}

	for _, formation := range data {
		if formation.Quantity == 0 {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\n",
				formation.Type, formation.Quantity, formation.Size)
			continue
		}

		formation.Quantity = 0

		err = cmd.Http("PATCH", &formation, nil, "/apps/%s/formation/%s", cmd.Config.Name, formation.Type)
		if err != nil {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\t\x1b[31;40;4;5m(failed to pause)\x1b[0m\n   error: %s\n",
				formation.Type, formation.Quantity, formation.Size, err)
		} else {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\t\x1b[32m(paused)\x1b[0m\n",
				formation.Type, formation.Quantity, formation.Size)
		}
	}

	tabw.Flush()

	cmd.Formation = data

	time.Sleep(12 * time.Second)

	return nil
}

func (cmd *Root) RestoreFormation() error {
	fmt.Println("Restore formation:")
	tabw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)

	for _, formation := range cmd.Formation {
		if formation.Quantity == 0 {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\n",
				formation.Type, formation.Quantity, formation.Size)
			continue
		}

		err := cmd.Http("PATCH", &formation, nil, "/apps/%s/formation/%s", cmd.Config.Name, formation.Type)
		if err != nil {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\t\x1b[31;40;4;5m(failed to restore)\x1b[0m\n   error: %s\n",
				formation.Type, formation.Quantity, formation.Size, err)
		} else {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%d\t\x1b[32m(restored)\x1b[0m\n",
				formation.Type, formation.Quantity, formation.Size)
		}
	}

	tabw.Flush()

	time.Sleep(12 * time.Second)

	return nil
}
