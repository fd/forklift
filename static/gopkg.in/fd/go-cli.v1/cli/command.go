package cli

import (
	"reflect"
)

/*

Example command:
  type Cmd struct {
    Root|ParentCommand
    Arg0

    FlagsAndEnvs
    Args
  }

*/

type Command interface {
	Main() error
}

type Arg0 string

type Root struct {
	Help bool `flag:"-h,--help"`

	Manual `

    .Help:
      Show the help message for the current command.

  `
}

type Manual struct {
	exec     *executable_t
	usage    string
	summary  string
	options  map[string]section_t
	sections []section_t
}

type section_t struct {
	Header string
	Body   string
}

var (
	typ_Command = reflect.TypeOf((*Command)(nil)).Elem()
	typ_Arg0    = reflect.TypeOf((*Arg0)(nil)).Elem()
	typ_Root    = reflect.TypeOf((*Root)(nil)).Elem()
	typ_Manual  = reflect.TypeOf((*Manual)(nil)).Elem()
)
