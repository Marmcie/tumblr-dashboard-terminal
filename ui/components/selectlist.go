package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selectlist struct {
	Scrollable

	OptionCallbacks []func()
	Cursor          int
	SizeList        []int
}

func NewSelectlist(name string) *Selectlist {
	s := &Selectlist{}
	s.Scrollable.Initialize(name)
	s.Cursor = 0
	s.ComponentName = "Selectlist"
	s.SizeList = append(s.SizeList, 0)

	s.AddEventListener("onUpdate", func(msg tea.Msg, time int) {

		switch msg := msg.(type) {

		// Is it a key press?
		case tea.KeyMsg:

			// Cool, what was the actual key pressed?
			switch msg.String() {

			// These keys should exit the program.
			case "j":
				s.Cursor = min(s.Cursor+1, len(s.OptionCallbacks)-1)
				s.DispatchEvent("onChange")

			case "k":
				s.Cursor = max(s.Cursor-1, 0)
				s.DispatchEvent("onChange")

			case "l":
				s.OptionCallbacks[s.Cursor]()

			case "enter":
				s.OptionCallbacks[s.Cursor]()
			}
		}

	})

	return s
}

func (s *Selectlist) UpdateOffset() {
	intended := s.SizeList[s.Cursor]

	if intended > s.OffsetY+s.GetInnerHeight() {
		s.OffsetY = s.GetInnerHeight() - intended
	}

}

func (c *Selectlist) AddOption(child Component, cb func()) {
	c.ComponentState.AddChild(child)
	c.OptionCallbacks = append(c.OptionCallbacks, cb)
}

func (s *Selectlist) Propagate() {

	children := s.GetChildren()
	if len(children) > 0 {
		s.SizeList = append(s.SizeList, s.SizeList[len(s.SizeList)-1]+children[len(children)-1].GetHeight())
	}
	for i, c := range s.GetChildren() {
		if i == s.Cursor {
			style := lipgloss.NewStyle().Background(lipgloss.Color("#444444"))
			c.SetStyle(style)
		} else {
			c.ClearStyle()
		}
	}
	s.UpdateOffset()

	s.Scrollable.Propagate()
}
