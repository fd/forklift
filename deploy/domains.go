package deploy

import (
	"fmt"

	"github.com/fd/forklift/util/syncset"
)

type (
	domain_t struct {
		Id       string `json:"id,omitempty"`
		Hostname string `json:"hostname"`
	}

	domain_set struct {
		ctx *Deploy

		requested []string
		current   []string
		domains   map[string]*domain_t
	}
)

func (cmd *Deploy) sync_domains() error {
	set := &domain_set{
		ctx:       cmd,
		requested: cmd.Config.Domains,
	}

	fmt.Printf("Domain:\n")

	err := set.LoadCurrentKeys()
	if err != nil {
		return err
	}

	syncset.Sync(set)
	return nil
}

func (set *domain_set) LoadCurrentKeys() error {
	var (
		data     []*domain_t
		mainhost = set.ctx.Config.Name + ".herokuapp.com"
	)

	err := set.ctx.Http("GET", nil, &data, "/apps/%s/domains", set.ctx.Config.Name)
	if err != nil {
		return err
	}

	set.domains = make(map[string]*domain_t, len(data))
	for _, domain := range data {
		host := domain.Hostname

		if host == mainhost {
			continue
		}

		set.domains[host] = domain
		set.current = append(set.current, host)
	}
	return nil
}

func (set *domain_set) RequestedKeys() []string {
	return set.requested
}

func (set *domain_set) CurrentKeys() []string {
	return set.current
}

func (set *domain_set) ShouldChange(key string) bool {
	return false
}

func (set *domain_set) Change(key string) (string, string, error) {
	return "", "", nil
}

func (set *domain_set) Add(host string) error {
	domain := domain_t{
		Hostname: host,
	}

	return set.ctx.Http("POST", &domain, nil, "/apps/%s/domains", set.ctx.Config.Name)
}

func (set *domain_set) Remove(host string) error {
	domain := set.domains[host]

	return set.ctx.Http("DELETE", nil, nil, "/apps/%s/domains/%s", set.ctx.Config.Name, domain.Id)
}
