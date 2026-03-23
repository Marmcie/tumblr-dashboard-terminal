package modules

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tumblr/tumblr.go"
)

type dashboardResponse struct {
	Response struct {
		Posts []tumblr.Post
	}
	meta struct {
		status int
		msg    string
	}
}

type TumblrClient struct {
	Client *http.Client
}

func NewTumblrClient() TumblrClient {
	c := TumblrClient{}
	c.Client = GetClient()
	return c
}

func (c *TumblrClient) GetDashboard(offset int) []tumblr.Post {

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			log.Fatal(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/user/dashboard")

	q := u.Query()
	q.Add("offset", strconv.Itoa(offset*20))

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	bytes, _ := io.ReadAll(resp.Body)

	dash := dashboardResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response.Posts
}
