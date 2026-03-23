package modules

import (
	"bytes"
	"strconv"

	"github.com/forPelevin/gomoji"
	"github.com/mattn/go-runewidth"
	"golang.org/x/text/width"
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
}

type Blog struct {
	Name                   string
	Title                  string
	Description            string
	Url                    string
	Uuid                   string
	Updated                int64
	Avatar                 []Avatar
	Tumblrmart_accessories struct {
		Badges               []Badge
		Blue_checkmark_count int64
	}
	Can_show_badges    bool
	Active             bool
	Show_follow_action bool
}

type Content struct {
	Type                        string
	Text                        string
	Width                       int64
	Height                      int64
	Original_dimensions_missing bool
	Cropped                     bool
	Has_original_dimensions     bool
	Subtype                     string
	Indent_level                int64
	Formatting                  []Formatting
	Alt_text                    string
	Caption                     string
	Media                       []Media
	Feedback_token              string
	Colors                      map[string]string
	Poster                      Media
	Author                      string

	Provider   string
	Artist     string
	Album      string
	Embed_html string
	Embed_url  string

	Can_autoplay_on_cellular bool

	Is_visible bool

	Metadata struct {
		Id string
	}
	Question string
	Answers  []struct {
		Client_id   string
		Answer_text string
	}
}

type Formatting struct {
	Start int64
	End   int64
	Type  string
	Url   string
	Blog  struct {
		Uuid string
		Name string
		Url  string
	}

	Hex string
}

type Media struct {
	Type   string
	Url    string
	Width  int64
	Height int64
	Hd     bool
}

type Layout struct {
	Type           string
	Display        []map[string][]int64
	Truncate_after int64
	Blocks         []int64

	Attribution struct {
		Type string
		Url  string
		Blog Blog
	}
}

type TrailPost struct {
	Post struct {
		Id            int64
		Timestamp     int64
		Is_commercial bool
	}
	Blog             Blog
	Content          []Content
	Layout           []Layout
	Broken_blog_name string
}

type Avatar struct {
	Width  int64
	Height int64
	Url    string
}

type Badge struct {
	Product_group   string
	Urls            []string
	Destination_url string
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

func (p *Post) BlogNames() []string {
	var result []string
	if len(p.Content) > 0 {
		result = append(result, p.Blog.GetName())
	}
	for _, t := range p.Trail {
		name := t.Blog.GetName()
		if len(name) == 0 {
			name = t.Broken_blog_name
		}
		result = append(result, name)
	}
	return result
}

func (p *Post) RenderString() string {
	var str = ""
	var reblog = ">"
	if len(p.Content) > 0 {
		str = "[" + p.Blog.Name + "]\n"
		for _, c := range p.Content {
			str += c.Render()
			str += "\n"
		}
	}
	for _, t := range p.Trail {
		str += reblog + "[" + t.Blog.Name + "]\n"
		for _, c := range t.Content {
			str += c.Render()
			str += "\n"
		}
		str += "\n"
		reblog += ">"
	}

	return str
}

func (c *Content) RenderWithData() struct {
	contentType string
	str         string
} {
	var str bytes.Buffer
	var cType = ""

	switch c.Type {
	case "image":
		alt := c.Alt_text
		if len(alt) == 0 {
			alt = "No alt"
		}
		str.WriteString("[Image : " + alt + "]")
		cType = "Image"
	case "text":

		cType = "Text"
		switch c.Subtype {

		case "heading1":
			str.WriteString("① " + c.Text)
			cType = "Heading1"

		case "heading2":
			str.WriteString("② " + c.Text)
			cType = "Heading2"

		case "ordered-list-item":
			str.WriteString(strconv.Itoa(orderedListIndex) + ". ")
			str.WriteString(c.Text)
			orderedListIndex = orderedListIndex + 1
			cType = "OrderedList"

		case "unordered-list-item":
			str.WriteString("- ")
			str.WriteString(c.Text)
			cType = "UnOrderedList"

		default:
			str.WriteString(c.Text)
		}

		if c.Subtype != "ordered-list-item" {
			orderedListIndex = 1
		}

	case "poll":
		str.WriteString("Question : " + c.Question + "\n")
		for _, a := range c.Answers {
			str.WriteString("- " + a.Answer_text + "\n")
		}
		str.WriteString(c.Text)
		cType = "Poll"

	default:
		str.WriteString(c.Text)
	}

	var result bytes.Buffer

	con := runewidth.Condition{
		EastAsianWidth:     false,
		StrictEmojiNeutral: false,
	}
	
	postStr:=gomoji.RemoveEmojis(str.String())
	for _, v := range postStr {
		info := width.LookupRune(v)
		runeWidth := con.RuneWidth(v)

		if info.Kind() == width.EastAsianFullwidth || info.Kind() == width.EastAsianWide {
			for range runeWidth - 1 {
				// INFO: Output 0 width character to account for full width chars
				result.WriteRune('\u200b')
				// result.WriteString(strconv.Itoa(runeWidth))
				// result.WriteRune('#')
			}
		}
		result.WriteRune(v)
	}

	return struct {
		contentType string
		str         string
	}{
		contentType: cType,
		str:         result.String(),
	}
}

func (p *Post) GetSummary() string {
	con := runewidth.Condition{
		EastAsianWidth:     true,
		StrictEmojiNeutral: false,
	}
	var result bytes.Buffer
	summary:=gomoji.RemoveEmojis(p.Summary)
	for _, v := range summary {
		info := width.LookupRune(v)
		runeWidth := con.RuneWidth(v)

		//INFO: SUS EMOJIS LIST : ⭐⭐⭐⭐
		if info.Kind() == width.EastAsianFullwidth || info.Kind() == width.EastAsianWide {
			for range runeWidth - 1 {
				// INFO: Output 0 width character to account for full width chars
				result.WriteRune('\u200b')
				// result.WriteRune('#')
			}
		}
		result.WriteRune(v)
	}
	return result.String()
}

func (c *Content) Render() string {
	var str = ""

	switch c.Type {
	case "image":
		alt := c.Alt_text
		if len(alt) == 0 {
			alt = "No alt"
		}
		str += "[Image : " + alt + "]"
	case "text":
		switch c.Subtype {

		case "heading1":
			str += "# " + c.Text

		case "heading2":
			str += "## " + c.Text

		case "ordered-list-item":
			str += strconv.Itoa(orderedListIndex) + ". "
			str += c.Text
			orderedListIndex = orderedListIndex + 1

		case "unordered-list-item":
			str += "- "
			str += c.Text

		default:
			str += c.Text
		}

		if c.Subtype != "ordered-list-item" {
			orderedListIndex = 1
		}

	case "poll":
		str += "Question : " + c.Question + "\n"
		for _, a := range c.Answers {
			str += "- " + a.Answer_text + "\n"
		}
		str += c.Text

	default:
		str += c.Text
	}

	return str
}

func (m *Media) Render() string {
	var str = ""
	str += "![Image]("
	str += m.Url
	str += ")"
	return str
}

func (b *Blog) GetName() string {
	if len(b.Name) > 0 {
		return b.Name
	}

	return b.Title
}
