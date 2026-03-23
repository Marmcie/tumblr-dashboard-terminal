package npf

import (
	"sort"
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
}

var orderedListIndex = 1
var renderResults map[string][]TrailData

type ContentData struct {
	ContentType string
	Str         string
}

type TrailData struct {
	Contents []ContentData
	Blog     Blog
	Layout   []Layout
	ID       int64
}

type SortPostByTimestamp []Post

func (t SortPostByTimestamp) Len() int               { return len(t) }
func (t SortPostByTimestamp) Swap(i int, j int)      { t[i], t[j] = t[j], t[i] }
func (t SortPostByTimestamp) Less(i int, j int) bool { return t[i].Timestamp > t[j].Timestamp }

type sortById []TrailData

func (t sortById) Len() int               { return len(t) }
func (t sortById) Swap(i int, j int)      { t[i], t[j] = t[j], t[i] }
func (t sortById) Less(i int, j int) bool { return t[i].ID < t[j].ID }

func (p *Post) Render() []TrailData {
	if renderResults == nil {
		renderResults = map[string][]TrailData{}
	}
	if renderResults[p.Id_string] != nil {
		return renderResults[p.Id_string]
	}
	var result []TrailData
	if len(p.Content) > 0 {
		var res []ContentData
		orderedListIndex = 1
		for _, c := range p.Content {
			data := c.RenderWithData()
			res = append(res, ContentData{
				ContentType: data.ContentType,
				Str:         data.Str,
			})
		}
		result = append(result, TrailData{
			Contents: res,
			Blog:     p.Blog,
			Layout:   p.Layout,
			ID:       p.Id,
		})
	}
	for _, t := range p.Trail {
		var res []ContentData
		orderedListIndex = 1
		for _, c := range t.Content {
			data := c.RenderWithData()
			res = append(res, ContentData{
				ContentType: data.ContentType,
				Str:         data.Str,
			})
		}
		tID, _ := strconv.ParseInt(t.Post.Id, 10, 64)
		result = append(result, TrailData{
			Contents: res,
			Blog:     t.Blog,
			Layout:   t.Layout,
			ID:       tID,
		})
	}
	sort.Sort(sortById(result))
	renderResults[p.Id_string] = result
	return renderResults[p.Id_string]
}

func (p *Post) GetSummary() string {
	return RenderUnicode(p.Summary)
}

func (p *Post) RemoveRenderResult() {
	delete(renderResults, p.Id_string)
}
