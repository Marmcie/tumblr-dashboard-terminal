package component

import (
	"bytes"
	"fmt"
	"image/color"
	"strconv"
	"tumblr-dt/ui/helper"

	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"
	"github.com/rivo/uniseg"
)

type Component interface {
	SetBorder(bool) *BaseComponent
	SetBorderForeground(string) *BaseComponent
	ClearBorderForeground() *BaseComponent
	SetBorderFocusForeground(string) *BaseComponent
	ClearBorderFocusForeground() *BaseComponent
	SetBorderCorner(bool) *BaseComponent
	SetPadding(int) *BaseComponent
	SetPaddings(int, int, int, int) *BaseComponent
	GetPaddings() (int, int, int, int)
	SetBorders(bool, bool, bool, bool) *BaseComponent
	GetCanvas() ([][]string, [][]string, [][]string)
	SetCanvas([][]string, [][]string, [][]string)
	AddChild(Component)
	GetChildren() []Component
	GetComponent() Component
	GetComponentName() string
	SetComponentName(string) *BaseComponent
	GetContentsHeight() int
	GetContentsSize() (int, int)
	SetDepth(int) *BaseComponent
	GetEventCallbacks(string) map[string]EventCb
	AddEventListener(string, func(tea.Msg), bool)
	GetFocusState() bool
	SetH(int) *BaseComponent
	GetHeight() int
	SetHeightInherit(bool) *BaseComponent
	GetInnerHeight() int
	GetInnerWidth() int
	GetIsFlexItem() bool
	SetIsFlexItem(bool) *BaseComponent
	GetName() string
	SetName(string) *BaseComponent
	GetParent() Component
	SetParent(*BaseComponent) *BaseComponent
	GetPos() (int, int)
	SetPos(int, int) *BaseComponent
	GetRect() (int, int, int, int)
	GetSiblings() []Component
	SetSize(int, int) *BaseComponent
	SetBackgroundGradient([]color.Color) *BaseComponent
	SetForegroundGradient([]color.Color) *BaseComponent
	GetBackgroundGradient() []color.Color
	GetForegroundGradient() []color.Color
	ClearBackgroundGradient() *BaseComponent
	ClearForegroundGradient() *BaseComponent
	GetTitle() string
	SetTitle(string) *BaseComponent
	GetTitleAlignment() string
	SetTitleAlignment(string) *BaseComponent
	GetTrace() []string
	GetUUID() string
	SetW(int) *BaseComponent
	GetWidth() int
	SetWidthInherit(bool) *BaseComponent
	GetX() int
	SetX(int) *BaseComponent
	GetY() int
	SetY(int) *BaseComponent
	ClearChildren()
	Update()
	IsAbsolute() bool
	Trace([]string) []string
	BeforeRender()
	RenderToCanvas()
	DispatchEvent(string)
	Blur()
	SetDoubleBorder(bool) *BaseComponent
	GetDoubleBorder() bool
	Initialize(string)
	Focus()
	SetVisibility(bool) *BaseComponent
	GetVisibility() bool
	ToString() string
	SetBorderLabel(string, string)
	SetBorderLabelColor(string, string)
	UpdateVisibility(int, int)
	Delete()
	SetGlobalIndex(int)
	SetAbsolute(bool) *BaseComponent
	SetCentered(bool) *BaseComponent
	GetCentered() bool
	SetForeground(string) *BaseComponent
	SetBackground(string) *BaseComponent
	ClearForeground() *BaseComponent
	ClearBackground() *BaseComponent
	GetForeground() string
	GetBackground() string
	SetFlexProportion(float32) *BaseComponent
	GetFlexProportion() float32
	SetMinWidth(int) *BaseComponent
	GetMinWidth() int
	SetMinHeight(int) *BaseComponent
	GetMinHeight() int
}

// Base class for all components
type BaseComponent struct {
	// X coordinates
	x int
	// Y coordinates
	y             int
	PaddingTop    int
	PaddingBottom int
	PaddingLeft   int
	PaddingRight  int
	Centered      bool
	UUID          string
	Width         int
	Height        int
	InheritWidth  bool
	InheritHeight bool
	Children      []Component
	Parent        Component
	Focused       bool
	Depth         int
	FitHeight     bool
	FitWidth      bool
	Canvas        [][]string
	BGSheet       [][]string
	FGSheet       [][]string

	ShowBorder bool
	// Name of an individual component
	Name string
	// Name of a component type
	ComponentName     string
	EventCallbacks    map[string]map[string]EventCb
	Absolute          bool
	Overflow          bool
	ShowTopBorder     bool
	ShowBottomBorder  bool
	ShowLeftBorder    bool
	ShowRightBorder   bool
	ShowBorderCorner  bool
	IsFlexItem        bool
	Title             string
	TitleAlignment    string
	ShowDoubleBorder  bool
	Visibility        bool
	BorderLabels      map[string]string
	BorderLabelColors map[string]string
	//Index of the component within global element list
	GlobalIndex           int
	BackgroundGradient    []color.Color
	ForegroundGradient    []color.Color
	Background            string
	Foreground            string
	BorderForeground      string
	BorderFocusForeground string
	FlexProportion        float32
	MinHeight             int
	MinWidth              int
}

type EventCb struct {
	Cb     func(tea.Msg)
	Bubble bool
}

// Initialized all shared values
func (c *BaseComponent) Initialize(name string) {
	c.x = 0
	c.y = 0
	c.Centered = false
	c.InheritWidth = false
	c.InheritWidth = false
	c.Focused = false
	c.Depth = 0
	c.FitHeight = false
	c.FitWidth = false
	c.ShowBorder = false
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
	c.EventCallbacks = map[string]map[string]EventCb{}
	c.TitleAlignment = "center"
	c.Foreground = ""
	c.Background = ""
	c.BorderForeground = ""
	c.BorderFocusForeground = ""
	c.SetPadding(0)

	c.BorderLabels = map[string]string{
		"TopLeft":     "",
		"Top":         "",
		"TopRight":    "",
		"BottomRight": "",
		"Bottom":      "",
		"BottomLeft":  "",
	}
	c.BorderLabelColors = map[string]string{
		"TopLeft":     "",
		"Top":         "",
		"TopRight":    "",
		"BottomRight": "",
		"Bottom":      "",
		"BottomLeft":  "",
	}
	c.SetFlexProportion(1)
	c.SetMinHeight(0)
	c.SetMinWidth(0)

	c.SetVisibility(true)

	c.ShowDoubleBorder = false
	c.GlobalIndex = Global.AddElement(c)
}

// #region Component relation

// Adds a child to an component
func (c *BaseComponent) AddChild(child Component) {
	child.SetDepth(c.Depth + 1)
	child.SetParent(c)
	c.Children = append(c.Children, child)
	c.DispatchEvent("onAddChild")
}

// Set parent of a component
func (c *BaseComponent) SetParent(parent *BaseComponent) *BaseComponent {
	c.Parent = parent
	return c
}

// Get array of child components
func (c *BaseComponent) GetChildren() []Component {
	return c.Children
}

// Get array of child components belonging to the parent component
func (c *BaseComponent) GetSiblings() []Component {
	if c.GetParent() != nil {
		return c.GetParent().GetChildren()
	}

	return []Component{}
}

// Get parent component
func (c *BaseComponent) GetParent() Component {
	return c.Parent
}

// Get root component
func (c *BaseComponent) GetComponent() Component {
	return c
}

// Remove all children
func (c *BaseComponent) ClearChildren() {
	for _, child := range c.GetChildren() {
		child.Delete()
	}

	c.Children = []Component{}
}

// Perform cleanup on elements and its children
func (c *BaseComponent) Delete() {
	for _, child := range c.GetChildren() {
		child.Delete()
	}
	Global.DeleteElement(c.GlobalIndex)
}

// Set the global index for the component
func (c *BaseComponent) SetGlobalIndex(i int) {
	c.GlobalIndex = i
}

//#endregion Component relation

// #region Component graphical properties

// Set the nest depth of a component
func (c *BaseComponent) SetDepth(v int) *BaseComponent {
	c.Depth = v
	return c
}

// Get the rect of the component. (X,Y,Width,Height)
func (c *BaseComponent) GetRect() (int, int, int, int) {
	return c.GetX(), c.GetY(), c.GetWidth(), c.GetHeight()
}

// Get X coordinates
func (c *BaseComponent) GetX() int {
	if c.Centered && c.Absolute {
		pw := c.GetParent().GetInnerWidth()
		w := c.GetWidth()
		return (pw - w) / 2
	}
	return c.x
}

// Get Y coordinates
func (c *BaseComponent) GetY() int {
	if c.Centered && c.Absolute {
		pw := c.GetParent().GetInnerHeight()
		w := c.GetHeight()
		return (pw - w) / 2
	}
	return c.y
}

// Get coordinates of the component. (X,Y)
func (c *BaseComponent) GetPos() (int, int) {
	return c.GetX(), c.GetY()
}

// Set X coordinate of the component
func (c *BaseComponent) SetX(v int) *BaseComponent {
	c.x = v
	return c
}

// Set Y coordinate of the component
func (c *BaseComponent) SetY(v int) *BaseComponent {
	c.y = v
	return c
}

// Set coordinates of the component
func (c *BaseComponent) SetPos(x int, y int) *BaseComponent {
	c.SetX(x)
	c.SetY(y)
	return c
}

// Set width of the component
func (c *BaseComponent) SetW(v int) *BaseComponent {
	if c.Width != v {
		c.DispatchEvent("onResize")
	}
	c.Width = v
	return c
}

// Set height of the component
func (c *BaseComponent) SetH(v int) *BaseComponent {
	if c.Height != v {
		c.DispatchEvent("onResize")
	}
	c.Height = v
	return c
}

// Set if component's height should be equal to the parent's inner height
func (c *BaseComponent) SetHeightInherit(v bool) *BaseComponent {

	if c.InheritHeight != v {
		c.DispatchEvent("onResize")
	}
	c.InheritHeight = v
	return c
}

// Set if component's width should be equal to the parent's inner width
func (c *BaseComponent) SetWidthInherit(v bool) *BaseComponent {
	if c.InheritWidth != v {
		c.DispatchEvent("onResize")
	}
	c.InheritWidth = v
	return c
}

// Set size of the component
func (c *BaseComponent) SetSize(w int, h int) *BaseComponent {
	c.SetW(w)
	c.SetH(h)
	return c
}

// Get width of a component. if InheritWidth is true, retrieve parent's inner width
func (c *BaseComponent) GetWidth() int {
	if c.InheritWidth == true && c.GetParent() != nil {
		return c.GetParent().GetInnerWidth()
	}
	return c.Width
}

// Get height of a component. if InheritHeight is true, retrieve parent's inner height
func (c *BaseComponent) GetHeight() int {
	if c.InheritHeight == true && c.GetParent() != nil {
		return c.GetParent().GetInnerHeight()
	}
	return c.Height
}

// Get inner width of a component. (width - side paddings).
func (c *BaseComponent) GetInnerWidth() int {
	if c.ShowBorder {
		_, _, l, r := c.GetPaddings()
		return c.GetWidth() - (l + r)
	}
	return c.GetWidth()
}

// Get inner height of a component. (height - top and bottom paddings)
func (c *BaseComponent) GetInnerHeight() int {
	if c.ShowBorder {
		t, b, _, _ := c.GetPaddings()
		return c.GetHeight() - (t + b)
	}
	return c.GetHeight()
}

// Check if element's position should be dictated by parent element
func (c *BaseComponent) IsAbsolute() bool {
	return c.Absolute
}

// Get smallest area that can fit all children. (width,height)
func (c *BaseComponent) GetContentsSize() (int, int) {
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

// Get smallest height that can fit all children.
func (c *BaseComponent) GetContentsHeight() int {
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

// Get the width of the padding for each sides. (top,bottom,left,right)
func (c *BaseComponent) GetPaddings() (int, int, int, int) {
	top := c.PaddingTop
	bottom := c.PaddingBottom
	left := c.PaddingLeft
	right := c.PaddingRight

	if c.ShowBorder {
		if c.ShowTopBorder {
			top = max(1, top)
		}

		if c.ShowBottomBorder {
			bottom = max(1, bottom)
		}

		if c.ShowRightBorder {
			right = max(1, right)
		}

		if c.ShowLeftBorder {
			left = max(1, left)
		}
	}
	return top, bottom, left, right
}

// Set a flag to check if a component is a child of a flex component
func (c *BaseComponent) SetIsFlexItem(flex bool) *BaseComponent {
	c.IsFlexItem = flex
	return c
}

// Check if a component is a child of a flex component
func (c *BaseComponent) GetIsFlexItem() bool {
	return c.IsFlexItem
}

func (c *BaseComponent) SetAbsolute(v bool) *BaseComponent {
	c.Absolute = v
	return c
}

func (c *BaseComponent) SetCentered(v bool) *BaseComponent {
	c.Centered = v
	return c
}

func (c *BaseComponent) GetCentered() bool {
	return c.Centered
}

// #endregion Component graphical properties

// #region Event handler

// Perform bubbletea's update event on all elements that has focus
func (c *BaseComponent) Update() {
	for _, child := range c.GetChildren() {
		child.Update()
	}

	if c.GetFocusState() {
		c.DispatchEvent("onUpdate")
	}
}

// Hook a callback to a specific event
func (c *BaseComponent) AddEventListener(event string, cb func(tea.Msg), bubble bool) {
	list := c.GetEventCallbacks(event)
	list[uuid.New().String()] = EventCb{
		Cb:     cb,
		Bubble: bubble,
	}
	c.EventCallbacks[event] = list
}

// Queue all functions hooked to an event to be executed at the end of the frame
func (c *BaseComponent) DispatchEvent(event string) {
	var bubble []Component
	pt := c.GetParent()
	bubble = append(bubble, c)

	for pt != nil {
		bubble = append(bubble, pt)
		pt = pt.GetParent()
	}

	for _, element := range bubble {
		continueBubble := true
		for callbackUUID, cb := range element.GetEventCallbacks(event) {
			if !cb.Bubble {
				continueBubble = false
			}
			Global.AddEventCallback(event, element.GetUUID(), callbackUUID, cb.Cb)
		}
		if !continueBubble {
			break
		}
	}
}

// Called on all function right before rendering
func (c *BaseComponent) BeforeRender() {
	for _, c := range c.GetChildren() {
		c.BeforeRender()
	}
}

// Get a list of callbacks hooked to a specific event
func (c *BaseComponent) GetEventCallbacks(event string) map[string]EventCb {
	if c.EventCallbacks[event] == nil {
		c.EventCallbacks[event] = map[string]EventCb{}
	}
	return c.EventCallbacks[event]
}

//#endregion Event handler

// #region Component non graphical properties

// Set focus on a component.
// Only the component with focus receives bubbletea's event.
func (c *BaseComponent) Focus() {
	Global.BlurAll()
	c.Focused = true
	c.DispatchEvent("onFocus")
	c.DispatchEvent("onFocusChange")
}

// Remove focus from a component
func (c *BaseComponent) Blur() {
	for _, child := range c.GetChildren() {
		child.Blur()
	}
	c.Focused = false
	c.DispatchEvent("onBlur")
	c.DispatchEvent("onFocusChange")
}

// Check if a component is focused
func (c *BaseComponent) GetFocusState() bool {
	return c.Focused
}

// Get a name of an individual component
func (c *BaseComponent) GetName() string {
	return c.Name
}

// Get a name of the component type
func (c *BaseComponent) GetComponentName() string {
	return c.ComponentName
}

// Set name of a component
func (c *BaseComponent) SetName(n string) *BaseComponent {
	c.Name = n
	return c
}

// Set name of a component type
func (c *BaseComponent) SetComponentName(n string) *BaseComponent {
	c.ComponentName = n
	return c
}

// Set title of a component.
// Title is displayed at the top of the component with border.
func (c *BaseComponent) SetTitle(str string) *BaseComponent {
	c.Title = str
	return c
}

// Get the title of a component
func (c *BaseComponent) GetTitle() string {
	return c.Title
}

// Get UUID of a component
func (c *BaseComponent) GetUUID() string {
	return c.UUID
}

func (c *BaseComponent) SetFlexProportion(v float32) *BaseComponent {
	c.FlexProportion = v
	return c
}
func (c *BaseComponent) GetFlexProportion() float32 {
	return c.FlexProportion
}
func (c *BaseComponent) SetMinWidth(v int) *BaseComponent {
	c.MinWidth = v
	return c
}
func (c *BaseComponent) GetMinWidth() int {
	return c.MinWidth
}
func (c *BaseComponent) SetMinHeight(v int) *BaseComponent {
	c.MinHeight = v
	return c
}
func (c *BaseComponent) GetMinHeight() int {
	return c.MinHeight
}

// #endregion Component non graphical properties

// #region Rendering

// Perform rendering for a component and all its child components.
// Rendered result is written to the Canvas property
func (b *BaseComponent) RenderToCanvas() {
	var result, fg, bg = b.CreateCanvas()
	if !b.Visibility {
		b.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}

	top, _, left, _ := b.GetPaddings()
	cursor := top

	innerWidth := b.GetInnerWidth() + 1

	for _, c := range b.GetChildren() {
		childHeight := c.GetHeight()
		childWidth := c.GetWidth()
		if !c.GetVisibility() {
			cursor += childHeight
			continue
		}
		c.RenderToCanvas()
		output, childFG, childBG := c.GetCanvas()

		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := left + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + top
				for index, char := range line {
					result[posY][globalX+index] = char

					if len(childFG[posY][index]) > 0 {
						fg[posY][globalX+index] = childFG[posY][index]
					}

					if len(childBG[posY][index]) > 0 {
						bg[posY][globalX+index] = childBG[posY][index]
					}
				}
			}
		} else {
			// Loop through lines
			pt := cursor
			for y := 0; y < min(childHeight, len(result), len(output)); y++ {
				line := output[y]
				if len(result) <= pt {
					break
				}
				// Loop through characters
				for x := range min(childWidth, innerWidth-left, len(line), len(result[pt])) {
					// If canvas is smaller than the horizontal pointer, break
					char := line[x]
					// Check if the character is over the drawable area
					result[pt][x+left] = char
					if len(childFG[y][x]) > 0 {
						fg[pt][x+left] = childFG[y][x]
					}
					if len(childBG[y][x]) > 0 {
						bg[pt][x+left] = childBG[y][x]
					}
				}
				pt++
			}
			cursor += childHeight
		}
	}

	b.SetCanvas(result, fg, bg)
}

// Create a 2D array of string the size of component
func (c *BaseComponent) CreateCanvas() ([][]string, [][]string, [][]string) {
	height := max(c.GetContentsHeight()+1, 1)
	width := max(c.GetWidth(), 1)

	var arr [][]string = make([][]string, height)
	var fg [][]string = make([][]string, height)
	var bg [][]string = make([][]string, height)

	for i := range height {
		arr[i] = make([]string, width)
		fg[i] = make([]string, width)
		bg[i] = make([]string, width)
		for a := range width {
			arr[i][a] = " "
			fg[i][a] = c.Foreground
			bg[i][a] = c.Background
		}
	}

	return arr, fg, bg
}

// Get the rendered canvas
func (c *BaseComponent) GetCanvas() ([][]string, [][]string, [][]string) {
	return c.Canvas, c.FGSheet, c.BGSheet
}

// Get the rendered canvas
func (c *BaseComponent) SetCanvas(
	canvas [][]string,
	fg [][]string,
	bg [][]string,
) {
	c.Canvas = canvas
	c.FGSheet = fg
	c.BGSheet = bg
}

// Add border to a component if applicable
func (c *BaseComponent) addBorder(arr [][]string, fg [][]string, bg [][]string) ([][]string, [][]string, [][]string) {
	if !c.ShowBorder || len(arr) <= 1 || len(arr[0]) <= 1 {
		return arr, fg, bg
	}

	foreground := c.BorderForeground
	side := (helper.Dictionary(helper.BorderSide))
	top := (helper.Dictionary(helper.BorderTop))
	tl := (helper.Dictionary(helper.BorderTopLeft))
	tr := (helper.Dictionary(helper.BorderTopRight))
	bl := (helper.Dictionary(helper.BorderBottomLeft))
	br := (helper.Dictionary(helper.BorderBottomRight))

	if c.GetFocusState() || c.GetDoubleBorder() {
		foreground = c.BorderFocusForeground
		side = (helper.Dictionary(helper.BorderSideDouble))
		top = (helper.Dictionary(helper.BorderTopDouble))
		tl = (helper.Dictionary(helper.BorderTopLeftDouble))
		tr = (helper.Dictionary(helper.BorderTopRightDouble))
		bl = (helper.Dictionary(helper.BorderBottomLeftDouble))
		br = (helper.Dictionary(helper.BorderBottomRightDouble))
	}

	wid := c.GetWidth()
	hei := c.GetHeight()

	if len(arr) > 0 && len(arr[0]) > 0 {

		for i := range c.GetWidth() {
			if c.ShowTopBorder {
				arr[0][i] = top
				fg[0][i] = foreground
			}
			if c.ShowBottomBorder {
				arr[hei-1][i] = top
				fg[hei-1][i] = foreground
			}
		}

		for i := range c.GetHeight() {
			if c.ShowLeftBorder {
				arr[i][0] = side
				fg[i][0] = foreground
			}
			if c.ShowRightBorder {
				arr[i][wid-1] = side
				fg[i][wid-1] = foreground
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

			i := 0
			var char string
			state := -1
			buf := title
			length := uniseg.StringWidth(title)
			for i < length && len(buf) > 0 {
				var clusterWidth int
				char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
				arr[0][i+1] = char
				i += clusterWidth
			}

		case "center":
			i := 0
			var char string
			state := -1
			buf := title
			length := uniseg.StringWidth(title)
			for i < wid-1 && i < length && len(buf) > 0 {
				var clusterWidth int
				char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
				arr[0][i+max(1, (wid-length)/2)] = char
				i += clusterWidth
			}

		case "right":
			i := min(uniseg.StringWidth(title), wid-2)
			var char string
			state := -1
			buf := title
			for i > 0 && len(buf) > 0 {
				var clusterWidth int
				char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
				arr[0][wid-(i+2)] = char
				i -= clusterWidth
			}
		}

		for key, str := range c.BorderLabels {
			if len(str) == 0 {
				continue
			}
			switch key {
			case "TopLeft":
				i := 0
				var char string
				state := -1
				buf := str
				length := uniseg.StringWidth(str)
				for i < length && len(buf) > 0 {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i+clusterWidth < wid {
						arr[0][i+1] = char
						fg[0][i+1] = c.BorderLabelColors[key]
					}
					i += clusterWidth
				}

			case "Top":
				i := 0
				var char string
				state := -1
				buf := str
				length := uniseg.StringWidth(str)
				half := max(1, (wid-length)/2)
				for i < wid-1 && i < length && len(buf) > 0 {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i+clusterWidth+half < wid {
						arr[0][i+max(1, (wid-length)/2)] = char
						fg[0][i+max(1, (wid-length)/2)] = c.BorderLabelColors[key]
					}
					i += clusterWidth
				}

			case "TopRight":
				i := min(uniseg.StringWidth(str), wid-2)
				var char string
				state := -1
				buf := str
				for i > 0 && len(buf) > 0 {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i-clusterWidth >= 0 {
						arr[0][wid-(i+2)] = char
						fg[0][wid-(i+2)] = c.BorderLabelColors[key]
					}
					i -= clusterWidth
				}

			case "BottomLeft":
				i := 0
				var char string
				state := -1
				buf := str
				length := uniseg.StringWidth(str)
				for i < length && len(buf) > 0 && i < wid {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i+clusterWidth < wid {
						arr[hei-1][i+1] = char
						fg[hei-1][i+1] = c.BorderLabelColors[key]
					}
					i += clusterWidth
				}

			case "Bottom":

				i := 0
				var char string
				state := -1
				buf := str
				length := uniseg.StringWidth(str)
				half := max(1, (wid-length)/2)
				for i < wid-1 && i < length && len(buf) > 0 {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i+clusterWidth+half < wid {
						arr[hei-1][i+half] = char
						fg[hei-1][i+half] = c.BorderLabelColors[key]
					}
					i += clusterWidth
				}

			case "BottomRight":

				i := min(uniseg.StringWidth(str), wid-2)
				var char string
				state := -1
				buf := str
				for i > 0 && len(buf) > 0 {
					var clusterWidth int
					char, buf, clusterWidth, state = uniseg.FirstGraphemeClusterInString(buf, state)
					if i-clusterWidth >= 0 {
						arr[hei-1][wid-(i+2)] = char
						fg[hei-1][wid-(i+2)] = c.BorderLabelColors[key]
					}
					i -= clusterWidth
				}
			}

		}

	}

	return arr, fg, bg
}

// Set if border should be visible
func (c *BaseComponent) SetBorder(show bool) *BaseComponent {
	c.ShowBorder = show
	return c
}
func (c *BaseComponent) SetBorderForeground(v string) *BaseComponent {
	c.BorderForeground = v
	return c
}
func (c *BaseComponent) ClearBorderForeground() *BaseComponent {
	c.BorderForeground = ""
	return c
}
func (c *BaseComponent) SetBorderFocusForeground(v string) *BaseComponent {
	c.BorderFocusForeground = v
	return c
}
func (c *BaseComponent) ClearBorderFocusForeground() *BaseComponent {
	c.BorderFocusForeground = ""
	return c
}

// Set where should the title be rendered
func (c *BaseComponent) SetTitleAlignment(str string) *BaseComponent {
	c.Title = str
	return c
}

// Get where should the title be rendered
func (c *BaseComponent) GetTitleAlignment() string {
	return c.TitleAlignment
}

// Set Visibility for each border edges
func (c *BaseComponent) SetBorders(top bool, bottom bool, left bool, right bool) *BaseComponent {
	c.ShowTopBorder = top
	c.ShowBottomBorder = bottom
	c.ShowLeftBorder = left
	c.ShowRightBorder = right
	return c
}

// Set if rounded border corder should be rendered
func (c *BaseComponent) SetBorderCorner(show bool) *BaseComponent {
	c.ShowBorderCorner = show
	return c
}

// Set padding width of the border
func (c *BaseComponent) SetPadding(v int) *BaseComponent {
	c.PaddingTop = v
	c.PaddingBottom = v
	c.PaddingLeft = v
	c.PaddingRight = v
	return c
}
func (c *BaseComponent) SetPaddings(t int, b int, l int, r int) *BaseComponent {
	c.PaddingTop = t
	c.PaddingBottom = b
	c.PaddingLeft = l
	c.PaddingRight = r
	return c
}

// Set if rendered border should be double border
func (c *BaseComponent) SetDoubleBorder(v bool) *BaseComponent {
	c.ShowDoubleBorder = v
	return c
}

// Get if rendered border should be double border
func (c *BaseComponent) GetDoubleBorder() bool {
	return c.ShowDoubleBorder
}

// Set if a component should be rendered
func (c *BaseComponent) SetVisibility(v bool) *BaseComponent {
	c.Visibility = v
	return c
}

func (c *BaseComponent) GetVisibility() bool {
	return c.Visibility
}

// Set a label to be displayed on corners of the border
func (c *BaseComponent) SetBorderLabel(key string, str string) {
	c.BorderLabels[key] = str
}

// Set a label to be displayed on corners of the border
func (c *BaseComponent) SetBorderLabelColor(key string, str string) {
	c.BorderLabelColors[key] = str
}

// Recursively hide all elements invisible to the parent element
func (c *BaseComponent) UpdateVisibility(ytop int, hei int) {
	top := 0
	y := ytop
	h := hei
	for _, child := range c.GetChildren() {
		childHeight := child.GetHeight()
		child.SetVisibility(!(top+childHeight < y || top > y+h))
		top += childHeight
	}
}

func (c *BaseComponent) SetBackgroundGradient(v []color.Color) *BaseComponent {
	c.BackgroundGradient = v
	return c
}
func (c *BaseComponent) SetForegroundGradient(v []color.Color) *BaseComponent {
	c.ForegroundGradient = v
	return c
}
func (c *BaseComponent) GetBackgroundGradient() []color.Color {
	return c.BackgroundGradient
}
func (c *BaseComponent) GetForegroundGradient() []color.Color {
	return c.ForegroundGradient
}
func (c *BaseComponent) ClearBackgroundGradient() *BaseComponent {
	c.BackgroundGradient = []color.Color{}
	return c
}
func (c *BaseComponent) ClearForegroundGradient() *BaseComponent {
	c.ForegroundGradient = []color.Color{}
	return c
}
func (c *BaseComponent) SetForeground(v string) *BaseComponent {
	c.Foreground = v
	return c
}
func (c *BaseComponent) SetBackground(v string) *BaseComponent {
	c.Background = v
	return c
}
func (c *BaseComponent) ClearForeground() *BaseComponent {
	c.Foreground = ""
	return c
}
func (c *BaseComponent) ClearBackground() *BaseComponent {
	c.Background = ""
	return c
}

func (c *BaseComponent) GetForeground() string {
	return c.Foreground
}
func (c *BaseComponent) GetBackground() string {
	return c.Background
}

//#endregion Rendering

// #region Debugging

// Retrieve the list of parents
func (c *BaseComponent) Trace(list []string) []string {

	if c.GetParent() != nil {
		list = append(list, c.GetParent().Trace(list)...)
	}

	list = append(list, strconv.Itoa(c.Depth)+":"+c.GetComponentName()+"("+c.GetName()+")")
	return list
}

// Retrieve the list of parents
func (c *BaseComponent) GetTrace() []string {
	return c.Trace([]string{})
}

// String representation of a component for debugging purpose.
func (c *BaseComponent) ToString() string {
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

//#endregion Debugging
