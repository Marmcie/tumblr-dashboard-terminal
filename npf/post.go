package npf


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
}









var orderedListIndex = 1

type ContentData struct {
	ContentType string
	Str         string
}

type TrailData struct {
	Contents []ContentData
	Blog     Blog
}

func (p *Post) Render() []TrailData {
	var result []TrailData
	if len(p.Content) > 0 {
		var res []ContentData
		orderedListIndex = 1
		for _, c := range p.Content {
			data := c.RenderWithData()
			res = append(res, ContentData{
				ContentType: data.contentType,
				Str:         data.str,
			})
		}
		result = append(result, TrailData{
			Contents: res,
			Blog:     p.Blog,
		})
	}
	for _, t := range p.Trail {
		var res []ContentData
		orderedListIndex = 1
		for _, c := range t.Content {
			data := c.RenderWithData()
			res = append(res, ContentData{
				ContentType: data.contentType,
				Str:         data.str,
			})
		}
		result = append(result, TrailData{
			Contents: res,
			Blog:     t.Blog,
		})
	}

	return result
}



func (p *Post) GetSummary() string {
	return RenderUnicode(p.Summary)
}






