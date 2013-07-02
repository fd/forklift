package deploy

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
		ctx *Deploy

		requested     []string
		current       []string
		collaborators map[string]*collaborator_t
	}
)

func (cmd *Deploy) sync_collaborators() error {
	set := &collaborator_set{
		ctx:       cmd,
		requested: cmd.Config.Collaborators,
	}

	fmt.Printf("Collaborators:\n")

	err := set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *collaborator_set) LoadCurrentKeys() error {
	var (
		data []*collaborator_t
	)

	err := set.ctx.Http("GET", nil, &data, "/apps/%s/collaborators", set.ctx.Config.Name)
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

	return set.ctx.Http("POST", &collaborator, nil, "/apps/%s/collaborators", set.ctx.Config.Name)
}

func (set *collaborator_set) Remove(email string) error {
	collaborator := set.collaborators[email]

	return set.ctx.Http("DELETE", nil, nil, "/apps/%s/collaborators/%s", set.ctx.Config.Name, collaborator.Id)
}
