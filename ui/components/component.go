package component

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"tumblr-dt/ui/helper"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
	"github.com/mattn/go-runewidth"
)

type Component interface {
	SetBorder(bool) *ComponentState
	SetBorderCorner(bool) *ComponentState
	GetBorderPadding() int
	SetBorderPadding(int) *ComponentState
	GetBorderPaddings() (int, int, int, int)
	GetBorderStyle() lipgloss.Style
	SetBorderStyle(lipgloss.Style) *ComponentState
	ResetBorderStyle() *ComponentState
	SetBorders(bool, bool, bool, bool) *ComponentState
	GetCanvas() [][]string
	AddChild(Component)
	GetChildren() []Component
	GetComponent() Component
	GetComponentName() string
	SetComponentName(string) *ComponentState
	GetContentsHeight() int
	GetContentsSize() (int, int)
	GetContentsWidth() int
	SetDepth(int) *ComponentState
	GetEventCallbacks(string) map[string]func(tea.Msg, int)
	AddEventListener(string, func(tea.Msg, int))
	GetFocusState() bool
	SetH(int) *ComponentState
	GetHeight() int
	SetHeightInherit(bool) *ComponentState
	GetInnerHeight() int
	GetInnerWidth() int
	GetIsFlexItem() bool
	SetIsFlexItem(bool) *ComponentState
	GetName() string
	SetName(string) *ComponentState
	GetParent() Component
	SetParent(*ComponentState) *ComponentState
	GetPos() (int, int)
	SetPos(int, int) *ComponentState
	GetRect() (int, int, int, int)
	GetSiblings() []Component
	SetSize(int, int) *ComponentState
	GetStyle() lipgloss.Style
	SetStyle(lipgloss.Style) *ComponentState
	GetTitle() string
	SetTitle(string) *ComponentState
	GetTitleAlignment() string
	SetTitleAlignment(string) *ComponentState
	GetTrace() []string
	GetUUID() string
	SetW(int) *ComponentState
	GetWidth() int
	SetWidthInherit(bool) *ComponentState
	GetX() int
	SetX(int) *ComponentState
	GetY() int
	SetY(int) *ComponentState
	ClearChildren()
	ClearStyle()
	Update()
	IsAbsolute() bool
	Trace([]string) []string
	Propagate()
	PrepareFrame()
	DispatchEvent(string)
	Blur()
	SetDoubleBorder(bool) *ComponentState
	GetDoubleBorder() bool
	Initialize(string)
	Focus()
	SetVisibility(bool)
	ToString() string
	SetBorderLabel(string, string)
	UpdateVisibility(int, int)
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
	EventCallbacks   map[string]map[string]func(tea.Msg, int)
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
	TitleAlignment   string
	BorderStyle      lipgloss.Style
	ShowDoubleBorder bool
	Visibility       bool
	BorderLabels     map[string]string
}

func (c *ComponentState) Initialize(name string) {
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
	c.Name = name
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
	c.EventCallbacks = map[string]map[string]func(tea.Msg, int){}
	c.TitleAlignment = "center"

	c.BorderLabels = map[string]string{
		"TopLeft":     "",
		"Top":         "",
		"TopRight":    "",
		"BottomRight": "",
		"Bottom":      "",
		"BottomLeft":  "",
	}

	c.ResetBorderStyle()
	c.ClearStyle()
	c.SetVisibility(true)

	c.ShowDoubleBorder = false
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
	c.DispatchEvent("onFocus")
	c.DispatchEvent("onFocusChange")
}

func (c *ComponentState) Blur() {
	for _, child := range c.GetChildren() {
		child.Blur()
	}
	c.Focused = false
	c.DispatchEvent("onBlur")
	c.DispatchEvent("onFocusChange")
}

func (c *ComponentState) GetFocusState() bool {
	return c.Focused
}

func (b *ComponentState) PrepareFrame() {
	var result = b.CreateCanvas()
	if !b.Visibility {
		b.Canvas = [][]string{{""}}
		return
	}

	top, _, left, _ := b.GetBorderPaddings()
	cursor := top

	innerWidth := b.GetInnerWidth() + 1
	// innerHeight := b.GetInnerHeight()

	style := b.GetStyle()
	for _, c := range b.GetChildren() {
		c.PrepareFrame()
		output := c.GetCanvas()

		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := left + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + top
				for index, char := range line {
					result[posY][globalX+index] = style.Render(char)
				}
			}
		} else {
			// Loop through lines
			pt := cursor
			for y := 0; y < min(c.GetHeight(), len(result), len(output)); y++ {
				line := output[y]
				// Loop through characters
				for x := range min(c.GetWidth(), innerWidth-left, len(line)) {
					// If canvas is smaller than the horizontal pointer, break
					char := line[x]
					// Check if the character is over the drawable area
					result[pt][x+left] = style.Render(char)
				}
				pt++
			}
			cursor += c.GetHeight()
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
	height := c.GetContentsHeight() + 1
	width := c.GetWidth()

	// height := c.GetHeight()

	for range height {
		arr = append(arr, strings.Split(strings.Repeat(" ", width), ""))
	}

	return arr
}

func (c *ComponentState) AddEventListener(event string, cb func(tea.Msg, int)) {
	list := c.GetEventCallbacks(event)
	list[uuid.New().String()] = cb
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
		for callbackUUID, cb := range element.GetEventCallbacks(event) {
			Global.AddEventCallback(event, element.GetUUID(), callbackUUID, cb)
		}
	}
}

func (c *ComponentState) bubbleEvent(event string) {

}

func (c *ComponentState) GetCanvas() [][]string {
	return c.Canvas
}

func (c *ComponentState) addBorder(arr [][]string) [][]string {
	if !c.ShowBorder || c.GetBorderPadding() == 0 || len(arr) < 3 || len(arr[0]) < 3 {
		return arr
	}

	style := c.GetBorderStyle()
	side := style.Render(helper.Dictionary(helper.BorderSide))
	top := style.Render(helper.Dictionary(helper.BorderTop))
	tl := style.Render(helper.Dictionary(helper.BorderTopLeft))
	tr := style.Render(helper.Dictionary(helper.BorderTopRight))
	bl := style.Render(helper.Dictionary(helper.BorderBottomLeft))
	br := style.Render(helper.Dictionary(helper.BorderBottomRight))

	if c.GetFocusState() || c.GetDoubleBorder() {
		side = style.Render(helper.Dictionary(helper.BorderSideDouble))
		top = style.Render(helper.Dictionary(helper.BorderTopDouble))
		tl = style.Render(helper.Dictionary(helper.BorderTopLeftDouble))
		tr = style.Render(helper.Dictionary(helper.BorderTopRightDouble))
		bl = style.Render(helper.Dictionary(helper.BorderBottomLeftDouble))
		br = style.Render(helper.Dictionary(helper.BorderBottomRightDouble))
	}

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
			arr[0][0] = tl
			arr[0][wid-1] = tr

			arr[hei-1][0] = bl
			arr[hei-1][wid-1] = br
		}

		title := c.GetTitle()

		switch c.GetTitleAlignment() {
		case "left":
			for i := range min(wid-1, runewidth.StringWidth(title)) {
				char := title[i]
				arr[0][i+1] = style.Render(string(char))
			}

		case "center":
			length := len(title)
			for i := range min(wid-1, runewidth.StringWidth(title)) {
				char := title[i]
				arr[0][i+max(1, (wid-length)/2)] = style.Render(string(char))
			}

		case "right":
			strWidth := len(title)
			for i := 0; i < min(wid-2, strWidth); i++ {
				char := title[strWidth-(i+1)]
				arr[0][wid-(i+2)] = style.Render(string(char))
			}
		}

		for key, str := range c.BorderLabels {
			if len(str) == 0 {
				continue
			}
			switch key {
			case "TopLeft":
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := title[i]
					arr[0][i+1] = style.Render(string(char))
				}

			case "Top":
				length := len(str)
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[0][i+max(1, (wid-length)/2)] = style.Render(string(char))
				}

			case "TopRight":
				strWidth := len(str)
				for i := 0; i < min(wid-2, strWidth); i++ {
					char := str[strWidth-(i+1)]
					arr[0][wid-(i+2)] = style.Render(string(char))
				}

			case "BottomLeft":
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[hei-1][i+1] = style.Render(string(char))
				}

			case "Bottom":
				length := len(str)
				center := wid - length/2
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[hei-1][i+max(1, center)] = style.Render(string(char))
				}

			case "BottomRight":
				strWidth := len(str)
				for i := 0; i < min(wid-2, strWidth); i++ {
					char := str[strWidth-(i+1)]
					arr[hei-1][wid-(i+2)] = style.Render(string(char))
				}
			}

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

func (c *ComponentState) GetContentsHeight() int {
	h := 0
	for _, child := range c.GetChildren() {
		_, cy, _, ch := child.GetRect()
		if child.IsAbsolute() {
			h = max(cy+ch, h)
		} else {
			h = h + ch
		}
	}
	return max(h, c.GetHeight())
}

func (c *ComponentState) GetContentsWidth() int {
	w, _ := c.GetContentsSize()
	return w
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
	if show && c.BorderPadWidth == 0 {
		c.SetBorderPadding(1)
	}
	return c
}

func (c *ComponentState) SetBorderStyle(style lipgloss.Style) *ComponentState {
	c.BorderStyle = style
	return c
}

func (c *ComponentState) GetBorderStyle() lipgloss.Style {
	return c.BorderStyle
}
func (c *ComponentState) ResetBorderStyle() *ComponentState {
	c.BorderStyle = lipgloss.NewStyle()
	return c
}

func (c *ComponentState) SetTitleAlignment(str string) *ComponentState {
	c.Title = str
	return c
}

func (c *ComponentState) GetTitleAlignment() string {
	return c.TitleAlignment
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

func (c *ComponentState) GetEventCallbacks(event string) map[string]func(tea.Msg, int) {
	if c.EventCallbacks[event] == nil {
		c.EventCallbacks[event] = map[string]func(tea.Msg, int){}
	}
	return c.EventCallbacks[event]
}

func (c *ComponentState) GetComponent() Component {
	return c
}

func (c *ComponentState) GetUUID() string {
	return c.UUID
}

func (c *ComponentState) ClearChildren() {
	c.Children = []Component{}
}

func (c *ComponentState) SetDoubleBorder(v bool) *ComponentState {
	c.ShowDoubleBorder = v
	return c
}

func (c *ComponentState) GetDoubleBorder() bool {
	return c.ShowDoubleBorder
}

func (c *ComponentState) SetVisibility(v bool) {
	c.Visibility = v
}

func (c *ComponentState) SetBorderLabel(key string, str string) {
	c.BorderLabels[key] = str
}
func (c *ComponentState) ToString() string {
	var res bytes.Buffer

	inherited := ""
	fmt.Fprintf(&res, "----------------\n")
	fmt.Fprintf(&res, "Name : %s", c.GetName())
	fmt.Fprintf(&res, "\n")
	fmt.Fprintf(&res, "Component Name : %s", c.GetComponentName())
	fmt.Fprintf(&res, "\n")
	fmt.Fprintf(&res, "Children count : %d", len(c.GetChildren()))
	fmt.Fprintf(&res, "\n")

	if c.InheritWidth {
		inherited = "(Inherited)"
	} else {
		inherited = ""
	}
	fmt.Fprintf(&res, "Width : %d %s", c.GetWidth(), inherited)
	fmt.Fprintf(&res, "\n")
	fmt.Fprintf(&res, "Inner width : %d %s", c.GetInnerWidth(), inherited)
	fmt.Fprintf(&res, "\n")

	if c.InheritHeight {
		inherited = "(Inherited)"
	} else {
		inherited = ""
	}
	fmt.Fprintf(&res, "Height : %d %s", c.GetHeight(), inherited)
	fmt.Fprintf(&res, "\n")
	fmt.Fprintf(&res, "Inner height : %d %s", c.GetInnerHeight(), inherited)
	fmt.Fprintf(&res, "\n")
	fmt.Fprintf(&res, "----------------\n")

	return res.String()
}

func (c *ComponentState) UpdateVisibility(ytop int, hei int) {
	top := 0
	y := ytop
	h := hei
	for _, child := range c.GetChildren() {
		childHeight := child.GetHeight()
		child.SetVisibility(!(top+childHeight < y || top > y+h))
		top += childHeight
	}
}
