package npf

import (
	"fmt"
	"slices"
	"strconv"

	mapset "github.com/deckarep/golang-set/v2"
)

type Post struct {
	Type                       string
	Original_type              string
	Is_blocks_post_format      bool
	Blog_name                  string
	Blog                       Blog
	Id                         int64
	Id_string                  string
	Is_blazed                  bool
	Is_blaze_pending           bool
	Can_ignite                 bool
	Can_blaze                  bool
	Post_url                   string
	Parent_post_url            string
	Slug                       string
	Date                       string
	Timestamp                  int64
	State                      string
	Reblog_key                 string
	Tags                       []string
	Short_url                  string
	Summary                    string
	Should_open_in_legacy      bool
	Recommended_source         string
	Recommended_color          string
	Followed                   bool
	Liked                      bool
	Note_count                 int64
	Content                    []Content
	Layout                     []Layout
	Trail                      []TrailPost
	Reblogged_from_id          int64
	Reblogged_from_url         string
	Reblogged_from_name        string
	Reblogged_from_title       string
	Reblogged_from_uuid        string
	Reblogged_from_can_message bool
	Reblogged_from_following   bool
	Reblogged_root_id          int64
	Reblogged_root_url         string
	Reblogged_root_name        string
	Reblogged_root_title       string
	Reblogged_root_uuid        string
	Reblogged_root_can_message bool
	Reblogged_root_following   bool
	Can_like                   bool
	Interactability_reblog     string
	Can_reblog                 bool
	Interactability_blaze      string
	Can_send_in_message        bool
	Can_reply                  bool
	Display_avatar             bool
	Rendered                   bool
	Result                     []TrailData
	IsFiltered                 bool
	FilteredContents           mapset.Set[string]
	FilteredTags               mapset.Set[string]

	// Keeps track of which link has been stored already.
	// Without it multiple links to same blog in reblog chain can end up in a list.
	LinkSet mapset.Set[string]
	// Link URLs
	Links []string
	// Title for each of the links to be displayed in links modal
	LinkTitles []string
}

var orderedListIndex = 1
var renderResults map[string][]TrailData

type ContentData struct {
	ContentType string
	Str         string
	Links       []string
	LinkTitles  []string
}

type TrailData struct {
	Contents []ContentData
	Blog     Blog
	BlogName string
	Layout   []Layout
	ID       int64
}

// Generic post object to unify post object and posts in trail
type GenericPost struct {
	Contents  []Content
	Timestamp int64
	Blog      Blog
	BlogName  string
	Layout    []Layout
	ID        int64
}

func (p *Post) Render() []TrailData {
	if renderResults == nil {
		renderResults = map[string][]TrailData{}
	}
	if renderResults[p.Id_string] != nil {
		return renderResults[p.Id_string]
	}
	var result []TrailData
	p.LinkSet = mapset.NewSet[string]()

	p.LinkSet.Add(p.Blog.Url)
	p.Links = append(p.Links, p.Blog.Url)
	p.LinkTitles = append(p.LinkTitles, fmt.Sprintf("%s's blog", p.Blog.Name))

	// Put all posts and trails into a slice to sort them correctly.
	// Needed because sometimes retrieved posts are ordered incorrectly and can lead to incorrectly ordered media link name in post link modal.
	posts := []GenericPost{}

	posts = append(posts, GenericPost{
		Contents:  p.Content,
		Timestamp: p.Timestamp,
		Blog:      p.Blog,
		BlogName:  p.Blog_name,
		Layout:    p.Layout,
		ID:        p.Id,
	})

	for _, trail := range p.Trail {
		blogName := trail.Blog.Name
		if len(blogName) == 0 {
			blogName = trail.Broken_blog_name
		}

		tID, _ := strconv.ParseInt(trail.Post.Id, 10, 64)
		posts = append(posts, GenericPost{
			Contents:  trail.Content,
			Timestamp: trail.Post.Timestamp,
			Blog:      trail.Blog,
			BlogName:  blogName,
			Layout:    trail.Layout,
			ID:        tID,
		})
	}

	slices.SortFunc(posts, func(a GenericPost, b GenericPost) int { return int(a.ID) - int(b.ID) })

	/** To differentiate images in post link list */
	imageCount := 1
	videoCount := 1
	audioCount := 1

	for _, post := range posts {
		if len(post.Contents) == 0 {
			continue
		}
		var res []ContentData

		blogName := post.BlogName
		if !p.LinkSet.Contains(post.Blog.Url) {
			p.LinkSet.Add(post.Blog.Url)
			p.Links = append(p.Links, post.Blog.Url)
			p.LinkTitles = append(p.LinkTitles, fmt.Sprintf("%s's blog", blogName))
		}

		orderedListIndex = 1
		for _, c := range post.Contents {
			data := c.RenderWithData(imageCount, videoCount, audioCount)
			res = append(res, ContentData{
				ContentType: data.ContentType,
				Str:         data.Str,
			})
			switch c.Type {
			case "image":
				imageCount++
			case "video":
				videoCount++
			case "audio":
				audioCount++
			}

			for i := range len(data.Links) {
				if !p.LinkSet.Contains(data.Links[i]) {
					p.LinkSet.Add(data.Links[i])
					p.Links = append(p.Links, data.Links[i])
					p.LinkTitles = append(p.LinkTitles, data.LinkTitles[i])
				}
			}
		}
		result = append(result, TrailData{
			Contents: res,
			Blog:     post.Blog,
			BlogName: blogName,
			Layout:   post.Layout,
			ID:       post.ID,
		})
	}

	renderResults[p.Id_string] = result
	return renderResults[p.Id_string]
}

func (p *Post) GetSummary() string {
	return RenderUnicode(p.Summary)
}

func (p *Post) RemoveRenderResult() {
	delete(renderResults, p.Id_string)
}

func (p *Post) GetLinks() []string {
	if p.Links == nil {
		p.Links = []string{}
	}
	return p.Links
}
func (p *Post) GetLinkTitles() []string {
	if p.LinkTitles == nil {
		p.LinkTitles = []string{}
	}
	return p.LinkTitles
}
