package main

import (
	"github.com/rivo/tview"
	"github.com/tumblr/tumblrclient.go"
	"tumblr-terminal/modules"
)

func main() {
	config := modules.GetConfig()

	client := tumblrclient.NewClientWithToken(config.Consumer_key, config.Secret_key, config.Oauth_token, config.Oauth_secret)

	app := tview.NewApplication()

	dashboard := modules.NewDashboard(client, app)

	dashboard.Update()

	if err := app.SetRoot(dashboard.Root, true).Run(); err != nil {
		panic(err)
	}
}
