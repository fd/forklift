package cli

import (
	// "fmt"
	"reflect"
	"strings"
)

var executables = map[reflect.Type]*executable_t{}

func init() {
	Register(Root{})
}

type executable_t struct {
	Type          reflect.Type
	IsGroup       bool
	Names         []string
	ParentExec    *executable_t
	Parent        reflect.StructField
	Arg0          reflect.StructField
	manual        reflect.StructField
	parsed_manual *Manual
	Variables     []reflect.StructField
	Flags         []reflect.StructField
	Args          []reflect.StructField
	SubCommands   []*executable_t
}

func Register(v interface{}) {
	RegisterType(reflect.TypeOf(v))
}

func RegisterType(t reflect.Type) {
	if _, p := executables[t]; p {
		return
	}

	var (
		exec = &executable_t{}
	)
	executables[t] = exec
	exec.Type = t

	if t.Kind() != reflect.Struct {
		panic("command must be a struct")
	}

	c := -1
	for i, l := 0, t.NumField(); i < l; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		c++
		f.Index = []int{i}

		// first public field
		if c == 0 && t != typ_Root {
			exec.Parent = f
			continue
		}

		if f.Type == typ_Arg0 {
			if tag := f.Tag.Get("name"); tag != "" {
				exec.Names = strings.Split(tag, ",")
			}
			exec.Arg0 = f
			continue
		}

		if f.Type == typ_Manual {
			exec.manual = f
			continue
		}

		if f.Tag.Get("env") != "" {
			exec.Variables = append(exec.Variables, f)
		}

		if f.Tag.Get("flag") != "" {
			exec.Flags = append(exec.Flags, f)
		}

		if strings.Contains(string(f.Tag), "arg") {
			exec.Args = append(exec.Args, f)
		}
	}

	if !reflect.PtrTo(t).Implements(typ_Command) {
		exec.IsGroup = true
	}

	if exec.Parent.Type != nil {
		RegisterType(exec.Parent.Type)
		parent_exec := executables[exec.Parent.Type]
		if !parent_exec.IsGroup {
			panic("the parent command must be a group")
		}
		exec.ParentExec = parent_exec
		parent_exec.SubCommands = append(parent_exec.SubCommands, exec)
	}
}
