package apps

import (
	"fmt"
	"path"

	"code.google.com/p/go-netrc/netrc"

	"github.com/fd/forklift/util/user"
)

type Env struct {
	config *Config

	CurrentUser Account   // the current heroku account
	OwnerPool   []Account // list of heroku accounts to use when creating a nuw app
}

func (env *Env) load_heroku_credentials() error {
	home, err := user.Home()
	if err != nil {
		return err
	}

	machines, _, err := netrc.ParseFile(path.Join(home, ".netrc"))
	if err != nil {
		return err
	}

	for _, machine := range machines {
		if machine.Name == "api.heroku.com" {
			env.CurrentUser.Email = machine.Login
			env.CurrentUser.ApiKey = machine.Password
			return nil
		}
	}

	return fmt.Errorf("Please run `heroku login`")
}
