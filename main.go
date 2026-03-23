package main

import (
	"os"
	"tumblr2/modules"

	"github.com/joho/godotenv"
	"github.com/tumblr/tumblrclient.go"

	"github.com/rivo/tview"
)

func main() {

	godotenv.Load("D:/Active projects/tumblr2/.env")
	var consumerkey = os.Getenv("consumer_key")
	var secretkey = os.Getenv("secret_key")
	var oauthtoken = os.Getenv("oauth_token")
	var oauthsecret = os.Getenv("oauth_secret")
	client := tumblrclient.NewClientWithToken(consumerkey, secretkey, oauthtoken, oauthsecret)
	app := tview.NewApplication()

	dashboard := modules.NewDashboard(client, app)

	dashboard.UpdateFeed()
	dashboard.AddPostsToList()
	dashboard.UpdateView()

	if err := app.SetRoot(dashboard.Flex, true).Run(); err != nil {
		panic(err)
	}
}
