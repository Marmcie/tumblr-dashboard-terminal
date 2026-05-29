package npf

import (
	"bytes"
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
	var str bytes.Buffer
	var cType = ""
	var links []string

	switch c.Type {
	case "image":
		alt := c.Alt_text

		if uniseg.StringWidth(alt) == 0 {
			alt = "No alt"
		}
		str.WriteString(fmt.Sprintf("[Image : %s]", alt))
		if len(c.Caption) > 0 {
			str.WriteString(fmt.Sprintf("\n%s", c.Caption))
		}
		cType = "Image"

	case "video":
		str.WriteString("[Video]" + "(" + c.Url + ")")
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

		str.WriteString(fmt.Sprintf("[Audio : %s By %s, From %s]", audioTitle, audioArtist, audioAlbum))
		cType = "Audio"

	case "text":
		text := c.Text
		offset := 0
		for _, f := range c.Formatting {
			switch f.Type {
			case "link":
				links = append(links, f.Url)
				t := strings.Split(text, "")
				urlString := " (" + f.Url + ")"
				text = strings.Join(t[:f.End+int64(offset)], "") + urlString + strings.Join(t[f.End+int64(offset):], "")
				offset += len(strings.Split(urlString, ""))
			}
		}

		cType = "Text"
		switch c.Subtype {

		case "heading1":
			str.WriteString("① " + text)
			cType = "Heading1"

		case "heading2":
			str.WriteString("② " + text)
			cType = "Heading2"

		case "heading3":
			str.WriteString("③ " + text)
			cType = "Heading3"

		case "quote":
			str.WriteString("> " + text)
			cType = "Quote"

		case "ordered-list-item":
			str.WriteString(strconv.Itoa(orderedListIndex) + ". ")
			str.WriteString(text)
			orderedListIndex = orderedListIndex + 1
			cType = "OrderedList"

		case "unordered-list-item":
			str.WriteString("- ")
			str.WriteString(text)
			cType = "UnOrderedList"

		default:
			str.WriteString(text)
		}

		if c.Subtype != "ordered-list-item" {
			orderedListIndex = 1
		}

	case "poll":
		str.WriteString("Poll : " + c.Question + "\n")
		for _, a := range c.Answers {
			str.WriteString("- " + a.Answer_text + "\n")
		}
		str.WriteString(c.Text)
		cType = "Poll"
	case "link":
		str.WriteString(c.Title + "(" + c.Url + ")")
		links = append(links, c.Url)

	default:
		str.WriteString(c.Text)
	}

	postStr := RenderUnicode(str.String())

	return ContentData{
		ContentType: cType,
		Str:         postStr,
		Links:       links,
	}
}
