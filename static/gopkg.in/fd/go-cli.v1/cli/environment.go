package cli

import (
	"strings"
)

type environment_t struct {
	args []string
	env  map[string]string
}

func new_environment(args []string, env []string) *environment_t {
	return &environment_t{args: args, env: parse_env(env)}
}

func parse_env(in []string) map[string]string {
	m := make(map[string]string, len(in))

	for _, l := range in {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 1 {
			m[parts[0]] = ""
		} else {
			m[parts[0]] = parts[1]
		}
	}

	return m
}
