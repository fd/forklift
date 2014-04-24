package apps

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"
)

type Formation struct {
	Type     string `json:"type"`
	Size     string `json:"size"`
	Quantity uint8  `json:"quantity"`
}

func (app *App) FormationLoad() error {
	var (
		data      []Formation
		formation map[string]Formation
	)

	err := app.HttpV3("GET", nil, &data, "/apps/%s/formation", app.AppName)
	if err != nil {
		return err
	}

	formation = make(map[string]Formation, len(data))
	for _, typ := range data {
		formation[typ.Type] = typ
	}

	app.formation = formation
	return nil
}

func (app *App) FormationBreak() error {
	var (
		tabw = tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
	)

	fmt.Println("Break formation:")

	err := app.FormationLoad()
	if err != nil {
		return err
	}

	// sort keys
	keys := make([]string, 0, len(app.formation))
	for key := range app.formation {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// scale formation
	for _, key := range keys {
		formation := app.formation[key]

		if formation.Quantity == 0 {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\n",
				formation.Type, formation.Quantity, formation.Size)
			continue
		}

		formation.Quantity = 0

		err = app.HttpV3("PATCH", &formation, nil, "/apps/%s/formation/%s", app.AppName, formation.Type)
		if err != nil {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\t\x1b[31;40;4;5m(failed to pause)\x1b[0m\n   error: %s\n",
				formation.Type, formation.Quantity, formation.Size, err)
		} else {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\t\x1b[32m(paused)\x1b[0m\n",
				formation.Type, formation.Quantity, formation.Size)
		}
	}

	tabw.Flush()

	if err = app.wait_until_dynos_are_down(); err != nil {
		return err
	}

	return nil
}

func (app *App) FormationRestore() error {
	var (
		tabw  = tabwriter.NewWriter(os.Stdout, 8, 8, 1, ' ', 0)
		count = 0
	)

	fmt.Println("Restore formation:")

	prev_formation := app.formation
	err := app.FormationLoad()
	if err != nil {
		return err
	}
	for key := range app.formation {
		prev, ok := prev_formation[key]
		if !ok {
			continue
		}

		formation := app.formation[key]
		formation.Quantity = prev.Quantity
		app.formation[key] = formation
	}

	// sort keys
	keys := make([]string, 0, len(app.formation))
	for key := range app.formation {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		formation := app.formation[key]
		count += int(formation.Quantity)

		if formation.Quantity == 0 {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\n",
				formation.Type, formation.Quantity, formation.Size)
			continue
		}

		err := app.HttpV3("PATCH", &formation, nil, "/apps/%s/formation/%s", app.AppName, formation.Type)
		if err != nil {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\t\x1b[31;40;4;5m(failed to restore)\x1b[0m\n   error: %s\n",
				formation.Type, formation.Quantity, formation.Size, err)
		} else {
			fmt.Fprintf(tabw, " - %s\tquantity=%d\tsize=%s\t\x1b[32m(restored)\x1b[0m\n",
				formation.Type, formation.Quantity, formation.Size)
		}
	}

	tabw.Flush()

	if err := app.wait_until_dynos_are_up(count); err != nil {
		return err
	}

	return nil
}

func (app *App) wait_until_dynos_are_down() error {
	if app.config.DryRun {
		return nil
	}

	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		ok, err := app.are_all_dynos_down()

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

func (app *App) wait_until_dynos_are_up(count int) error {
	if app.config.DryRun {
		return nil
	}

	deadline := time.Now().Add(5 * time.Minute)

	for time.Now().Before(deadline) {
		ok, err := app.are_all_dynos_up(count)

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

func (app *App) are_all_dynos_down() (bool, error) {
	var (
		data []struct {
			State string `json:"state"`
		}
	)

	err := app.HttpV3("GET", nil, &data, "/apps/%s/dynos", app.AppName)
	if err != nil {
		return false, err
	}

	if len(data) == 0 {
		return true, nil
	}

	return false, nil
}

func (app *App) are_all_dynos_up(count int) (bool, error) {
	var (
		data []struct {
			AttachUrl *string `json:"attach_url"`
			State     string  `json:"state"`
		}
	)

	err := app.HttpV3("GET", nil, &data, "/apps/%s/dynos", app.AppName)
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
