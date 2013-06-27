package root

import (
	"fmt"
)

func (cmd *Root) Pause() {
	cmd.SetMaintenance(true)
	fmt.Println("")
	cmd.BreakFormation()
	fmt.Println("")
}

func (cmd *Root) Unpause() {
	fmt.Println("")
	cmd.RestoreFormation()
	fmt.Println("")
	cmd.SetMaintenance(false)
}
