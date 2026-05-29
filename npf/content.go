package npf

import (
	"fmt"
	"github.com/rivo/uniseg"
	"strconv"
	"strings"
)

type Content struct {
	Type                        string
	Text                        string
	Title                       string
	Link                        string
	Url                         string
	Uisplay_url                 string
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

func (c *Content) RenderWithData() ContentData {
	var b strings.Builder
	var cType = ""
	var links []string

	switch c.Type {
	case "image":
		alt := c.Alt_text

		if uniseg.StringWidth(alt) == 0 {
			alt = "No alt"
		}
		fmt.Fprintf(&b, "[Image : %s]", alt)
		if len(c.Caption) > 0 {
			fmt.Fprintf(&b, "\n%s", c.Caption)
		}
		cType = "Image"

	case "video":
		fmt.Fprintf(&b, "[Video](%s)", c.Url)
		cType = "Video"

	case "audio":
		audioTitle := c.Title
		audioArtist := c.Artist
		audioAlbum := c.Album
		if len(audioTitle) == 0 {
			audioTitle = "Unknown audio"
		}

		if len(audioArtist) == 0 {
			audioArtist = "Unknown artist"
		}

		if len(audioAlbum) == 0 {
			audioAlbum = "Unknown album"
		}
		fmt.Fprintf(&b, "[Audio : %s By %s, From %s]", audioTitle, audioArtist, audioAlbum)
		cType = "Audio"

	case "text":
		text := c.Text
		offset := 0
		for _, f := range c.Formatting {
			switch f.Type {
			case "link":
				links = append(links, f.Url)
				t := strings.Split(text, "")
				
				urlString := fmt.Sprintf(" (%s)",f.Url)
				text = strings.Join(t[:f.End+int64(offset)], "") + urlString + strings.Join(t[f.End+int64(offset):], "")
				offset += len(strings.Split(urlString, ""))
			}
		}

		cType = "Text"
		switch c.Subtype {

		case "heading1":
			fmt.Fprintf(&b, "① %s", text)
			cType = "Heading1"

		case "heading2":
			fmt.Fprintf(&b, "② %s", text)
			cType = "Heading2"

		case "heading3":
			fmt.Fprintf(&b, "③ %s", text)
			cType = "Heading3"

		case "quote":
			fmt.Fprintf(&b, "> %s", text)
			cType = "Quote"

		case "ordered-list-item":
			fmt.Fprintf(&b, "%s. %s", strconv.Itoa(orderedListIndex), text)
			orderedListIndex = orderedListIndex + 1
			cType = "OrderedList"

		case "unordered-list-item":
			fmt.Fprintf(&b, "- %s", text)
			cType = "UnOrderedList"

		default:
			fmt.Fprintf(&b, "%s", text)
		}

		if c.Subtype != "ordered-list-item" {
			orderedListIndex = 1
		}

	case "poll":
		fmt.Fprintf(&b, "Poll : %s\n", c.Question)
		for _, a := range c.Answers {
			fmt.Fprintf(&b, "- %s\n", a.Answer_text)
		}
		fmt.Fprintf(&b, "%s", c.Text)
		cType = "Poll"
	case "link":
		fmt.Fprintf(&b, "%s(%s)", c.Title, c.Url)
		links = append(links, c.Url)

	default:
		fmt.Fprintf(&b, "%s", c.Text)
	}

	postStr := RenderUnicode(b.String())

	return ContentData{
		ContentType: cType,
		Str:         postStr,
		Links:       links,
	}
}
