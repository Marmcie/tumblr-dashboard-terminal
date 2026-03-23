package modules

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"tumblr-dt/npf"
)

type dashboardResponse struct {
	Response struct {
		Posts []npf.Post
	}
	meta struct {
		status int
		msg    string
	}
}
type taggedResponse struct {
	Response []npf.Post
	meta     struct {
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

func (c *TumblrClient) GetDashboard(offset int) []npf.Post {

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			panic(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/user/dashboard")

	q := u.Query()
	q.Add("offset", strconv.Itoa(offset*20))
	q.Add("reblog_info", "true")
	q.Add("npf", "true")

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	bytes, _ := io.ReadAll(resp.Body)

	dash := dashboardResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response.Posts
}

func (c *TumblrClient) GetTaggedPosts(before int, tag string) []npf.Post {

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			panic(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/tagged")

	q := u.Query()
	q.Add("before", strconv.Itoa(before))
	q.Add("tag", tag)
	q.Add("npf", "true")

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	bytes, _ := io.ReadAll(resp.Body)

	dash := taggedResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response
}
