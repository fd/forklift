package cli

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type context_t int

const (
	ctx_env context_t = 1 + iota
	ctx_arg
)

func Main(args []string, environ []string) {
	if args == nil {
		args = os.Args
	}

	if environ == nil {
		environ = os.Environ()
	}

	env := new_environment(args, environ)
	exec := executables[typ_Root]
	err := execute(env, exec, reflect.Value{})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func execute(env *environment_t, exec *executable_t, parent reflect.Value) error {
	var (
		pv  = reflect.New(exec.Type)
		v   = reflect.Indirect(pv)
		fv  reflect.Value
		err error
	)

	if parent.IsValid() {
		fv = v.Field(exec.Parent.Index[0])
		fv.Set(parent)

		if len(env.args) == 0 {
			return fmt.Errorf("expected a command")
		}
		arg := env.args[0]
		if exec.Arg0.Type != nil {
			fv = v.Field(exec.Arg0.Index[0])
			fv.SetString(arg)
		}
		env.args = env.args[1:]
	}

	for _, f := range exec.Variables {
		fv = v.Field(f.Index[0])
		err = handle_env(env, f, fv)
		if err != nil {
			return err
		}
	}

	for _, f := range exec.Flags {
		fv = v.Field(f.Index[0])
		err = handle_flag(env, f, fv)
		if err != nil {
			return err
		}
	}

	if m := exec.Manual(); m != nil {
		f := exec.manual
		fv = v.Field(f.Index[0])
		fv.Set(reflect.ValueOf(*m))
	}

	if exec.IsGroup {
		return execute_group(env, exec, pv)
	} else {
		return execute_command(env, exec, pv)
	}

	panic("not reached")
}

func execute_command(env *environment_t, exec *executable_t, pv reflect.Value) error {
	var (
		v   = reflect.Indirect(pv)
		fv  reflect.Value
		err error
	)

	for _, f := range exec.Args {
		fv = v.Field(f.Index[0])
		err = handle_arg(env, f, fv)
		if err != nil {
			return err
		}
	}

	if len(env.args) != 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(env.args, " "))
	}

	if hv := v.FieldByName("Help"); hv.Bool() {
		exec.Manual().Open()
		return nil
	}

	return pv.Interface().(Command).Main()
}

func execute_group(env *environment_t, exec *executable_t, pv reflect.Value) error {
	var (
		v   = reflect.Indirect(pv)
		arg string
	)

	if len(env.args) == 0 {
		if hv := v.FieldByName("Help"); hv.Bool() {
			exec.Manual().Open()
			return nil
		}
		return fmt.Errorf("expected a command")
	}
	arg = env.args[0]

	for _, sub_exec := range exec.SubCommands {
		if len(sub_exec.Names) > 0 {
			found := false
			for _, n := range sub_exec.Names {
				if n == arg {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}

		return execute(env, sub_exec, v)
	}

	return fmt.Errorf("unexpected arguments: %s", strings.Join(env.args, " "))
}

func handle_env(env *environment_t, f reflect.StructField, fv reflect.Value) error {
	names := strings.Split(f.Tag.Get("env"), ",")

	for _, n := range names {
		val, p := env.env[n]
		if p {
			_, err := set_value_with_args(fv, []string{val}, ctx_env)
			return err
		}
	}

	return nil
}

func handle_flag(env *environment_t, f reflect.StructField, fv reflect.Value) error {
	var (
		names        = strings.Split(f.Tag.Get("flag"), ",")
		args         = env.args
		skipped_args []string
	)

	for _, n := range names {
		for i, l := 0, len(args); i < l; i++ {
			arg := args[i]

			if arg != n {
				skipped_args = append(skipped_args, arg)
				continue
			}

			n, err := set_value_with_args(fv, args[i+1:], ctx_arg)
			if err != nil {
				return err
			}

			i += n
		}

		args = skipped_args
		skipped_args = nil
	}

	env.args = args
	return nil
}

func handle_arg(env *environment_t, f reflect.StructField, fv reflect.Value) error {
	var (
		args = env.args
	)

	n, err := set_value_with_args(fv, args, ctx_arg)
	if err != nil {
		return err
	}

	env.args = args[n:]
	return nil
}

func set_value_with_args(fv reflect.Value, args []string, ctx context_t) (int, error) {
	switch fv.Kind() {

	case reflect.Bool:
		if ctx == ctx_arg {
			fv.SetBool(true)
			return 0, nil
		}
		if args[0] == "t" || args[0] == "true" || args[0] == "yes" || args[0] == "y" {
			fv.SetBool(true)
			return 1, nil
		}
		return 0, nil

	case reflect.String:
		if len(args) > 0 {
			fv.SetString(args[0])
			return 1, nil
		}
		return 0, fmt.Errorf("expected an argument")

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(args) > 0 {
			i, err := strconv.ParseInt(args[0], 10, fv.Type().Bits())
			if err != nil {
				return 0, err
			}
			fv.SetInt(i)
			return 1, nil
		}
		return 0, fmt.Errorf("expected an argument")

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(args) > 0 {
			i, err := strconv.ParseUint(args[0], 10, fv.Type().Bits())
			if err != nil {
				return 0, err
			}
			fv.SetUint(i)
			return 1, nil
		}
		return 0, fmt.Errorf("expected an argument")

	case reflect.Float32, reflect.Float64:
		if len(args) > 0 {
			i, err := strconv.ParseFloat(args[0], fv.Type().Bits())
			if err != nil {
				return 0, err
			}
			fv.SetFloat(i)
			return 1, nil
		}
		return 0, fmt.Errorf("expected an argument")

	}

	panic("Unsupported type: " + fv.Type().String())
}
