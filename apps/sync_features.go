package apps

import (
	"fmt"

	"github.com/fd/forklift/util/syncset"
)

type (
	feature_set struct {
		ctx *App

		requested []string
		current   []string
	}
)

func (app *App) sync_features() error {
	fmt.Printf("Features:\n")

	set := &feature_set{
		ctx:       app,
		requested: app.Features,
	}

	err := set.LoadCurrentFeatures()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *feature_set) LoadCurrentFeatures() error {
	var (
		data []struct {
			Name    string
			Enabled bool
			Kind    string
		}
	)

	err := set.ctx.HttpV2("GET", nil, &data, "/features?app=%s", set.ctx.AppName)
	if err != nil {
		return err
	}

	set.current = make([]string, 0, len(data))
	for _, feature := range data {
		if !feature.Enabled {
			continue
		}
		if feature.Kind != "app" {
			continue
		}
		set.current = append(set.current, feature.Name)
	}

	return nil
}

func (set *feature_set) RequestedKeys() []string {
	return set.requested
}

func (set *feature_set) CurrentKeys() []string {
	return set.current
}

func (set *feature_set) ShouldChange(key string) bool {
	return false
}

func (set *feature_set) Change(key string) (string, string, error) {
	return "", "", nil
}

func (set *feature_set) Add(feature string) error {
	return set.ctx.HttpV2("POST", nil, nil, "/features/%s?app=%s", feature, set.ctx.AppName)
}

func (set *feature_set) Remove(feature string) error {
	return set.ctx.HttpV2("DELETE", nil, nil, "/features/%s?app=%s", feature, set.ctx.AppName)
}
