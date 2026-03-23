package component

import (
	"strconv"
	"strings"
	"tumblr-dt/modules/ui/helper"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type GlobalValues struct {
	Msg             tea.Msg
	Time            int
	Elements        []Component
	EventDispatches map[string]map[string]func(tea.Msg, int)
}

var Global = &GlobalValues{}

func (g *GlobalValues) BlurAll() {
	for _, v := range g.Elements {
		v.Blur()
	}
}

func (g *GlobalValues) AddEventCallback(event string, uuid string, cb func(tea.Msg, int)) {
	if g.EventDispatches == nil {
		g.EventDispatches = map[string]map[string]func(tea.Msg, int){}
	}
	if g.EventDispatches[event] == nil {
		g.EventDispatches[event] = map[string]func(tea.Msg, int){}
	}
	g.EventDispatches[event][uuid] = cb
}

func (g *GlobalValues) CallEvents() {
	for _, v := range g.EventDispatches {
		for _, cb := range v {
			cb(g.Msg, g.Time)
		}
	}
	g.EventDispatches = map[string]map[string]func(tea.Msg, int){}
}

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
	SetX(int) *ComponentState
	SetY(int) *ComponentState
	SetPos(int, int) *ComponentState
	SetW(int) *ComponentState
	SetH(int) *ComponentState
	SetWidthInherit(bool) *ComponentState
	SetHeightInherit(bool) *ComponentState
	SetSize(int, int) *ComponentState
	Update()
	Focus()
	Blur()
	GetFocusState() bool
	SetDepth(int) *ComponentState
	SetParent(*ComponentState) *ComponentState
	GetSiblings() []Component
	GetChildren() []Component
	GetParent() Component
	GetCanvas() [][]string
	IsAbsolute() bool
	GetBorderPadding() int
	Propagate()
	GetName() string
	GetComponentName() string

	SetName(string) *ComponentState
	SetComponentName(string) *ComponentState
	AddEventListener(string, func(tea.Msg, int))
	SetStyle(lipgloss.Style) *ComponentState
	ClearStyle()
	GetStyle() lipgloss.Style

	SetBorder(bool) *ComponentState
	SetBorders(bool, bool, bool, bool) *ComponentState
	SetBorderCorner(bool) *ComponentState
	SetBorderPadding(int) *ComponentState

	GetBorderPaddings() (int, int, int, int)

	SetTitle(string) *ComponentState
	GetTitle() string
	DispatchEvent(string)

	SetIsFlexItem(bool) *ComponentState
	GetIsFlexItem() bool

	Trace([]string) []string
	GetTrace() []string

	GetEventCallbacks(string) []func(tea.Msg, int)
	GetComponent() Component

	GetUUID() string
}

type ComponentState struct {
	x                int
	y                int
	UUID             string
	Width            int
	Height           int
	InheritWidth     bool
	InheritHeight    bool
	Children         []Component
	Parent           Component
	Focused          bool
	Depth            int
	FitHeight        bool
	FitWidth         bool
	OnRenderReady    ([]func(Component))
	Canvas           [][]string
	ShowBorder       bool
	BorderPadWidth   int
	Name             string
	ComponentName    string
	EventCallbacks   map[string][]func(tea.Msg, int)
	Absolute         bool
	Overflow         bool
	Style            lipgloss.Style
	ShowTopBorder    bool
	ShowBottomBorder bool
	ShowLeftBorder   bool
	ShowRightBorder  bool
	ShowBorderCorner bool
	IsFlexItem       bool
	Title            string
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
	c.ShowTopBorder = true
	c.ShowBottomBorder = true
	c.ShowLeftBorder = true
	c.ShowRightBorder = true
	c.ShowBorderCorner = true
	c.IsFlexItem = false
	c.UUID = uuid.New().String()
	c.EventCallbacks = map[string][]func(tea.Msg, int){}

	Global.Elements = append(Global.Elements, c)
}

func (c *ComponentState) AddChild(child Component) {
	child.SetDepth(c.Depth + 1)
	child.SetParent(c)
	c.Children = append(c.Children, child)
	c.DispatchEvent("onAddChild")
}

func (c *ComponentState) SetParent(parent *ComponentState) *ComponentState {
	c.Parent = parent
	return c
}

func (c *ComponentState) SetDepth(v int) *ComponentState {
	c.Depth = v
	return c
}

func (c *ComponentState) GetRect() (int, int, int, int) {
	return c.GetX(), c.GetY(), c.GetWidth(), c.GetHeight()
}

func (c *ComponentState) GetRenderArea() (int, int, int, int) {
	if c.ShowBorder {
		t, _, l, _ := c.GetBorderPaddings()
		return c.GetX() + l, c.GetY() + t, c.GetInnerWidth(), c.GetInnerHeight()
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

func (c *ComponentState) SetX(v int) *ComponentState {
	c.x = v
	return c
}

func (c *ComponentState) SetY(v int) *ComponentState {
	c.y = v
	return c
}

func (c *ComponentState) SetPos(x int, y int) *ComponentState {
	c.SetX(x)
	c.SetY(y)
	return c
}

func (c *ComponentState) SetW(v int) *ComponentState {
	c.Width = v
	return c
}

func (c *ComponentState) SetH(v int) *ComponentState {
	c.Height = v
	return c
}

func (c *ComponentState) SetHeightInherit(v bool) *ComponentState {
	c.InheritHeight = v
	return c
}

func (c *ComponentState) SetWidthInherit(v bool) *ComponentState {
	c.InheritWidth = v
	return c
}

func (c *ComponentState) SetSize(w int, h int) *ComponentState {
	c.SetW(w)
	c.SetH(h)
	return c
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
		_, _, l, r := c.GetBorderPaddings()
		return c.GetWidth() - (l + r)
	}
	return c.GetWidth()
}

func (c *ComponentState) GetInnerHeight() int {
	if c.ShowBorder {
		t, b, _, _ := c.GetBorderPaddings()
		return c.GetHeight() - (t + b)
	}
	return c.GetHeight()
}

func (c *ComponentState) Update() {
	for _, child := range c.GetChildren() {
		child.Update()
	}

	if c.GetFocusState() {
		c.DispatchEvent("onUpdate")
	}
}

func (c *ComponentState) Focus() {
	Global.BlurAll()
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

	top, _, left, _ := b.GetBorderPaddings()
	cursor := top

	innerWidth := b.GetInnerWidth()
	innerHeight := b.GetInnerHeight()

	for _, c := range b.GetChildren() {
		c.PrepareFrame()
		output := c.GetCanvas()
		style := c.GetStyle()

		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := left + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + top
				for index, char := range line {
					result[posY][globalX+index] = char
				}
			}
		} else {
			for _, line := range output {
				if cursor > innerHeight {
					break
				}
				for i, char := range line {
					index := i
					if index >= innerWidth {
						break
					}
					result[cursor][index+left] = style.Render(char)
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
	if c.GetParent() != nil {
		return c.GetParent().GetChildren()
	}

	return []Component{}
}

func (c *ComponentState) GetParent() Component {
	return c.Parent
}

func (c *ComponentState) CreateCanvas() [][]string {
	var arr [][]string
	// width, height := c.GetContentsSize()
	width := c.GetWidth()
	height := c.GetHeight()

	for range height {
		arr = append(arr, strings.Split(strings.Repeat(" ", width), ""))
	}

	return arr
}

func (c *ComponentState) AddEventListener(event string, cb func(tea.Msg, int)) {
	list := c.GetEventCallbacks(event)
	list = append(list, cb)
	c.EventCallbacks[event] = list
}

func (c *ComponentState) DispatchEvent(event string) {
	var bubble []Component
	bubble = append(bubble, c)
	pt := c.GetParent()

	for pt != nil {
		bubble = append(bubble, pt)
		pt = pt.GetParent()
	}

	for _, element := range bubble {
		for _, cb := range element.GetEventCallbacks(event) {
			Global.AddEventCallback(event, element.GetUUID(), cb)
		}
	}
}

func (c *ComponentState) bubbleEvent(event string) {

}

func (c *ComponentState) GetCanvas() [][]string {
	return c.Canvas
}

func (c *ComponentState) addBorder(arr [][]string) [][]string {
	if !c.ShowBorder || c.GetBorderPadding() == 0 {
		return arr
	}

	side := helper.Dictionary(helper.BorderSide)
	top := helper.Dictionary(helper.BorderTop)

	wid := c.GetWidth()
	hei := c.GetHeight()

	if len(arr) > 0 && len(arr[0]) > 0 {

		for i := range c.GetWidth() {
			if c.ShowTopBorder {
				arr[0][i] = top
			}
			if c.ShowBottomBorder {
				arr[hei-1][i] = top
			}
		}

		for i := range c.GetHeight() {
			if c.ShowLeftBorder {
				arr[i][0] = side
			}
			if c.ShowRightBorder {
				arr[i][wid-1] = side
			}
		}
		if c.ShowBorderCorner {
			arr[0][0] = helper.Dictionary(helper.BorderTopLeft)
			arr[0][wid-1] = helper.Dictionary(helper.BorderTopRight)

			arr[hei-1][0] = helper.Dictionary(helper.BorderBottomLeft)
			arr[hei-1][wid-1] = helper.Dictionary(helper.BorderBottomRight)
		}

		title := c.Title
		for i, char := range title {
			arr[0][i+1] = string(char)
		}

		if c.GetFocusState() {
			arr[0][0] = "!"
		}

	}

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
			w = max(w, cw)
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

func (c *ComponentState) GetBorderPaddings() (int, int, int, int) {
	if c.ShowBorder {
		pad := c.GetBorderPadding()
		top := 0
		if c.ShowTopBorder {
			top = pad
		}

		bottom := 0
		if c.ShowBottomBorder {
			bottom = pad
		}

		right := 0
		if c.ShowRightBorder {
			right = pad
		}

		left := 0
		if c.ShowLeftBorder {
			left = pad
		}
		return top, bottom, left, right
	}
	return 0, 0, 0, 0
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

func (c *ComponentState) SetName(n string) *ComponentState {
	c.Name = n
	return c
}

func (c *ComponentState) SetComponentName(n string) *ComponentState {
	c.ComponentName = n
	return c
}

func (c *ComponentState) SetStyle(s lipgloss.Style) *ComponentState {
	c.Style = s
	return c
}

func (c *ComponentState) ClearStyle() {
	c.Style = lipgloss.NewStyle()
}

func (c *ComponentState) GetStyle() lipgloss.Style {
	return c.Style
}

func (c *ComponentState) SetBorder(show bool) *ComponentState {
	c.ShowBorder = show
	return c
}
func (c *ComponentState) SetBorders(top bool, bottom bool, left bool, right bool) *ComponentState {
	c.ShowTopBorder = top
	c.ShowBottomBorder = bottom
	c.ShowLeftBorder = left
	c.ShowRightBorder = right
	return c
}

func (c *ComponentState) SetBorderCorner(show bool) *ComponentState {
	c.ShowBorderCorner = show
	return c
}

func (c *ComponentState) SetBorderPadding(v int) *ComponentState {
	c.BorderPadWidth = v
	return c
}

func (c *ComponentState) SetTitle(str string) *ComponentState {
	c.Title = str
	return c
}

func (c *ComponentState) GetTitle() string {
	return c.Title
}

func (c *ComponentState) SetIsFlexItem(flex bool) *ComponentState {
	c.IsFlexItem = flex
	return c
}

func (c *ComponentState) GetIsFlexItem() bool {
	return c.IsFlexItem
}

func (c *ComponentState) Trace(list []string) []string {

	if c.GetParent() != nil {
		list = append(list, c.GetParent().Trace(list)...)
	}

	list = append(list, strconv.Itoa(c.Depth)+":"+c.GetComponentName()+"("+c.GetName()+")")
	return list
}

func (c *ComponentState) GetTrace() []string {
	return c.Trace([]string{})
}

func (c *ComponentState) GetEventCallbacks(event string) []func(tea.Msg, int) {
	return c.EventCallbacks[event]
}

func (c *ComponentState) GetComponent() Component {
	return c
}

func (c *ComponentState) GetUUID() string {
	return c.UUID
}
