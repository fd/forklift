package sentry

import (
	"fmt"
	"net/url"
	"path"
)

var (
	public_key string
	secret_key string
	endpoint   string
)

func SetDSN(dsn string) error {
	u, err := url.Parse(dsn)
	if err != nil {
		return err
	}

	if u.User == nil {
		return fmt.Errorf("missing keys in: %s", dsn)
	}

	public_key = u.User.Username()
	secret_key, _ = u.User.Password()
	u.User = nil

	dir, project_id := path.Split(u.Path)
	u.Path = path.Join(dir, "api", project_id, "store") + "/"

	endpoint = u.String()
	return nil
}
