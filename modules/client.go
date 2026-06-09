package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"tumblr-dt/npf"
	"tumblr-dt/ui/helper"
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
type timelimeResponse struct {
	Response struct {
		Timeline struct {
			Elements []npf.Post
			Links   struct {
				Next struct {
					Href string
				}
			} `json:"_links"`
		}
	}
	meta struct {
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
	return c.GetSearchedPosts(before, tag)
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
	q.Add("npf", "true")
	q.Add("tag", tag)

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	dash := taggedResponse{}
	json.Unmarshal(bytes, &dash)
	return dash.Response
}

func (c *TumblrClient) GetSearchedPosts(before int, term string) []npf.Post {
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

	u, _ := url.Parse("https://www.tumblr.com/api/v2/timeline/search")

	q := u.Query()
	q.Add("before", strconv.Itoa(before))
	q.Add("npf", "true")
	q.Add("timeline_type", "post")
	q.Add("query_source", "search_box_typed_query")
	q.Add("post_role", "any")
	q.Add("query", term)
	q.Add("limit", "20")
	q.Add("days", "0")

	u.RawQuery = q.Encode()

	resp, _ := c.Client.Get(u.String())
	defer resp.Body.Close()
	bytes, _ := io.ReadAll(resp.Body)

	helper.Log(string(bytes))
	dash := timelimeResponse{}
	json.Unmarshal(bytes, &dash)
	helper.Log(dash.Response.Timeline.Links.Next.Href)
	return dash.Response.Timeline.Elements
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

	u, _ := url.Parse(fmt.Sprintf("https://api.tumblr.com/v2/blog/%s.tumblr.com/posts?notes_info=true", blogName))

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
