package setup

import (
	"fmt"
	"strings"

	"github.com/fd/forklift/util/syncset"
)

type (
	config_set struct {
		ctx *Setup

		requested map[string]string
		current   map[string]string
	}
)

func (cmd *Setup) sync_config() error {
	set := &config_set{
		ctx:       cmd,
		requested: cmd.Config.Config,
	}

	fmt.Printf("Config:\n")

	err := set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *config_set) LoadCurrentKeys() error {
	var (
		data map[string]string
	)

	err := set.ctx.Http("GET", nil, &data, "/apps/%s/config-vars", set.ctx.Config.Name)
	if err != nil {
		return err
	}

	for key := range data {
		ignore := false

		if strings.HasPrefix(key, "HEROKU_POSTGRESQL_") {
			ignore = true
		}

		if key == "DATABASE_URL" {
			ignore = true
		}

		if key == "GEM_PATH" {
			ignore = true
		}

		if key == "LANG" {
			ignore = true
		}

		if key == "PATH" {
			ignore = true
		}

		if key == "PGBACKUPS_URL" {
			ignore = true
		}

		if key == "REDISTOGO_URL" {
			ignore = true
		}

		if ignore {
			delete(data, key)
		}
	}

	set.current = data
	return nil
}

func (set *config_set) RequestedKeys() []string {
	keys := make([]string, 0, len(set.requested))
	for key := range set.requested {
		keys = append(keys, key)
	}
	return keys
}

func (set *config_set) CurrentKeys() []string {
	keys := make([]string, 0, len(set.current))
	for key := range set.current {
		keys = append(keys, key)
	}
	return keys
}

func (set *config_set) ShouldChange(key string) bool {
	return set.current[key] != set.requested[key]
}

func (set *config_set) Add(key string) error {
	data := map[string]interface{}{
		key: set.requested[key],
	}

	return set.ctx.Http("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.Config.Name)
}

func (set *config_set) Change(key string) (string, string, error) {
	var (
		before = set.current[key]
		after  = set.requested[key]
		data   = map[string]interface{}{key: after}
		err    = set.ctx.Http("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.Config.Name)
	)

	return before, after, err
}

func (set *config_set) Remove(key string) error {
	data := map[string]interface{}{
		key: nil,
	}

	return set.ctx.Http("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.Config.Name)
}
