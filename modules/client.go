package modules

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
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
type filteredTagsResponse struct {
	meta struct {
		status int
		msg    string
	}
	Response struct {
		Filtered_Tags []string
	}
}
type filteredContentsResponse struct {
	meta struct {
		status int
		msg    string
	}
	Response struct {
		Filtered_content []string
	}
}

type TumblrClient struct {
	Client *http.Client
	Config Config
}

func NewTumblrClient(config Config) TumblrClient {
	c := TumblrClient{}
	c.Config = config
	if !c.Config.Testing {
		c.Client = GetClient()
	}
	return c
}

func (c *TumblrClient) GetDashboard(offset int) []npf.Post {
	if c.Config.Testing {
		return npf.TestPosts(20)
	}

	if TokenExpired() {
		c.Client = GetClient()
	}

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
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := dashboardResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response.Posts
}

func (c *TumblrClient) GetTaggedPosts(before int, tag string) []npf.Post {
	if c.Config.Testing {
		return npf.TestPosts(20)
	}

	if TokenExpired() {
		c.Client = GetClient()
	}

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
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := taggedResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response
}

func (c *TumblrClient) GetBlogPosts(before int, blogName string) []npf.Post {
	if c.Config.Testing {
		return npf.TestPosts(20)
	}

	if TokenExpired() {
		c.Client = GetClient()
	}

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			panic(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/blog/" + blogName + ".tumblr.com/posts/text?notes_info=true")

	q := u.Query()
	q.Add("before", strconv.Itoa(before))
	q.Add("npf", "true")

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := dashboardResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response.Posts
}

func (c *TumblrClient) GetFilteredTags(ch chan []string) {
	if c.Config.Testing {
		ch <- []string{}
		return
	}

	if TokenExpired() {
		c.Client = GetClient()
	}

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			panic(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/user/filtered_tags")

	resp, _ := c.Client.Get(u.String())
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := filteredTagsResponse{}
	json.Unmarshal(bytes, &dash)
	ch <- dash.Response.Filtered_Tags
}

func (c *TumblrClient) GetFilteredContents(ch chan []string) {
	if c.Config.Testing {
		ch <- []string{}
		return
	}

	if TokenExpired() {
		c.Client = GetClient()
	}

	defer func() {
		if err := recover(); err != nil {
			RemoveToken()
			print("Failed to retrieve posts\n")
			panic(err)
		}
	}()

	u, _ := url.Parse("https://api.tumblr.com/v2/user/filtered_content")

	resp, _ := c.Client.Get(u.String())
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := filteredContentsResponse{}
	json.Unmarshal(bytes, &dash)
	ch <- dash.Response.Filtered_content
}

func (c *TumblrClient) GetTutorial() []npf.Post {
	tutorialData, err := os.ReadFile("./doc/tutorial.json")
	if err != nil {
		panic(err)
	}
	dash := dashboardResponse{}
	json.Unmarshal(tutorialData, &dash)
	return dash.Response.Posts
}
