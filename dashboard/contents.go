package dashboard

import (
	"strings"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Contents struct {
	contentElem *component.Selectlist
	dashboard   *Dashboard
}

func NewContents(dashboard *Dashboard) *Contents {
	f := &Contents{}
	f.contentElem = component.NewSelectlist("Contents")
	f.contentElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetHeightInherit(true).SetWidthInherit(true)
	f.dashboard = dashboard

	f.contentElem.SelectBgStyle = lipgloss.NewStyle()
	f.InitEvents()

	return f
}
func (f *Contents) InitEvents() {
	f.contentElem.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "h":
				f.dashboard.FocusFeed()
			}
		}
	})

}

func (f *Contents) DisplayPost(post modules.Post) {
	f.contentElem.ClearChildren()
	for _, reblog := range post.Render() {

		box := component.NewBox("Post")
		box.SetBorder(true).SetBorderPadding(1)
		box.SetWidthInherit(true)
		f.contentElem.AddOption(box, func() {})
		innerWidth := box.GetInnerWidth()
		str := ""
		style := lipgloss.NewStyle()

		parts := []string{}
		styles := []lipgloss.Style{}
		for _, contents := range reblog {
			contentType := contents.ContentType
			switch contentType {
			case "Heading1":
				style = style.Foreground(lipgloss.Color("#40f0f0"))

			case "Image":
				style = style.Foreground(lipgloss.Color("#40f0f0"))
			case "Heading2":
				style = style.Foreground(lipgloss.Color("#4000f0"))
			}

			for lines := range strings.SplitSeq(contents.Str, "\n") {
				for word := range strings.SplitSeq(lines, " ") {
					if len(str)+len(word)+1 >= innerWidth {
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

		for i, line := range parts {
			style := styles[i]
			l := component.NewLine("Post text")
			l.SetText(line)
			l.SetStyle(style)
			l.SetWidthInherit(true)
			box.AddChild(l)
		}
		box.SetH(len(parts) + 1)
	}

}

func (f *Contents) Focus() {
	f.contentElem.Focus()
}
