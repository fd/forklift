package apps

import (
	"fmt"
)

func (app *App) Pause() {
	app.SetMaintenance(true)
	fmt.Println("")
	app.FormationBreak()
	fmt.Println("")
}

func (app *App) Unpause() {
	fmt.Println("")
	app.FormationRestore()
	fmt.Println("")
	app.SetMaintenance(false)
}
