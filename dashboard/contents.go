package dashboard

import (
	"bytes"
	"fmt"
	"strings"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	"github.com/rivo/uniseg"
)

type Contents struct {
	contentElem *component.Scrollable
	dashboard   *Dashboard
	prev        string
}

func NewContents(dashboard *Dashboard) *Contents {
	f := &Contents{}
	f.contentElem = component.NewScrollable("Contents")
	f.contentElem.SetBorder(true).SetWidthInherit(true)
	f.contentElem.SetForeground(ui.GetColorStr(ui.ColorWhite))
	f.contentElem.SetBorderFocusForeground(ui.GetColorStr(ui.ColorFocusBorder))
	f.dashboard = dashboard

	f.InitEvents()

	return f
}
func (f *Contents) InitEvents() {
	f.contentElem.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case f.dashboard.config.Keymaps.Navigation.Left:
				f.dashboard.FocusFeed()
			case f.dashboard.config.Keymaps.Navigation.Down:
				f.contentElem.OffsetY = min(f.contentElem.Bottom-1, f.contentElem.OffsetY+1)

			case f.dashboard.config.Keymaps.Navigation.Up:
				f.contentElem.OffsetY = max(0, f.contentElem.OffsetY-1)

			case strings.ToUpper(f.dashboard.config.Keymaps.Navigation.Down):
				f.contentElem.OffsetY = min(f.contentElem.Bottom-1, f.contentElem.OffsetY+3)

			case strings.ToUpper(f.dashboard.config.Keymaps.Navigation.Up):
				f.contentElem.OffsetY = max(0, f.contentElem.OffsetY-3)

			case f.dashboard.config.Keymaps.Navigation.JumpNext:
				pos := 0
				for _, c := range f.contentElem.Children {
					if pos > f.contentElem.OffsetY {
						f.contentElem.OffsetY = min(f.contentElem.Bottom-1, pos)
						break
					}
					pos += c.GetHeight()
				}
			case f.dashboard.config.Keymaps.Navigation.JumpPrev:
				pos := 0
				for i, c := range f.contentElem.Children {
					if pos >= f.contentElem.OffsetY {
						f.contentElem.OffsetY = max(0, pos-f.contentElem.Children[max(0, i-1)].GetHeight())
						return
					}
					pos += c.GetHeight()
				}
				f.contentElem.OffsetY = max(0, pos-f.contentElem.Children[len(f.contentElem.Children)-1].GetHeight())

			case f.dashboard.config.Keymaps.Navigation.JumpBottom:
				f.contentElem.OffsetY = f.contentElem.Bottom - 1

			case f.dashboard.config.Keymaps.Navigation.JumpTop:
				if f.prev == f.dashboard.config.Keymaps.Navigation.JumpTop {
					f.contentElem.OffsetY = 0
				}
			}
			f.prev = msg.String()
		}
	}, true)
}

func (f *Contents) DisplayPost(post *npf.Post, showFiltered bool) {
	isSlimMode := f.dashboard.config.Post_theme == "slim"

	f.contentElem.ClearChildren()

	if post.IsFiltered && !showFiltered {
		box := component.NewBox("Post")
		box.SetBorder(true)
		box.SetWidthInherit(true)
		box.SetH(6)
		f.contentElem.AddChild(box)
		l := component.NewText("Filtered content")
		l.SetText("This post contains filtered content.\nPress enter to read.")
		l.SetWidthInherit(true)
		l.SetHeightInherit(true)
		box.AddChild(l)
		return
	}

	for _, reblog := range post.Render() {
		box := component.NewBox("Post")
		box.SetBorder(true).SetWidthInherit(true)
		if isSlimMode {
			box.SetBorders(true, false, false, false).SetBorderCorner(false)
		}
		f.contentElem.AddChild(box)
		innerWidth := box.GetInnerWidth()
		var str bytes.Buffer

		//INFO: Array of each lines
		parts := []string{}
		//INFO: Array of styles for each lines
		colors := []string{}
		col := " "

		var askLayout npf.Layout
		var askStart = -100
		var askEnd = -100
		for _, layout := range reblog.Layout {
			if layout.Type == "ask" {
				askLayout = layout
				askStart = int(layout.Blocks[0])
				askEnd = int(layout.Blocks[len(layout.Blocks)-1])
				break
			}
		}
		isAsk := false
		// INFO: Loop through rendered NPF content blocks
		for i, contents := range reblog.Contents {
			contentStr := contents.Str
			if i == askStart {
				isAsk = true
				blogName := "Anonymous"
				if askLayout.Attribution != nil {
					blogName = askLayout.Attribution.Blog.GetName()
				}
				contentStr = fmt.Sprintf("%s asked :\n %s", blogName, contentStr)
			}

			if i == askEnd+1 {
				isAsk = false
				contentStr = fmt.Sprintf("%s answered :\n %s", reblog.Blog.GetName(), contentStr)
			}

			contentType := contents.ContentType
			//INFO: Change text color based on content type
			switch contentType {
			case "Heading1":
				col = ui.GetColorStr(ui.ColorH1)
			case "Image", "Video", "Audio":
				col = ui.GetColorStr(ui.ColorImage)
			case "Heading2":
				col = ui.GetColorStr(ui.ColorH2)
			case "Quote":
				col = ui.GetColorStr(ui.ColorQuote)
			case "Poll":
				col = ui.GetColorStr(ui.ColorQuote)
			default:
				col = ""
			}
			if isAsk {
				col = ui.GetColorStr(ui.ColorQuote)
			}

			for line := range strings.SplitSeq(contentStr, "\n") {

				state := -1
				var word string
				for len(line) > 0 {
					word, line, state = uniseg.FirstWordInString(line, state)
					//INFO: Divide the text into lines, while preventing word break
					// for word := range strings.SplitSeq(contentStr, " ") {
					// word = strings.Trim(word, " ")
					strString := str.String()
					if uniseg.StringWidth(strString)+uniseg.StringWidth(word)+1 >= innerWidth {
						parts = append(parts, strString)
						colors = append(colors, col)
						//INFO: If the single word is wider than the box,
						//or the language doesn't use white space as separator,
						//split the word into smaller chunks
						if innerWidth > 1 && uniseg.StringWidth(word) >= innerWidth {
							w := strings.ReplaceAll(word, " ", "")
							//INFO: Loop through each characters to determine real width of string split
							for uniseg.StringWidth(w) >= innerWidth {
								l := uniseg.StringWidth(w)
								parts = append(parts, w[:l])
								colors = append(colors, col)
								w = w[l:]
							}
							parts = append(parts, w)
							colors = append(colors, col)
							str.Reset()
						} else {
							str.Reset()
							str.WriteString(word)
						}
					} else {
						str.WriteString(word)
					}
				}

				parts = append(parts, str.String())
				colors = append(colors, col)
				str.Reset()
			}

			if len(strings.Trim(str.String(), " ")) > 0 {
				parts = append(parts, strings.Trim(str.String(), " "))
				colors = append(colors, col)
				str.Reset()
			}

			parts = append(parts, "")
			colors = append(colors, "")
		}
		colors = append(colors, col)
		parts = append(parts, str.String())

		labelAlignment := "Top"
		if isSlimMode {
			labelAlignment = "TopLeft"
		}

		if f.dashboard.config.Use_blog_avatar_color {
			box.SetBorderLabelColor(labelAlignment, reblog.Blog.GetBlogColor())
		}
		box.SetBorderLabel(labelAlignment, reblog.BlogName)
		box.SetH(max(3, len(parts)))

		//INFO: Convert each line into Line object, then apply corresponding style
		for i := 0; i < min(len(parts), box.GetInnerHeight()); i++ {
			line := strings.Trim(parts[i], " ")
			col := colors[i]
			l := component.NewLine("Post text")
			l.SetText(line)
			if col != "" {
				l.SetForeground(col)
			}

			l.SetWidthInherit(true)
			box.AddChild(l)
		}
	}

}

func (f *Contents) Focus() {
	f.contentElem.Focus()
}
