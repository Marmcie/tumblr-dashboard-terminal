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
	SelectBgStyle   lipgloss.Style
}

func NewSelectlist(name string) *Selectlist {
	s := &Selectlist{}
	s.Scrollable.Initialize(name)
	s.Cursor = 0
	s.ComponentName = "Selectlist"
	s.SizeList = append(s.SizeList, 0)

	s.AddEventListener("onUpdate", func(msg tea.Msg, time int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "j":
				s.Cursor = min(s.Cursor+1, len(s.OptionCallbacks)-1)
				s.DispatchEvent("onChange")

			case "k":
				s.Cursor = max(s.Cursor-1, 0)
				s.DispatchEvent("onChange")

			case "l":
				if len(s.OptionCallbacks) > 0 {
					s.RunSelectedOption()
				}

			case "enter":
				if len(s.OptionCallbacks) > 0 {
					s.RunSelectedOption()
				}
			}
		}
	})
	s.SelectBgStyle = lipgloss.NewStyle().Background(lipgloss.Color("#444444"))
	return s
}

func (s *Selectlist) UpdateOffset() {
	if len(s.SizeList) > 1 {
		intended := s.SizeList[s.Cursor+1]
		innerHeight := s.GetInnerHeight()
		if intended > s.OffsetY+innerHeight {
			s.OffsetY = intended - innerHeight
		} else {
			if s.SizeList[s.Cursor] < s.OffsetY {
				s.OffsetY = s.SizeList[s.Cursor]
			}
		}
	}
}

func (c *Selectlist) AddOption(child Component, cb func()) {
	c.ComponentState.AddChild(child)
	c.OptionCallbacks = append(c.OptionCallbacks, cb)
}

func (s *Selectlist) RunSelectedOption() {
	if len(s.OptionCallbacks) >= s.Cursor {
		s.OptionCallbacks[s.Cursor]()
	}
}

func (s *Selectlist) Propagate() {

	children := s.GetChildren()
	if len(children) > 0 {
		s.SizeList = append(s.SizeList, s.SizeList[len(s.SizeList)-1]+children[len(children)-1].GetHeight())
	}
	for i, c := range s.GetChildren() {
		if i == s.Cursor {
			style := s.SelectBgStyle
			c.SetStyle(style)
		} else {
			c.ClearStyle()
		}
	}
	s.UpdateOffset()

	s.Scrollable.Propagate()
}

func (c *Selectlist) ClearChildren() {
	c.Children = []Component{}
	c.OptionCallbacks = []func(){}
}
