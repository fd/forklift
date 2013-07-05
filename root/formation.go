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

	if !cmd.DryRun {
		err := cmd.wait_until_dynos_are_down()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Root) RestoreFormation() error {
	fmt.Println("Restore formation:")
	tabw := tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
	count := 0

	for _, formation := range cmd.Formation {
		count += int(formation.Quantity)

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

	if !cmd.DryRun {
		err := cmd.wait_until_dynos_are_up(count)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Root) wait_until_dynos_are_down() error {
	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		ok, err := cmd.are_all_dynos_down()

		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		if ok {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("formation: unable to stop dynos")
}

func (cmd *Root) wait_until_dynos_are_up(count int) error {
	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		ok, err := cmd.are_all_dynos_up(count)

		if err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		if ok {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("formation: unable to start dynos")
}

func (cmd *Root) are_all_dynos_down() (bool, error) {
	var (
		data []struct {
			State string `json:"state"`
		}
	)

	err := cmd.Http("GET", nil, &data, "/apps/%s/dynos", cmd.Config.Name)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	return false, nil
}

func (cmd *Root) are_all_dynos_up(count int) (bool, error) {
	var (
		data []struct {
			AttachUrl *string `json:"attach_url"`
			State     string  `json:"state"`
		}
	)

	err := cmd.Http("GET", nil, &data, "/apps/%s/dynos", cmd.Config.Name)
	if err != nil {
		return false, err
	}

	for _, dyno := range data {
		if dyno.AttachUrl == nil {
			if dyno.State == "up" {
				count--
			}
		}
	}

	return count == 0, nil
}
