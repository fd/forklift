package syncset

import (
	"fmt"
	"sort"
)

type Set interface {
	RequestedKeys() []string
	CurrentKeys() []string
	ShouldChange(key string) bool

	Add(key string) error
	Change(key string) (before, after string, erro error)
	Remove(key string) error
}

const (
	state_add    = 1
	state_del    = 2
	state_change = 3
	state_keep   = 4
)

func Sync(set Set) bool {
	var (
		requested_keys = set.RequestedKeys()
		current_keys   = set.CurrentKeys()
		key_state      = make(map[string]uint8, len(requested_keys)+len(current_keys))
		keys           = make([]string, 0, len(key_state))
		ok             = true
	)

	for _, key := range requested_keys {
		key_state[key] = state_add
	}

	for _, key := range current_keys {
		if _, p := key_state[key]; p {
			if set.ShouldChange(key) {
				key_state[key] = state_change
			} else {
				key_state[key] = state_keep
			}
		} else {
			key_state[key] = state_del
		}
	}

	for key, _ := range key_state {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		fmt.Printf(" - (empty)\n")
		return ok
	}

	sort.Strings(keys)

	for _, key := range keys {
		switch key_state[key] {

		case state_keep:
			fmt.Printf(" - %s\n", key)

		case state_change:
			before, after, err := set.Change(key)
			if err != nil {
				ok = false
				fmt.Printf(" - %s \x1b[31;40;4;5m(failed to change)\x1b[0m\n   before: %s\n   after:  %s\n   error:  %s\n",
					key, before, after, err)
			} else {
				fmt.Printf(" - %s \x1b[33m(changed)\x1b[0m\n  before: %s\n  after:  %s\n",
					key, before, after)
			}

		case state_add:
			err := set.Add(key)
			if err != nil {
				ok = false
				fmt.Printf(" - %s \x1b[31;40;4;5m(failed to add)\x1b[0m\n   error: %s\n", key, err)
			} else {
				fmt.Printf(" - %s \x1b[32m(added)\x1b[0m\n", key)
			}

		case state_del:
			err := set.Remove(key)
			if err != nil {
				ok = false
				fmt.Printf(" - %s \x1b[31;40;4;5m(failed to remove)\x1b[0m\n   error: %s\n", key, err)
			} else {
				fmt.Printf(" - %s \x1b[31m(removed)\x1b[0m\n", key)
			}

		}
	}

	return ok
}
