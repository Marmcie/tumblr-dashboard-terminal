package main

import (
	"github.com/rivo/tview"
	"tumblr-dt/modules"
)

func main() {
	client := modules.NewTumblrClient()

	app := tview.NewApplication()

	dashboard := modules.NewDashboard(&client, app)

	// Load first sets of posts.
	dashboard.Update()

	if err := app.SetRoot(dashboard.Root, true).Run(); err != nil {
		print("Error in tview loop\n")
		panic(err)
	}
}
