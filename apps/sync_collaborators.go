package apps

import (
	"fmt"

	"github.com/fd/forklift/util/syncset"
)

type (
	user_t struct {
		Email string `json:"email"`
	}

	collaborator_t struct {
		Id   string `json:"id,omitempty"`
		User user_t `json:"user"`
	}

	collaborator_set struct {
		ctx *App

		requested     []string
		current       []string
		collaborators map[string]*collaborator_t
	}
)

func (app *App) sync_collaborators() error {
	fmt.Printf("Collaborators:\n")

	err := app.add_owner_to_collaborators()
	if err != nil {
		return err
	}

	set := &collaborator_set{
		ctx:       app,
		requested: app.Collaborators,
	}

	err = set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (app *App) add_owner_to_collaborators() error {
	var (
		resp struct {
			Owner struct {
				Email string `json:"email"`
			} `json:"owner"`
		}
	)

	err := app.OwnerHttpV3("GET", nil, &resp, "/apps/%s", app.AppName)
	if err != nil {
		return err
	}

	if app.Owner == "" {
		app.Owner = resp.Owner.Email
	}

	app.Collaborators = append(app.Collaborators, resp.Owner.Email)
	return nil
}

func (set *collaborator_set) LoadCurrentKeys() error {
	var (
		data []*collaborator_t
	)

	err := set.ctx.OwnerHttpV3("GET", nil, &data, "/apps/%s/collaborators", set.ctx.AppName)
	if err != nil {
		return err
	}

	set.collaborators = make(map[string]*collaborator_t, len(data))
	for _, collaborator := range data {
		email := collaborator.User.Email
		set.collaborators[email] = collaborator
		set.current = append(set.current, email)
	}
	return nil
}

func (set *collaborator_set) RequestedKeys() []string {
	return set.requested
}

func (set *collaborator_set) CurrentKeys() []string {
	return set.current
}

func (set *collaborator_set) ShouldChange(key string) bool {
	return false
}

func (set *collaborator_set) Change(key string) (string, string, error) {
	return "", "", nil
}

func (set *collaborator_set) Add(email string) error {
	collaborator := collaborator_t{
		User: user_t{
			Email: email,
		},
	}

	return set.ctx.OwnerHttpV3("POST", &collaborator, nil, "/apps/%s/collaborators", set.ctx.AppName)
}

func (set *collaborator_set) Remove(email string) error {
	collaborator := set.collaborators[email]

	return set.ctx.OwnerHttpV3("DELETE", nil, nil, "/apps/%s/collaborators/%s", set.ctx.AppName, collaborator.Id)
}
