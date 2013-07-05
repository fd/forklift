package apps

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	toml "github.com/pelletier/go-toml"
)

func DecodeJSON(r io.Reader) (*Config, error) {
	var (
		v   interface{}
		err error
		cnf = &Config{}
	)

	err = json.NewDecoder(r).Decode(&v)
	if err != nil {
		return nil, err
	}

	err = populate_cnf(v, cnf)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

func EncodeJSON(w io.Writer, cnf *Config) error {
	root := make(map[string]interface{}, 50)
	pack_cnf(root, cnf)
	return json.NewEncoder(w).Encode(root)
}

func DecodeTOML(r io.Reader) (*Config, error) {
	var (
		v    interface{}
		err  error
		data []byte
		tree *toml.TomlTree
		cnf  = &Config{}
	)

	data, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tree, err = toml.Load(string(data))
	if err != nil {
		return nil, err
	}

	v = map[string]interface{}(*tree)

	err = populate_cnf(v, cnf)
	if err != nil {
		return nil, err
	}

	return cnf, nil
}

func populate_cnf(v interface{}, cnf *Config) (err error) {
	defer func(err *error) {
		if e := recover(); e != nil {
			*err = fmt.Errorf("parse error: %s", e)
		}
	}(&err)

	root, ok := v.(map[string]interface{})
	if !ok {
		panic("expected a map for `.`")
	}

	populate_app(root, &cnf.App)
	populate_env(root, &cnf.Env)

	for key, v := range root {
		switch key {

		case "deploypack", "Deploypack":
			cnf.Deploypack = must_string("deploypack", v)
			delete(root, key)

		default:
			// ignore unknown keys

		}
	}

	cnf.App.config = cnf
	cnf.Env.config = cnf
	cnf.Unused = root

	return
}

func populate_app(root map[string]interface{}, app *App) {
	for key, v := range root {
		switch key {

		case "name", "Name", "app", "app_name", "AppName":
			app.AppName = must_string("name", v)
			delete(root, key)

		case "owner", "Owner":
			app.Owner = must_string("owner", v)
			delete(root, key)

		case "upstream", "Upstream":
			app.Upstream = must_string("upstream", v)
			delete(root, key)

		case "addons", "Addons":
			app.Addons = must_string_slice("addons", v)
			delete(root, key)

		case "collaborators", "Collaborators":
			app.Collaborators = must_string_slice("collaborators", v)
			delete(root, key)

		case "domains", "Domains":
			app.Domains = must_string_slice("domains", v)
			delete(root, key)

		case "post_push_commands", "PostPushCommands":
			app.PostPushCommands = must_string_slice("post_push_commands", v)
			delete(root, key)

		case "environment", "Environment":
			app.Environment = must_string_map("environment", v)
			delete(root, key)

		default:
			// ignore unknown keys

		}
	}
}

func populate_env(root map[string]interface{}, env *Env) {
	for key, v := range root {
		switch key {

		case "owners", "Owners", "owner_pool", "OwnerPool":
			env.OwnerPool = must_account_slice("owners", v)
			delete(root, key)

		default:
			// ignore unknown keys

		}
	}
}

func pack_cnf(root map[string]interface{}, cnf *Config) {
	if cnf.Deploypack != "" {
		root["deploypack"] = cnf.Deploypack
	}

	pack_app(root, &cnf.App)
	pack_env(root, &cnf.Env)

	for k, v := range cnf.Unused {
		root[k] = v
	}
}

func pack_app(root map[string]interface{}, app *App) {
	if app.AppName != "" {
		root["name"] = app.AppName
	}

	if app.Owner != "" {
		root["owner"] = app.Owner
	}

	if app.Upstream != "" {
		root["upstream"] = app.Upstream
	}

	if len(app.Addons) > 0 {
		root["addons"] = app.Addons
	}

	if len(app.Collaborators) > 0 {
		root["collaborators"] = app.Collaborators
	}

	if len(app.Domains) > 0 {
		root["domains"] = app.Domains
	}

	if len(app.PostPushCommands) > 0 {
		root["post_push_commands"] = app.PostPushCommands
	}

	if len(app.Environment) > 0 {
		root["environment"] = app.Environment
	}
}

func pack_env(root map[string]interface{}, env *Env) {
	if len(env.OwnerPool) > 0 {
		root["owner_pool"] = env.OwnerPool
	}
}

func must_string(key string, v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	panic(fmt.Sprintf("expected a string for `%s` (instead of %T)", key, v))
}

func must_string_slice(key string, v interface{}) []string {
	if vl, ok := v.([]interface{}); ok {
		var (
			sl   = make([]string, len(vl))
			skey = key + "[]"
		)

		for i, v := range vl {
			sl[i] = must_string(skey, v)
		}

		return sl
	}
	panic(fmt.Sprintf("expected an array for `%s` (instead of %T)", key, v))
}

func must_string_map(key string, v interface{}) map[string]string {
	if vl, ok := v.(map[string]interface{}); ok {
		var (
			sl   = make(map[string]string, len(vl))
			skey = key + "[k]"
		)

		for k, v := range vl {
			sl[k] = must_string(skey, v)
		}

		return sl
	}
	panic(fmt.Sprintf("expected a map for `%s` (instead of %T)", key, v))
}

func must_account(key string, v interface{}) Account {
	if root, ok := v.(map[string]interface{}); ok {
		account := Account{}

		for key, v := range root {
			switch key {

			case "email", "Email", "user", "User", "account", "Account":
				account.Email = must_string("email", v)

			case "api_key", "apikey", "ApiKey":
				account.ApiKey = must_string("api_key", v)

			default:
				// ignore unknown keys

			}
		}

		return account
	}
	panic(fmt.Sprintf("expected a map for `%s` (instead of %T)", key, v))
}

func must_account_slice(key string, v interface{}) []Account {
	if vl, ok := v.([]interface{}); ok {
		var (
			sl   = make([]Account, len(vl))
			skey = key + "[]"
		)

		for i, v := range vl {
			sl[i] = must_account(skey, v)
		}

		return sl
	}
	panic(fmt.Sprintf("expected an array for `%s` (instead of %T)", key, v))
}
