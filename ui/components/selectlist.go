package component

import (
	"charm.land/lipgloss/v2"
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

	s.SelectBgStyle = lipgloss.NewStyle()
	return s
}

func (s *Selectlist) IncrementCursor() {
	s.Cursor = min(s.Cursor+1, len(s.OptionCallbacks)-1)
	s.DispatchEvent("onChange")
}

func (s *Selectlist) DecrementCursor() {
	s.Cursor = max(s.Cursor-1, 0)
	s.DispatchEvent("onChange")
}

func (s *Selectlist) SetCursor(v int) {
	if s.Cursor != v {
		s.Cursor = max(min(len(s.OptionCallbacks)-1, v), 0)
		s.DispatchEvent("onChange")
	} else {
		s.Cursor = max(min(len(s.OptionCallbacks)-1, v), 0)
	}
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
