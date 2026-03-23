package component

// Component that displays a list of elements that can be selected and has a corresponding callbacks.
type Selectlist struct {
	Scrollable

	OptionCallbacks []func()
	Cursor          int
	SizeList        []int
	SelectedBG      string
	SelectedFG      string
}

func NewSelectlist(name string) *Selectlist {
	s := &Selectlist{}
	s.Scrollable.Initialize(name)
	s.Cursor = 0
	s.ComponentName = "Selectlist"
	s.SizeList = append(s.SizeList, 0)

	return s
}

func (s *Selectlist) SetSelectedOptionForeground(fg string) {
	s.SelectedFG = fg
}

func (s *Selectlist) SetSelectedOptionBackground(bg string) {
	s.SelectedBG = bg
}

func (s *Selectlist) IncrementCursor() {
	s.SetCursor(min(s.Cursor+1, len(s.OptionCallbacks)-1))
}

func (s *Selectlist) DecrementCursor() {
	s.SetCursor(max(s.Cursor-1, 0))
}

func (s *Selectlist) SetCursor(v int) {
	prev := s.Cursor
	children := s.GetChildren()
	children[prev].ClearBackground()
	children[prev].ClearForeground()
	children[v].SetBackground(s.SelectedBG)
	children[v].SetForeground(s.SelectedFG)
	s.Cursor = max(min(len(s.OptionCallbacks)-1, v), 0)
	s.DispatchEvent("onChange")
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

func (s *Selectlist) AddOption(child Component, cb func()) {
	s.ComponentState.AddChild(child)
	s.OptionCallbacks = append(s.OptionCallbacks, cb)

	children := s.GetChildren()
	if len(children) > 0 {
		s.SizeList = append(s.SizeList, s.SizeList[len(s.SizeList)-1]+children[len(children)-1].GetHeight())
	}
}

func (s *Selectlist) RunSelectedOption() {
	if len(s.OptionCallbacks) >= s.Cursor {
		s.OptionCallbacks[s.Cursor]()
	}
}

func (s *Selectlist) Propagate() {

	s.UpdateOffset()

	s.Scrollable.Propagate()
}

func (c *Selectlist) ClearChildren() {
	c.ComponentState.ClearChildren()
	c.OptionCallbacks = []func(){}
	c.Cursor = 0
}
