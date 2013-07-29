package apps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fd/go-util/errors"
)

func (app *App) HttpV2(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.config.Env.CurrentUser.Email
	return app.config.Env.HttpV2(account, method, in, out, path, args...)
}

func (app *App) OwnerHttpV2(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.lookup_owner()
	return app.config.Env.HttpV2(account, method, in, out, path, args...)
}

func (env *Env) HttpV2(account, method string, in, out interface{}, path string, args ...interface{}) error {
	if !env.perform_request(method) {
		return nil
	}

	req, err := env.new_request(account, method, in, path, args...)
	if err != nil {
		return err
	}

	err = env.do_request(req, account, in, out)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) HttpV3(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.config.Env.CurrentUser.Email
	return app.config.Env.HttpV3(account, method, in, out, path, args...)
}

func (app *App) OwnerHttpV3(method string, in, out interface{}, path string, args ...interface{}) error {
	account := app.lookup_owner()
	return app.config.Env.HttpV3(account, method, in, out, path, args...)
}

func (env *Env) HttpV3(account, method string, in, out interface{}, path string, args ...interface{}) error {
	if !env.perform_request(method) {
		return nil
	}

	req, err := env.new_request(account, method, in, path, args...)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.heroku+json; version=3")

	err = env.do_request(req, account, in, out)
	if err != nil {
		return err
	}

	return nil
}

func (env *Env) perform_request(method string) bool {
	if method == "GET" || method == "HEAD" {
		return true
	}
	if env.config.DryRun {
		return false
	}
	return true
}

func (env *Env) new_request(account, method string, in interface{}, path string, args ...interface{}) (*http.Request, error) {
	api_key, err := env.lookup_api_key(account)
	if e := errors.Annotate(err, "API error"); e != nil {
		e.AddContext("account=%s", account)
		return nil, e
	}

	var (
		body_in io.Reader
		rawurl  string
	)

	if in != nil {
		buf := bytes.NewBuffer(nil)

		err := json.NewEncoder(buf).Encode(in)
		if e := errors.Annotate(err, "API error"); e != nil {
			e.AddContext("request.body=%+v", in)
			return nil, e
		}

		body_in = bytes.NewReader(buf.Bytes())
	}

	rawurl = fmt.Sprintf("https://api.heroku.com"+path, args...)

	req, err := http.NewRequest(method, rawurl, body_in)
	if e := errors.Annotate(err, "API error"); e != nil {
		e.AddContext("account=%s", account)
		e.AddContext("url=%s", rawurl)
		e.AddContext("method=%s", method)
		e.AddContext("request.body=%+v", in)
		return nil, e
	}

	req.SetBasicAuth(account, api_key)
	req.Header.Set("User-Agent", "forklift; version=0")

	return req, nil
}

func (env *Env) do_request(req *http.Request, account string, in, out interface{}) error {
	var (
		body_out bytes.Buffer
	)

	resp, err := http.DefaultClient.Do(req)
	if e := errors.Annotate(err, "API error"); e != nil {
		e.AddContext("account=%s", account)
		e.AddContext("url=%s", req.URL.String())
		e.AddContext("method=%s", req.Method)
		e.AddContext("request.body=%+v", in)
		return e
	}

	_, err = io.Copy(&body_out, resp.Body)
	resp.Body.Close()
	if e := errors.Annotate(err, "API error"); e != nil {
		e.AddContext("account=%s", account)
		e.AddContext("url=%s", req.URL.String())
		e.AddContext("method=%s", req.Method)
		e.AddContext("status=%d", resp.StatusCode)
		e.AddContext("content_type=%d", resp.Header.Get("content-type"))
		e.AddContext("request.body=%+v", in)
		return e
	}

	if resp.StatusCode/100 != 2 {
		err_resp := error_resp{}
		json.Unmarshal(body_out.Bytes(), &err_resp)

		e := errors.New("API error")
		e.AddContext("account=%s", account)
		e.AddContext("url=%s", req.URL.String())
		e.AddContext("method=%s", req.Method)
		e.AddContext("status=%d", resp.StatusCode)
		e.AddContext("content_type=%d", resp.Header.Get("content-type"))
		e.AddContext("request.body=%+v", in)
		e.AddContext("error.id=%s", err_resp.Id)
		e.AddContext("error.message=%s", err_resp.Message)
		return e
	}

	if out != nil {
		err := json.Unmarshal(body_out.Bytes(), out)
		if e := errors.Annotate(err, "API error"); e != nil {
			e.AddContext("account=%s", account)
			e.AddContext("url=%s", req.URL.String())
			e.AddContext("method=%s", req.Method)
			e.AddContext("status=%d", resp.StatusCode)
			e.AddContext("content_type=%d", resp.Header.Get("content-type"))
			e.AddContext("request.body=%+v", in)
			e.AddContext("response.body=%s", body_out.Bytes())
			return e
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

type error_resp struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}
