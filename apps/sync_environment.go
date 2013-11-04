package apps

import (
	"fmt"
	"strings"

	"github.com/fd/forklift/util/syncset"
)

type (
	env_var_set struct {
		ctx *App

		requested map[string]string
		current   map[string]string
	}
)

var ignored_keys = map[string]bool{
	"DATABASE_URL": true,
	"GEM_PATH":     true,
	"LANG":         true,
	"PATH":         true,

	"HEROKU_POSTGRESQL_*": true,
	"MEMCACHIER_*":        true,
	"PGBACKUPS_*":         true,
	"REDISCLOUD_*":        true,
	"REDISTOGO_*":         true,
	"CLEARDB_*":           true,
	"NEW_RELIC_L*":        true,
}

func (app *App) sync_config() error {
	set := &env_var_set{
		ctx:       app,
		requested: app.Environment,
	}

	fmt.Printf("Environment:\n")

	err := set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *env_var_set) LoadCurrentKeys() error {
	var (
		data map[string]string
	)

	err := set.ctx.HttpV3("GET", nil, &data, "/apps/%s/config-vars", set.ctx.AppName)
	if err != nil {
		return err
	}

	for key := range data {
		ignore := false

		for pattern := range ignored_keys {
			if strings.HasSuffix(pattern, "*") {
				ignore = strings.HasPrefix(key, strings.TrimSuffix(pattern, "*"))
			} else {
				ignore = pattern == key
			}
			if ignore {
				break
			}
		}

		if ignore {
			delete(data, key)
		}
	}

	set.current = data
	return nil
}

func (set *env_var_set) RequestedKeys() []string {
	keys := make([]string, 0, len(set.requested))
	for key := range set.requested {
		keys = append(keys, key)
	}
	return keys
}

func (set *env_var_set) CurrentKeys() []string {
	keys := make([]string, 0, len(set.current))
	for key := range set.current {
		keys = append(keys, key)
	}
	return keys
}

func (set *env_var_set) ShouldChange(key string) bool {
	return set.current[key] != set.requested[key]
}

func (set *env_var_set) Add(key string) error {
	data := map[string]interface{}{
		key: set.requested[key],
	}

	return set.ctx.HttpV3("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.AppName)
}

func (set *env_var_set) Change(key string) (string, string, error) {
	var (
		before = set.current[key]
		after  = set.requested[key]
		data   = map[string]interface{}{key: after}
		err    = set.ctx.HttpV3("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.AppName)
	)

	return before, after, err
}

func (set *env_var_set) Remove(key string) error {
	data := map[string]interface{}{
		key: nil,
	}

	return set.ctx.HttpV3("PATCH", &data, nil, "/apps/%s/config-vars", set.ctx.AppName)
}
