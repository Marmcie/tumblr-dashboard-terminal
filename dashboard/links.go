package dashboard

import (
	"tumblr-dt/modules"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
)

type LinkWindow struct {
	LinkList  []string
	Window    *component.Box
	List      *component.Selectlist
	prev      string
	dashboard *Dashboard
}

func NewLinkWindow(dashboard *Dashboard) *LinkWindow {
	l := &LinkWindow{}
	l.Window = component.NewBox("Link window")
	l.List = component.NewSelectlist("Link list")
	l.dashboard = dashboard

	l.Window.SetAbsolute(true).
		SetW(80).
		SetH(30).
		SetBorder(true).
		SetCentered(true)

	l.Window.SetTitle("Links")

	l.Window.SetBorderLabelColor("BottomRight", ui.GetColorStr(ui.ColorWhite))
	l.Window.SetBorderLabel("BottomRight", "Esc to close")

	l.List.SetWidthInherit(true).
		SetHeightInherit(true)

	l.List.SetSelectedOptionBackground(ui.GetColorStr(ui.ColorFocus))
	l.List.SetSelectedOptionForeground(ui.GetColorStr(ui.ColorWhite))
	l.List.SetBorderFocusForeground(ui.GetColorStr(ui.ColorFocusBorder))

	l.Window.AddChild(l.List)

	l.Window.SetVisibility(false)

	l.List.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {

		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter":
				if l.List.Cursor < len(l.List.GetChildren()) {
					l.List.RunSelectedOption()
				}
			case "j":
				l.List.IncrementCursor()
			case "k":
				l.List.DecrementCursor()
			case "G":
				l.List.SetCursor(len(l.LinkList) - 1)
				l.List.RunSelectedOption()
			case "g":
				if l.prev == "g" {
					l.List.SetCursor(0)
					l.List.RunSelectedOption()
				}
			case "esc":
				l.dashboard.toggleLinkWindow()
			}
			l.prev = msg.String()
		}
	}, true)

	return l
}

func (l *LinkWindow) SetLinks(links []string) {
	l.LinkList = links
	l.UpdateLinks()
	if len(links) > 0 {
		l.List.SetCursor(0)
	}
}

func (l *LinkWindow) UpdateLinks() {

	l.List.ClearChildren()
	for _, link := range l.LinkList {
		box := component.NewBox("Link item")

		box.SetBorder(true).
			SetBorders(false, true, false, false).
			SetBorderCorner(false).
			SetH(2).
			SetWidthInherit(true)
		title := component.NewLine("Link Title")
		title.SetWidthInherit(true).SetH(1)
		title.SetText(link)
		box.AddChild(title)

		l.List.AddOption(box, func() {
			modules.OpenInBrowser(link)
			component.Global.SetCmd(tea.ClearScreen)
		})
	}

}

func (l *LinkWindow) Focus() {
	l.Window.SetVisibility(true)
	l.List.Focus()
}

func (l *LinkWindow) Blur() {
	l.Window.SetVisibility(false)
}
