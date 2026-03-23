package dashboard

import (
	"strings"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	f.contentElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetWidthInherit(true)
	f.dashboard = dashboard

	f.contentElem.AddEventListener("onFocusChange", func(m tea.Msg, i int) {
		if f.contentElem.GetFocusState() {
			f.contentElem.SetStyle(lipgloss.NewStyle().Foreground(ui.GetColor("ColorWhite")))
		} else {
			f.contentElem.ClearStyle()
		}
	})
	// f.contentElem.SelectBgStyle = lipgloss.NewStyle()
	f.InitEvents()

	return f
}
func (f *Contents) InitEvents() {
	f.contentElem.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
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
	})
}

func (f *Contents) DisplayPost(post npf.Post) {
	f.contentElem.ClearChildren()
	for _, reblog := range post.Render() {
		box := component.NewBox("Post")
		box.SetBorder(true).SetBorderPadding(1)
		box.SetWidthInherit(true)
		f.contentElem.AddChild(box)
		innerWidth := box.GetInnerWidth()
		str := ""
		style := lipgloss.NewStyle()

		//INFO: Array of each lines
		parts := []string{}
		//INFO: Array of styles for each lines
		styles := []lipgloss.Style{}

		// INFO: Loop through rendered NPF content blocks
		for _, contents := range reblog.Contents {
			contentType := contents.ContentType
			//INFO: Change text color based on content type
			switch contentType {
			case "Heading1":
				style = style.Foreground(ui.GetColor("ColorH1"))
			case "Image":
				style = style.Foreground(ui.GetColor("ColorImage"))
			case "Heading2":
				style = style.Foreground(ui.GetColor("ColorH2"))
			case "Quote":
				style = style.Foreground(ui.GetColor("ColorQuote"))
			default:
				style = lipgloss.NewStyle()
			}

			//INFO: Divide the text into lines, while preventing word break
			for lines := range strings.SplitSeq(contents.Str, "\n") {
				for word := range strings.SplitSeq(lines, " ") {
					if runewidth.StringWidth(str)+runewidth.StringWidth(word)+1 >= innerWidth {
						parts = append(parts, str)
						styles = append(styles, style)
						str = word + " "
					} else {
						str += word + " "
					}
				}
			}
			parts = append(parts, str)
			styles = append(styles, style)
			str = ""
		}
		styles = append(styles, style)
		parts = append(parts, str)

		//INFO: Convert each line into Line object, then apply corresponding style
		for i, line := range parts {
			style := styles[i]
			l := component.NewLine("Post text")
			l.SetText(line)
			l.SetStyle(style)
			l.SetWidthInherit(true)
			box.AddChild(l)
		}
		box.SetTitle(reblog.Blog.Name)
		box.SetH(len(parts) + 1)
	}

}

func (f *Contents) Focus() {
	f.contentElem.Focus()
}
