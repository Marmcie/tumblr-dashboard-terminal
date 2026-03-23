package main

import (
	"github.com/rivo/tview"
	"tumblr-terminal/modules"
)

func main() {
	client := modules.NewTumblrClient()

	app := tview.NewApplication()

	dashboard := modules.NewDashboard(&client, app)

	dashboard.Update()

	if err := app.SetRoot(dashboard.Root, true).Run(); err != nil {
		print("Error in tview loop\n")
		panic(err)
	}
}
