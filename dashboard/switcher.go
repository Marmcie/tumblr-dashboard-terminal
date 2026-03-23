package dashboard

import (
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
)

type Switcher struct {
	Window     *component.Flex
	DashOption *component.Box
	DashLabel  *component.Line
	TagOption  *component.Flex
	TagInput   *component.Line
	TagLabel   *component.Line
	BlogOption *component.Box
	dashboard  *Dashboard
}

func NewSwitcher(dashboard *Dashboard) *Switcher {
	s := &Switcher{}
	s.dashboard = dashboard
	s.Window = component.NewFlex("Switcher window")
	s.Window.
		SetAbsolute(true).
		SetCentered(true).
		SetSize(30, 7).
		SetBorder(true).
		SetBorderPadding(2)
	s.Window.SetTitle("Feed picker")
	s.Window.SetBorderLabel("BottomRight", "Esc to close")

	s.DashOption = component.NewBox("Dash option")
	s.DashOption.SetBorder(true).
		SetWidthInherit(true).
		SetBorders(false, false, true, false).
		SetBorderCorner(false)

	s.DashLabel = component.NewLine("Dash label")
	s.DashLabel.SetText("Dashboard")

	s.DashOption.AddChild(s.DashLabel)

	s.BlogOption = component.NewBox("Blog option")
	s.BlogOption.SetBorder(true).
		SetWidthInherit(true).
		SetBorders(false, false, true, false).
		SetBorderCorner(false)

	BlogLabel := component.NewLine("Dash label")
	BlogLabel.SetText("Blog")
	s.BlogOption.AddChild(BlogLabel)

	s.TagLabel = component.NewLine("Tag label")
	s.TagLabel.SetText("Tag : ")

	s.TagOption = component.NewFlex("Tag option")
	s.TagOption.Direction = 1
	s.TagOption.SetBorder(true).
		SetWidthInherit(true).
		SetBorders(false, false, true, false).
		SetBorderCorner(false)

	s.TagInput = component.NewLine("Tag input")
	s.TagInput.SetWidthInherit(true)
	s.TagInput.SetBackground(ui.GetColorStr(ui.ColorFocus))

	s.TagOption.AddItem(s.TagLabel, component.NewFlexDescriptor(0, 1))
	s.TagOption.AddItem(s.TagInput, component.NewFlexDescriptor(0, 3))

	s.Window.AddItem(s.DashOption, component.NewFlexDescriptor(1, 0))
	s.Window.AddItem(s.TagOption, component.NewFlexDescriptor(1, 0))
	s.InitEvents()

	return s
}

func (s *Switcher) ToggleOption() {

	if s.DashOption.GetFocusState() {
		s.TagOption.Focus()
	} else {
		s.TagInput.SetText("")
		s.DashOption.Focus()
	}
}

func (s *Switcher) InitEvents() {
	s.Window.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "tab", "up", "down":
				s.ToggleOption()
			case "esc":
				s.dashboard.toggleSwitcher()
			}

		}
	}, false)

	s.DashOption.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter":
				s.dashboard.SwitchMode("dashboard", "")
			}
		}
	}, true)

	s.TagOption.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter":
				s.dashboard.SwitchMode("tag", s.TagInput.Text)
				s.TagInput.Text = ""
			case "backspace":
				if len(s.TagInput.Text) > 0 {
					s.TagInput.SetText(s.TagInput.Text[:len(s.TagInput.Text)-1])
				}

			default:
				str := string(msg.Code)
				if len(str) == 1 {
					s.TagInput.SetText(s.TagInput.Text + str)
				}
			}
		}
	}, true)

}
