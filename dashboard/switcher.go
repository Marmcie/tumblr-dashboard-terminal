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
	TagInput   *component.Input
	BlogInput  *component.Input
	BlogOption *component.Flex

	SearchInput  *component.Input
	SearchOption *component.Flex
	dashboard    *Dashboard
	index        int
}

func NewSwitcher(dashboard *Dashboard) *Switcher {
	s := &Switcher{}
	s.dashboard = dashboard
	s.Window = component.NewFlex("Switcher window")
	s.Window.
		SetAbsolute(true).
		SetCentered(true).
		SetSize(50, 11).
		SetBorder(true)
	s.Window.SetForeground(ui.GetColorStr(ui.ColorWhite))
	s.Window.SetBorderForeground(ui.GetColorStr(ui.ColorWhite))
	s.Window.SetBorderLabelColor("BottomRight", ui.GetColorStr(ui.ColorWhite))

	s.Window.SetTitle("Feed picker")
	s.Window.SetBorderLabel("BottomRight", "Esc to close")

	s.DashOption = component.NewBox("Dash option")
	s.DashOption.SetBorder(true).
		SetWidthInherit(true).
		SetBorders(false, false, true, false).
		SetPaddings(0, 0, 1, 0).
		SetBorderCorner(false)

	s.DashLabel = component.NewLine("Dash label")
	s.DashLabel.SetText("Dashboard")

	s.DashOption.AddChild(s.DashLabel)

	tagLabel := component.NewLine("Tag label")
	tagLabel.SetText("Tag : ")

	s.TagOption = component.NewFlex("Tag option")
	s.TagOption.SetDirection(1).SetBorder(true).SetWidthInherit(true).SetBorders(false, false, true, false).SetBorderCorner(false).SetPadding(0)

	s.TagInput = component.NewInput("Tag input")
	s.TagInput.SetPlaceholder("Input tag here").SetWidthInherit(true).SetBackground(ui.GetColorStr(ui.ColorFocus))

	s.TagOption.AddItem(tagLabel, 0, 1)
	s.TagOption.AddItem(s.TagInput, 0, 3)

	blogLabel := component.NewLine("Blog label")
	blogLabel.SetText("Blog name : ")

	s.BlogOption = component.NewFlex("Blog option")
	s.BlogOption.SetDirection(1).SetBorder(true).SetWidthInherit(true).SetBorders(false, false, true, false).SetBorderCorner(false).SetPadding(0)

	s.BlogInput = component.NewInput("Blog input")
	s.BlogInput.SetPlaceholder("Input blog name here").SetWidthInherit(true).SetBackground(ui.GetColorStr(ui.ColorFocus))

	s.BlogOption.AddItem(blogLabel, 0, 1)
	s.BlogOption.AddItem(s.BlogInput, 0, 3)

	SearchLabel := component.NewLine("Search label")
	SearchLabel.SetText("Search : ")

	s.SearchOption = component.NewFlex("Search option")
	s.SearchOption.SetDirection(1).SetBorder(true).SetWidthInherit(true).SetBorders(false, false, true, false).SetBorderCorner(false).SetPadding(0)

	s.SearchInput = component.NewInput("Search input")
	s.SearchInput.SetPlaceholder("Input Search term here").SetWidthInherit(true).SetBackground(ui.GetColorStr(ui.ColorFocus))

	s.SearchOption.AddItem(SearchLabel, 0, 1)
	s.SearchOption.AddItem(s.SearchInput, 0, 3)

	s.Window.AddItem(s.DashOption, 1, 0)
	s.Window.AddItem(s.TagOption, 1, 0)
	s.Window.AddItem(s.BlogOption, 1, 0)
	s.Window.AddItem(s.SearchOption, 1, 0)
	s.index = 0
	s.InitEvents()

	return s
}

func (s *Switcher) ToggleOption() {
	switch s.index {
	case 0:
		s.DashOption.Focus()
		s.TagInput.ClearInput()
		s.BlogInput.ClearInput()
		s.SearchInput.ClearInput()

	case 1:
		s.TagOption.Focus()
		s.BlogInput.ClearInput()
		s.SearchInput.ClearInput()

	case 2:
		s.BlogOption.Focus()
		s.TagInput.ClearInput()
		s.SearchInput.ClearInput()

	case 3:
		s.SearchOption.Focus()
		s.TagInput.ClearInput()
		s.BlogInput.ClearInput()
	}

}

func (s *Switcher) InitEvents() {
	s.Window.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case s.dashboard.config.Keymaps.Switcher.Down:
				s.index = (s.index + 1) % 4
				s.ToggleOption()

			case s.dashboard.config.Keymaps.Switcher.Up:
				if s.index == 0 {
					s.index = 3
				} else {
					s.index = (s.index - 1) % 4
				}
				s.ToggleOption()

			case s.dashboard.config.Keymaps.Switcher.Close:
				s.index = 0
				s.TagInput.ClearInput()
				s.dashboard.toggleSwitcher()
			}

		}
	}, false)

	s.DashOption.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case s.dashboard.config.Keymaps.Confirm:
				s.index = 0
				s.dashboard.SwitchMode("dashboard", "")
			}
		}
	}, true)

	s.TagOption.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			default:
				s.TagInput.ParseInput(msg)
			case s.dashboard.config.Keymaps.Confirm:
				s.dashboard.SwitchMode("tag", s.TagInput.Value)
				s.TagInput.ClearInput()
				s.index = 0

			case s.dashboard.config.Keymaps.Switcher.Suggestion:
				s.TagInput.ApplyTopSuggestion()

			}
		}
	}, true)

	s.BlogOption.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			default:
				s.BlogInput.ParseInput(msg)
			case s.dashboard.config.Keymaps.Confirm:
				s.dashboard.SwitchMode("blog", s.BlogInput.Value)
				s.BlogInput.ClearInput()
				s.index = 0
			case s.dashboard.config.Keymaps.Switcher.Suggestion:
				s.BlogInput.ApplyTopSuggestion()

			}
		}
	}, true)

	s.SearchOption.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			default:
				s.SearchInput.ParseInput(msg)
			case s.dashboard.config.Keymaps.Confirm:
				s.dashboard.SwitchMode("search", s.SearchInput.Value)
				s.SearchInput.ClearInput()
				s.index = 0
			case s.dashboard.config.Keymaps.Switcher.Suggestion:
				// s.SearchInput.ApplyTopSuggestion()

			}
		}
	}, true)

}
