package npf

import (
	"bytes"
	"strconv"

	"github.com/mattn/go-runewidth"
)

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

func (c *Content) RenderWithData() struct {
	contentType string
	str         string
} {
	var str bytes.Buffer
	var cType = ""

	switch c.Type {
	case "image":
		alt := c.Alt_text
		if runewidth.StringWidth(alt) == 0 {
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

		case "heading3":
			str.WriteString("③ " + c.Text)
			cType = "Heading3"

		case "quote":
			str.WriteString("> " + c.Text)
			cType = "Quote"

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

	postStr := RenderUnicode(str.String())

	return struct {
		contentType string
		str         string
	}{
		contentType: cType,
		str:         postStr,
	}
}
func (c *Content) Render() string {
	var str bytes.Buffer

	switch c.Type {
	case "image":
		alt := c.Alt_text
		if runewidth.StringWidth(alt) == 0 {
			alt = "No alt"
		}
		str.WriteString("[Image : " + alt + "]")
	case "text":
		switch c.Subtype {

		case "heading1":
			str.WriteString("# " + c.Text)

		case "heading2":
			str.WriteString("## " + c.Text)

		case "ordered-list-item":
			str.WriteString(strconv.Itoa(orderedListIndex) + ". " + c.Text)
			orderedListIndex = orderedListIndex + 1

		case "unordered-list-item":
			str.WriteString("- " + c.Text)

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

	default:
		str.WriteString(c.Text)
	}

	return str.String()
}
