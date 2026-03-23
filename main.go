package main

import (
	"os"
	"tumblr-dt/modules"

	"github.com/rivo/tview"
)

func main() {
	client := modules.NewTumblrClient()

	app := tview.NewApplication()

	dashboard := modules.NewDashboard(&client, app)

	// Black and white dashboard
	dashboard.BWMode = os.Args[1] == "BW"

	// Load first sets of posts.
	dashboard.Update()
	dashboard.RenderPost()

	if err := app.SetRoot(dashboard.Root, true).Run(); err != nil {
		print("Error in tview loop\n")
		panic(err)
	}
}
