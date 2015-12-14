package apps

import (
	"fmt"
	"strings"

	"github.com/fd/forklift/util/syncset"
)

type (
	plan_t struct {
		Name string `json:"name"`
		name string
	}

	addon_t struct {
		Id   string `json:"id,omitempty"`
		name string
		Plan plan_t `json:"plan"`
	}

	addon_set struct {
		ctx *App

		requested map[string]*addon_t
		current   map[string]*addon_t
	}
)

func (app *App) sync_addons() error {
	set := &addon_set{
		ctx:       app,
		requested: make(map[string]*addon_t, len(app.Addons)),
	}

	for _, name := range app.Addons {
		parts := strings.SplitN(name, ":", 2)
		addon := &addon_t{
			name: parts[0],
			Plan: plan_t{Name: name, name: parts[1]},
		}

		if addon.name == "pgbackups" {
			continue
		}

		set.requested[addon.name] = addon
	}

	fmt.Printf("Addons:\n")

	err := set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *addon_set) LoadCurrentKeys() error {
	var (
		data []*addon_t
	)

	err := set.ctx.HttpV3("GET", nil, &data, "/apps/%s/addons", set.ctx.AppName)
	if err != nil {
		return err
	}

	set.current = make(map[string]*addon_t, len(data))
	for _, addon := range data {
		parts := strings.SplitN(addon.Plan.Name, ":", 2)
		addon.name = parts[0]
		addon.Plan.name = parts[1]
		set.current[addon.name] = addon
	}

	return nil
}

func (set *addon_set) RequestedKeys() []string {
	addons := make([]string, 0, len(set.requested))
	for name := range set.requested {
		addons = append(addons, name)
	}
	return addons
}

func (set *addon_set) CurrentKeys() []string {
	addons := make([]string, 0, len(set.current))
	for name := range set.current {
		addons = append(addons, name)
	}
	return addons
}

func (set *addon_set) ShouldChange(key string) bool {
	if set.current[key].name == "heroku-postgresql" {
		return false
	}

	return set.requested[key].Plan.name != set.current[key].Plan.name
}

func (set *addon_set) Change(key string) (string, string, error) {
	var (
		before = set.current[key].Plan.name
		after  = set.requested[key].Plan.name
		addon  = set.requested[key]
		id     = set.current[key].Id
		err    = set.ctx.HttpV3("PATCH", addon, nil, "/apps/%s/addons/%s", set.ctx.AppName, id)
	)

	return before, after, err
}

func (set *addon_set) Add(name string) error {
	addon := set.requested[name]

	return set.ctx.HttpV3("POST", &addon, nil, "/apps/%s/addons", set.ctx.AppName)
}

func (set *addon_set) Remove(name string) error {
	addon := set.current[name]

	return set.ctx.HttpV3("DELETE", nil, nil, "/apps/%s/addons/%s", set.ctx.AppName, addon.Id)
}
