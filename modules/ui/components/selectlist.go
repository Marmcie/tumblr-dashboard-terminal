package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Selectlist struct {
	Scrollable

	OptionCallbacks []func()
	Cursor          int
}

func NewSelectlist() *Selectlist {
	s := &Selectlist{}
	s.Scrollable.Initialize()
	s.Cursor = 0

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
			}
		}
		
	})
	return s
}

func (c *Selectlist) AddOption(child Component, cb func()) {
	c.ComponentState.AddChild(child)
	c.OptionCallbacks = append(c.OptionCallbacks, cb)
}

func (s *Selectlist) Propagate() {
	for i, c := range s.GetChildren() {
		if i == s.Cursor {
			style := lipgloss.NewStyle().Background(lipgloss.Color("#444444"))
			c.SetStyle(style)
		} else {
			c.ClearStyle()
		}
	}

}
