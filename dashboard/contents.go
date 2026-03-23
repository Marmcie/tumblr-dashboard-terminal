package dashboard

import (
	"strings"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-runewidth"
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
	f.dashboard = dashboard

	f.InitEvents()

	return f
}
func (f *Contents) InitEvents() {
	f.contentElem.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "h":
				f.dashboard.FocusFeed()
			case "j":
				f.contentElem.OffsetY = min(f.contentElem.Bottom-1, f.contentElem.OffsetY+1)
			case "k":
				f.contentElem.OffsetY = max(0, f.contentElem.OffsetY-1)

			case "G":
				f.contentElem.OffsetY = f.contentElem.Bottom - 1

			case "g":
				if f.prev == "g" {
					f.contentElem.OffsetY = 0
				}
			}
			f.prev = msg.String()
		}
	}, true)
}

func (f *Contents) DisplayPost(post npf.Post) {
	f.contentElem.ClearChildren()
	for _, reblog := range post.Render() {
		box := component.NewBox("Post")
		box.SetBorder(true)
		box.SetWidthInherit(true)
		f.contentElem.AddChild(box)
		innerWidth := box.GetInnerWidth()
		str := ""

		//INFO: Array of each lines
		parts := []string{}
		//INFO: Array of styles for each lines
		colors := []string{}
		col := " "

		// INFO: Loop through rendered NPF content blocks
		for _, contents := range reblog.Contents {
			contentType := contents.ContentType
			//INFO: Change text color based on content type
			switch contentType {
			case "Heading1":
				col = ui.GetColorStr(ui.ColorH1)

			case "Image":
				col = ui.GetColorStr(ui.ColorImage)

			case "Video":
				col = ui.GetColorStr(ui.ColorImage)
			case "Heading2":
				col = ui.GetColorStr(ui.ColorH2)
			case "Quote":
				col = ui.GetColorStr(ui.ColorQuote)
			default:
				col = ""
			}

			//INFO: Divide the text into lines, while preventing word break
			for lines := range strings.SplitSeq(contents.Str, "\n") {
				for word := range strings.SplitSeq(lines, " ") {
					if runewidth.StringWidth(str)+runewidth.StringWidth(word)+1 >= innerWidth {
						parts = append(parts, str)
						colors = append(colors, col)
						//INFO: If the single word is wider than the box,
						//or the language doesn't use white space as separator,
						//split the word into smaller chunks
						if innerWidth > 1 && runewidth.StringWidth(word) >= innerWidth {
							w := word
							//INFO: Loop through each characters to determine real width of string split
							for runewidth.StringWidth(w) >= innerWidth {
								l := 0
								for i := 0; l < innerWidth && i < len(w); i++ {
									l += runewidth.StringWidth(string(w[i]))
								}
								parts = append(parts, w[:l])
								colors = append(colors, col)
								w = w[l:]
							}
							parts = append(parts, w)
							colors = append(colors, col)
							str = ""
						} else {
							str = word + " "
						}
					} else {
						str += word + " "
					}
				}
			}
			parts = append(parts, str)
			colors = append(colors, col)
			str = ""

			parts = append(parts, "")
			colors = append(colors, "")
		}
		colors = append(colors, col)
		parts = append(parts, str)

		top, _, _, _ := box.GetBorderPaddings()

		box.SetTitle(reblog.Blog.Name)
		box.SetH(max(3, len(parts)+1))

		//INFO: Convert each line into Line object, then apply corresponding style
		for i := 0; i < min(len(parts), box.GetInnerHeight()-top); i++ {
			line := parts[i]
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
