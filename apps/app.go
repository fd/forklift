package apps

type App struct {
	Target string

	AppName          string            // name of the heroku app
	Owner            string            // owner of the heroku app
	Upstream         string            // name of the upstream target
	Addons           []string          // list of addons for the app
	Collaborators    []string          // list of collaborators for the app
	Domains          []string          // list of domains for the app
	PostPushCommands []string          // list of commands to execute after pushing to heroku for the app
	Environment      map[string]string // list of config variables for the app
}
