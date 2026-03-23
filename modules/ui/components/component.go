package component

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

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
	Update(tea.Msg, int)
	Focus()
	Blur()
	GetFocusState() bool
	SetDepth(int)
	SetParent(ComponentState)
	GetSiblings() []Component
	GetChildren() []Component
	GetCanvas() [][]string
	OnUpdate(tea.Msg, int)

	AddEventListener(string, func(Component))
}

type ComponentState struct {
	x             int
	y             int
	Width         int
	Height        int
	Children      []Component
	Parent        Component
	Focused       bool
	Depth         int
	FitHeight     bool
	FitWidth      bool
	OnRenderReady ([]func(Component))
	Canvas        [][]string
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

func (c *ComponentState) OnUpdate(msg tea.Msg, time int) {}

func (c *ComponentState) Update(msg tea.Msg, time int) {
	c.OnUpdate(msg, time)
	for _, child := range c.Children {
		if child.GetFocusState() {
			child.Update(msg, time)
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

func (c *ComponentState) PrepareFrame() {
}

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

func (c *ComponentState) AddEventListener(event string, cb func(Component)) {
	switch event {
	case "onRenderReady":
		c.OnRenderReady = append(c.OnRenderReady, cb)
	}
}

func (c *ComponentState) DispatchEvent(event string) {
	switch event {
	case "onRenderReady":
		for _, cb := range c.OnRenderReady {
			cb(c)
		}

	}
}

func (c *ComponentState) GetCanvas() [][]string {
	return c.Canvas
}
