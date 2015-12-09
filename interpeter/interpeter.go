package interpeter

import (
	"fmt"

	"github.com/fd/forklift/static/github.com/zhemao/glisp/interpreter"
)

type Interperter struct {
	Env map[string]string
	env *glisp.Glisp
}

func (i *Interperter) SexpString() string {
	return "forklift"
}

func (i *Interperter) setup() {
	i.env.AddGlobal("forklift", i)

	i.env.AddFunction("include", includeDeploypack)
	i.env.AddFunction("env-set", envSet)
	i.env.AddFunction("env-get", envGet)
}

func envSet(env *glisp.Glisp, name string, args []glisp.Sexp) (glisp.Sexp, error) {
	if len(args) == 0 || len(args)%2 != 0 {
		return nil, fmt.Errorf("usage: (env-set <key> <value> ...)")
	}

	for i := 0; i < len(args); i += 2 {
		if !glisp.IsString(args[i]) || !glisp.IsString(args[i+1]) {
			return nil, fmt.Errorf("usage: (env-set <key> <value> ...)")
		}

		var (
			interp = getInterperter(env)
			key    = string(args[i].(glisp.SexpStr))
			value  = string(args[i+1].(glisp.SexpStr))
		)

		interp.Env[key] = value
	}

	return glisp.SexpBool(true), nil
}

func envGet(env *glisp.Glisp, name string, args []glisp.Sexp) (glisp.Sexp, error) {
	if len(args) != 1 || !glisp.IsString(args[0]) {
		return nil, fmt.Errorf("usage: (env-get <key>)")
	}

	var (
		interp = getInterperter(env)
		key    = string(args[0].(glisp.SexpStr))
	)

	return glisp.SexpStr(interp.Env[key]), nil
}

func getInterperter(env *glisp.Glisp) *Interperter {
	x, ok := env.FindObject("forklift")
	if !ok {
		panic("forklift not found")
	}

	i, ok := x.(*Interperter)
	if !ok {
		panic("forklift not found")
	}

	return i
}
