# Forklift

More managable application deployment for heroku.

## Configuration

All forklift configuration happens in _target files_. Target files live in the `.forklift` directory at the root of your project. Target files use the [TOML](https://github.com/fd/forklift/static/github.com/mojombo/toml) format.

| Key | Type | Description |
| --- | ---- | ----------- |
| `name` | `string` | The name of the Heroku application. |
| `owner` | `string` | The email address of the Owner of the application. |
| `upstream` | `string` | An upstream target to follow. When this key is set deploys will only deploy code that has already been deployed in the _upstream target_. |
| `addons` | `[]string` | A list of addons that must be provisioned. |
| `collaborators` | `[]string` | A list of collaborators that must have access to this app. |
| `domains` | `[]string` | A list of domains that must be added to this app. |
| `features` | `[]string` | A list of heroku labs features. |
| `post_push_commands` | `[]string` | A list of commands that must be run (on heroku using `heroku run ...`) after new code has been pushed. |
| `environment` | `map[string]string` | A map of configuration variables that must be added to the environment of the application. |
| `deploypack` | `string` | The Deploypack to use when deploying this application. |

## TODO

- `pre push hooks`
- `check`
