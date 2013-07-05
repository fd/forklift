package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (app *App) HttpV3(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.config.Env.CurrentUser.Email
	return app.config.Env.HttpV3(account, method, in, out, path, args...)
}

func (app *App) OwnerHttpV3(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.lookup_owner()
	return app.config.Env.HttpV3(account, method, in, out, path, args...)
}

func (env *Env) HttpV3(account, method string, in, out interface{}, path string, args ...interface{}) error {
	if env.config.DryRun && (method != "GET" && method != "HEAD") {
		return nil
	}

	api_key, err := env.lookup_api_key(account)
	if err != nil {
		return err
	}

	var (
		body_in  io.Reader
		body_out bytes.Buffer
		rawurl   string
	)

	if in != nil {
		buf := bytes.NewBuffer(nil)

		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return err
		}

		body_in = bytes.NewReader(buf.Bytes())
	}

	rawurl = fmt.Sprintf("https://api.heroku.com"+path, args...)

	req, err := http.NewRequest(method, rawurl, body_in)
	if err != nil {
		return err
	}

	req.SetBasicAuth(account, api_key)
	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
	req.Header.Set("User-Agent", "forklift; version=0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_, err = io.Copy(&body_out, resp.Body)
	if err != nil {
		resp.Body.Close()
		return err
	}
	resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		err := &api_error{status: resp.StatusCode}
		if e := json.Unmarshal(body_out.Bytes(), &err); e != nil {
			return err
		}
		return err
	}

	if out != nil {
		err := json.Unmarshal(body_out.Bytes(), out)
		if err != nil {
			return err
		}
	}

	return err
}

func (app *App) lookup_owner() string {
	if app.Owner != "" {
		return app.Owner
	}

	return app.config.Env.CurrentUser.Email
}

func (env *Env) lookup_api_key(email string) (string, error) {
	if email == env.CurrentUser.Email {
		return env.CurrentUser.ApiKey, nil
	}

	for _, owner := range env.OwnerPool {
		if owner.Email == email {
			return owner.ApiKey, nil
		}
	}

	return "", fmt.Errorf("api: unknown heroku account %s", email)
}

type api_error struct {
	status  int
	Id      string `json:"id"`
	Message string `json:"message"`
}

func (a *api_error) Error() string {
	return fmt.Sprintf("api: %s: %s (%d)", a.Id, a.Message, a.status)
}
