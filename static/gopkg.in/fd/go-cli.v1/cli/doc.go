/*

Package cli provides an advanced CLI framework.

Commands are defined as structs wich must be registerd with Register(T).

  type Command struct {
    // The parent command must always be present and the first field
    // It defaults to cli.Root ban it can be any command type.
    cli.Root

    // The name of the command as it was passed to exec()
    cli.Arg0

    // flag: env: and arg can be combined
    Field Type `flag:"-f,--flag"`
    Field Type `env:"VAR"`
    Field Type `arg`

    // The manual for this command
    cli.Manual `
      Usage:   a usage example
      Summary: a short decription of the command

      .Field:
        A description for the field Field.

      Section:
        A manual section and header.
    `
  }

  func init() {
    // register the command. (must happen during init())
    cli.Register(Command{})
  }

*/
package cli
