package dashboard

import (
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
)

type Switcher struct {
	Window     *component.Flex
	DashOption *component.Box
	TagOption  *component.Box
	TagInput   *component.Line
	dashboard  *Dashboard
}

func NewSwitcher(dashboard *Dashboard) *Switcher {
	s := &Switcher{}
	s.dashboard = dashboard
	s.Window = component.NewFlex("Switcher window")
	s.Window.
		SetAbsolute(true).
		SetCentered(true).
		SetSize(40, 8).
		SetBorder(true)
	s.DashOption = component.NewBox("Dash option")
	s.DashOption.SetBorder(true).
		SetWidthInherit(true)
	s.DashOption.SetTitle("Dashboard")

	s.TagOption = component.NewBox("Tag option")
	s.TagOption.SetBorder(true).
		SetWidthInherit(true)
	s.TagOption.SetTitle("Tag")

	s.TagInput = component.NewLine("Tag input")
	s.TagInput.SetWidthInherit(true)

	s.TagOption.AddChild(s.TagInput)
	s.Window.AddItem(s.DashOption, component.NewFlexDescriptor(0, 1))
	s.Window.AddItem(s.TagOption, component.NewFlexDescriptor(0, 2))
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
			case "tab":
				s.ToggleOption()
			}
		}
	}, false)

	s.DashOption.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "tab":
				s.ToggleOption()
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
			case "tab":
				s.ToggleOption()
			case "backspace":
				if len(s.TagInput.Text) > 0 {
					s.TagInput.SetText(s.TagInput.Text[:len(s.TagInput.Text)-1])
				}

			default:
				s.TagInput.SetText(s.TagInput.Text + string(msg.Code))
			}
		}
	}, false)

}
