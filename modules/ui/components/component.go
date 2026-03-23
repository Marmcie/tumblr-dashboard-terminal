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
	Initialize()
	//X,Y,Width,Height
	GetRect() (int, int, int, int)
	PrepareFrame()
	GetPos() (int, int)
	GetX() int
	GetY() int
	GetWidth() int
	GetHeight() int
	GetInnerWidth() int
	GetInnerHeight() int
	GetContentsSize() (int, int)
	GetContentsWidth() int
	GetContentsHeight() int
	SetX(int)
	SetY(int)
	SetPos(int, int)
	SetW(int)
	SetH(int)
	SetSize(int, int)
	Update()
	Focus()
	Blur()
	GetFocusState() bool
	SetDepth(int)
	SetParent(*ComponentState)
	GetSiblings() []Component
	GetChildren() []Component
	GetParent() Component
	GetCanvas() [][]string
	IsAbsolute() bool
	GetBorderPadding() int
	Propagate()
	GetName() string
	GetComponentName() string

	SetName(string)
	SetComponentName(string)
	AddEventListener(string, func(tea.Msg, int))
}

type ComponentState struct {
	x              int
	y              int
	Width          int
	Height         int
	InheritWidth   bool
	InheritHeight  bool
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
	ComponentName  string
	EventCallbacks map[string][]func(tea.Msg, int)
	Absolute       bool
	Overflow       bool
}

func (c *ComponentState) Initialize() {
	c.x = 0
	c.y = 0
	c.InheritWidth = false
	c.InheritWidth = false
	c.Focused = false
	c.Depth = 0
	c.FitHeight = false
	c.FitWidth = false
	c.ShowBorder = false
	c.BorderPadWidth = 0
	c.Name = "Name"
	c.ComponentName = "Base component"
	c.Absolute = false
	c.Parent = nil
}

func (c *ComponentState) AddChild(child Component) {
	child.SetDepth(c.Depth + 1)
	child.SetParent(c)
	c.Children = append(c.Children, child)
	c.DispatchEvent("onAddChild")
}

func (c *ComponentState) SetParent(parent *ComponentState) {
	c.Parent = parent
}

func (c *ComponentState) SetDepth(v int) {
	c.Depth = v
}

func (c *ComponentState) GetRect() (int, int, int, int) {
	return c.GetX(), c.GetY(), c.GetWidth(), c.GetHeight()
}

func (c *ComponentState) GetRenderArea() (int, int, int, int) {
	if c.ShowBorder {
		pad := c.GetBorderPadding()
		return c.GetX() + pad, c.GetY() + pad, c.GetInnerWidth(), c.GetInnerHeight()
	}

	return c.GetX(), c.GetY(), c.GetInnerWidth(), c.GetInnerHeight()
}

func (c *ComponentState) GetX() int {
	return c.x
}

func (c *ComponentState) GetY() int {
	return c.y
}

func (c *ComponentState) GetPos() (int, int) {
	return c.GetX(), c.GetY()
}

func (c *ComponentState) SetX(v int) {
	c.x = v
}

func (c *ComponentState) SetY(v int) {
	c.y = v
}

func (c *ComponentState) SetPos(x int, y int) {
	c.SetX(x)
	c.SetY(y)
}

func (c *ComponentState) SetW(v int) {
	c.Width = v
}

func (c *ComponentState) SetH(v int) {
	c.Height = v
}

func (c *ComponentState) SetSize(w int, h int) {
	c.SetW(w)
	c.SetH(h)
}

func (c *ComponentState) GetWidth() int {
	if c.InheritWidth == true && c.GetParent() != nil {
		return c.GetParent().GetInnerWidth()
	}
	return c.Width
}

func (c *ComponentState) GetHeight() int {
	if c.InheritHeight == true && c.GetParent() != nil {
		return c.GetParent().GetInnerHeight()
	}
	return c.Height
}

func (c *ComponentState) GetInnerWidth() int {
	if c.ShowBorder {
		return c.GetWidth() - (2 + ((c.GetBorderPadding() - 1) * 2))
	}
	return c.GetWidth()
}

func (c *ComponentState) GetInnerHeight() int {
	if c.ShowBorder {
		return c.GetHeight() - (2 + ((c.GetBorderPadding() - 1) * 2))
	}
	return c.GetHeight()
}

func (c *ComponentState) Update() {
	if c.GetFocusState() {
		c.DispatchEvent("onUpdate")
	}
	for _, child := range c.GetChildren() {
		if child.GetFocusState() {
			child.Update()
		}
	}
}

func (c *ComponentState) Focus() {
	for _, child := range c.GetChildren() {
		child.Focus()
	}
	c.Focused = true
}

func (c *ComponentState) Blur() {
	for _, child := range c.GetChildren() {
		child.Blur()
	}
	c.Focused = false
}

func (c *ComponentState) GetFocusState() bool {
	return c.Focused
}

func (b *ComponentState) PrepareFrame() {
	var result = b.CreateCanvas()

	cursor := b.GetBorderPadding()

	pad := b.GetBorderPadding() - 1
	for _, c := range b.GetChildren() {
		c.PrepareFrame()
		output := c.GetCanvas()
		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := pad + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + pad
				for index, char := range line {
					result[posY][globalX+index] = char
				}
			}
		} else {
			for _, line := range output {
				if cursor > b.GetInnerHeight() {
					break
				}
				for i, char := range line {
					index := i + b.GetBorderPadding()
					if index > b.GetInnerWidth() {
						break
					}
					result[cursor][index+pad] = char
				}
				cursor++
			}
		}
	}

	b.Canvas = result
	b.DispatchEvent("onRenderReady")
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

func (c *ComponentState) GetParent() Component {
	return c.Parent
}

func (c *ComponentState) CreateCanvas() [][]string {
	var arr [][]string
	width, height := c.GetContentsSize()

	for range height {
		arr = append(arr, strings.Split(strings.Repeat(" ", width), ""))
	}

	return arr
}

func (c *ComponentState) AddEventListener(event string, cb func(tea.Msg, int)) {
	if c.EventCallbacks == nil {
		c.EventCallbacks = map[string][]func(tea.Msg, int){}
	}
	list := c.EventCallbacks[event]
	list = append(list, cb)
	c.EventCallbacks[event] = list

	if c.Parent != nil {
		c.Parent.AddEventListener(event, cb)
	}

}

func (c *ComponentState) DispatchEvent(event string) {
	for _, v := range c.EventCallbacks[event] {
		v(Global.Msg, Global.Time)
	}
}

func (c *ComponentState) GetCanvas() [][]string {
	return c.Canvas
}

func (c *ComponentState) addBorder(arr [][]string) [][]string {
	if !c.ShowBorder || c.GetBorderPadding() == 0 || c.GetHeight() < (c.GetBorderPadding()*2) {
		return arr
	}

	side := helper.Dictionary(helper.BorderSide)
	top := helper.Dictionary(helper.BorderTop)
	wid := c.GetWidth()
	hei := c.GetHeight()
	for i := range c.GetHeight() {
		arr[i][0] = side
		arr[i][wid-1] = side
	}

	for i := range c.GetWidth() {
		arr[0][i] = top
		arr[hei-1][i] = top
	}

	arr[0][0] = helper.Dictionary(helper.BorderTopLeft)
	arr[0][wid-1] = helper.Dictionary(helper.BorderTopRight)

	arr[hei-1][0] = helper.Dictionary(helper.BorderBottomLeft)
	arr[hei-1][wid-1] = helper.Dictionary(helper.BorderBottomRight)

	return arr
}

func (c *ComponentState) IsAbsolute() bool {
	return c.Absolute
}

func (c *ComponentState) GetContentsSize() (int, int) {
	w := 0
	h := 0
	for _, child := range c.GetChildren() {
		cx, cy, cw, ch := child.GetRect()
		if child.IsAbsolute() {
			w = max(cx+cw, w)
			h = max(cy+ch, h)
		} else {
			// w = w + cw
			h = h + ch
		}
	}
	return max(w, c.GetWidth()), max(h, c.GetHeight())
}

func (c *ComponentState) GetContentsWidth() int {
	w, _ := c.GetContentsSize()
	return w
}

func (c *ComponentState) GetContentsHeight() int {
	_, h := c.GetContentsSize()
	return h
}

func (c *ComponentState) GetBorderPadding() int {
	if c.ShowBorder {
		return c.BorderPadWidth
	}
	return 0
}

func (c *ComponentState) Propagate() {
	for _, c := range c.GetChildren() {
		c.Propagate()
	}
}

func (c *ComponentState) GetName() string {
	return c.Name
}

func (c *ComponentState) GetComponentName() string {
	return c.ComponentName
}

func (c *ComponentState) SetName(n string) {
	c.Name = n
}

func (c *ComponentState) SetComponentName(n string) {
	c.ComponentName = n
}
