package component

import (
	"strings"
	"tumblr-dt/modules/ui/helper"

	tea "github.com/charmbracelet/bubbletea"
)

type GlobalValues struct {
	Msg  tea.Msg
	Time int
}

var Global = &GlobalValues{}

func UpdateGlobalValues(msg tea.Msg, time int) {
	Global.Msg = msg
	Global.Time = time
}

type Component interface {
	//X,Y,Width,Height
	GetRect() (int, int, int, int)
	PrepareFrame()
	GetPos() (int, int)
	GetX() int
	GetY() int
	SetX(int)
	SetY(int)
	SetPos(int, int)
	Update()
	Focus()
	Blur()
	GetFocusState() bool
	SetDepth(int)
	SetParent(ComponentState)
	GetSiblings() []Component
	GetChildren() []Component
	GetCanvas() [][]string

	AddEventListener(string, func(tea.Msg, int))
}

type ComponentState struct {
	x              int
	y              int
	Width          int
	Height         int
	Children       []Component
	Parent         Component
	Focused        bool
	Depth          int
	FitHeight      bool
	FitWidth       bool
	OnRenderReady  ([]func(Component))
	Canvas         [][]string
	ShowBorder     bool
	BorderPadWidth int
	Name           string
	EventCallbacks map[string][]func(tea.Msg, int)
}

func (c *ComponentState) AddChild(child Component) {
	child.SetDepth(c.Depth + 1)
	child.SetParent(*c)
	c.Children = append(c.Children, child)
}

func (c *ComponentState) SetParent(parent ComponentState) {
	c.Parent = &parent
}

func (c *ComponentState) SetDepth(v int) {
	c.Depth = v
}

func (c *ComponentState) GetRect() (int, int, int, int) {
	return c.x, c.y, c.Width, c.Height
}

func (c *ComponentState) GetRenderArea() (int, int, int, int) {
	return c.x, c.y, c.Width, c.Height
}

func (c *ComponentState) GetX() int {
	return c.x
}

func (c *ComponentState) GetY() int {
	return c.y
}

func (c *ComponentState) GetPos() (int, int) {
	return c.x, c.y
}

func (c *ComponentState) SetX(v int) {
	c.x = v
}

func (c *ComponentState) SetY(v int) {
	c.y = v
}

func (c *ComponentState) SetPos(x int, y int) {
	c.x = x
	c.y = y
}

func (c *ComponentState) Update() {
	if c.GetFocusState() {
		c.DispatchEvent("onUpdate")
	}
	for _, child := range c.Children {
		if child.GetFocusState() {
			child.Update()
		}
	}
}

func (c *ComponentState) Focus() {
	for _, child := range c.Children {
		child.Focus()
	}
	c.Focused = true
}

func (c *ComponentState) Blur() {
	for _, child := range c.Children {
		child.Blur()
	}
	c.Focused = false
}

func (c *ComponentState) GetFocusState() bool {
	return c.Focused
}

func (c *ComponentState) PrepareFrame() {}

func (c *ComponentState) GetChildren() []Component {
	return c.Children
}

func (c *ComponentState) GetSiblings() []Component {
	if c.Parent != nil {
		return c.Parent.GetChildren()
	}

	return []Component{}
}

func (c *ComponentState) CreateCanvas() [][]string {
	var arr [][]string

	for range c.Height {
		arr = append(arr, strings.Split(strings.Repeat(" ", c.Width), ""))
	}

	return arr
}

func (c *ComponentState) AddEventListener(event string, cb func(tea.Msg, int)) {
	if c.EventCallbacks == nil{
		c.EventCallbacks = map[string][]func(tea.Msg, int){}
	}
	list := c.EventCallbacks[event]
	list = append(list, cb)
	c.EventCallbacks[event] = list
}

func (c *ComponentState) DispatchEvent(event string) {
	for _, v := range c.EventCallbacks[event] {
		v(Global.Msg,Global.Time)
	}
}

func (c *ComponentState) GetCanvas() [][]string {
	return c.Canvas
}

func (c *ComponentState) addBorder(arr [][]string) [][]string {
	if !c.ShowBorder || c.BorderPadWidth == 0 {
		return arr
	}

	side := helper.Dictionary(helper.BorderSide)
	top := helper.Dictionary(helper.BorderTop)
	for i := range c.Height {
		arr[i][0] = side
		arr[i][c.Width-1] = side
	}

	for i := range c.Width {
		arr[0][i] = top
		arr[c.Height-1][i] = top
	}

	arr[0][0] = helper.Dictionary(helper.BorderTopLeft)
	arr[0][c.Width-1] = helper.Dictionary(helper.BorderTopRight)

	arr[c.Height-1][0] = helper.Dictionary(helper.BorderBottomLeft)
	arr[c.Height-1][c.Width-1] = helper.Dictionary(helper.BorderBottomRight)

	return arr
}
