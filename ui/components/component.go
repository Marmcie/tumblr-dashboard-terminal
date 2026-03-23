package component

import (
	"bytes"
	"fmt"
	"image/color"
	"strconv"
	"tumblr-dt/ui/helper"

	tea "charm.land/bubbletea/v2"
	"github.com/google/uuid"
	"github.com/mattn/go-runewidth"
)

type Component interface {
	SetBorder(bool) *ComponentState
	SetBorderCorner(bool) *ComponentState
	GetBorderPadding() int
	SetBorderPadding(int) *ComponentState
	GetBorderPaddings() (int, int, int, int)
	SetBorders(bool, bool, bool, bool) *ComponentState
	GetCanvas() ([][]string, [][]string, [][]string)
	SetCanvas([][]string, [][]string, [][]string)
	AddChild(Component)
	GetChildren() []Component
	GetComponent() Component
	GetComponentName() string
	SetComponentName(string) *ComponentState
	GetContentsHeight() int
	GetContentsSize() (int, int)
	SetDepth(int) *ComponentState
	GetEventCallbacks(string) map[string]EventCb
	AddEventListener(string, func(tea.Msg, int), bool)
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
	SetBackgroundGradient([]color.Color) *ComponentState
	SetForegroundGradient([]color.Color) *ComponentState
	GetBackgroundGradient() []color.Color
	GetForegroundGradient() []color.Color
	ClearBackgroundGradient()
	ClearForegroundGradient()
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
	SetVisibility(bool) *ComponentState
	GetVisibility() bool
	ToString() string
	SetBorderLabel(string, string)
	UpdateVisibility(int, int)
	Delete()
	SetGlobalIndex(int)
	SetAbsolute(bool) *ComponentState
	SetCentered(bool) *ComponentState
	GetCentered() bool
	SetForeground(string)
	SetBackground(string)
	ClearForeground()
	ClearBackground()
	GetForeground() string
	GetBackground() string
}

// Base class for all components
type ComponentState struct {
	// X coordinates
	x int
	// Y coordinates
	y             int
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

	ShowBorder     bool
	BorderPadWidth int
	// Name of an individual component
	Name string
	// Name of a component type
	ComponentName    string
	EventCallbacks   map[string]map[string]EventCb
	Absolute         bool
	Overflow         bool
	ShowTopBorder    bool
	ShowBottomBorder bool
	ShowLeftBorder   bool
	ShowRightBorder  bool
	ShowBorderCorner bool
	IsFlexItem       bool
	Title            string
	TitleAlignment   string
	ShowDoubleBorder bool
	Visibility       bool
	BorderLabels     map[string]string
	//Index of the component within global element list
	GlobalIndex        int
	BackgroundGradient []color.Color
	ForegroundGradient []color.Color
	Background         string
	Foreground         string
}

type EventCb struct {
	Cb     func(tea.Msg, int)
	Bubble bool
}

// Initialized all shared values
func (c *ComponentState) Initialize(name string) {
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
	c.EventCallbacks = map[string]map[string]EventCb{}
	c.TitleAlignment = "center"
	c.Foreground = ""
	c.Background = ""

	c.BorderLabels = map[string]string{
		"TopLeft":     "",
		"Top":         "",
		"TopRight":    "",
		"BottomRight": "",
		"Bottom":      "",
		"BottomLeft":  "",
	}

	c.SetVisibility(true)

	c.ShowDoubleBorder = false
	c.GlobalIndex = Global.AddElement(c)
}

// #region Component relation

// Adds a child to an component
func (c *ComponentState) AddChild(child Component) {
	child.SetDepth(c.Depth + 1)
	child.SetParent(c)
	c.Children = append(c.Children, child)
	c.DispatchEvent("onAddChild")
}

// Set parent of a component
func (c *ComponentState) SetParent(parent *ComponentState) *ComponentState {
	c.Parent = parent
	return c
}

// Get array of child components
func (c *ComponentState) GetChildren() []Component {
	return c.Children
}

// Get array of child components belonging to the parent component
func (c *ComponentState) GetSiblings() []Component {
	if c.GetParent() != nil {
		return c.GetParent().GetChildren()
	}

	return []Component{}
}

// Get parent component
func (c *ComponentState) GetParent() Component {
	return c.Parent
}

// Get root component
func (c *ComponentState) GetComponent() Component {
	return c
}

// Remove all children
func (c *ComponentState) ClearChildren() {
	for _, child := range c.GetChildren() {
		child.Delete()
	}

	c.Children = []Component{}
}

// Perform cleanup on elements and its children
func (c *ComponentState) Delete() {
	for _, child := range c.GetChildren() {
		child.Delete()
	}
	Global.DeleteElement(c.GlobalIndex)
}

// Set the global index for the component
func (c *ComponentState) SetGlobalIndex(i int) {
	c.GlobalIndex = i
}

//#endregion Component relation

// #region Component graphical properties

// Set the nest depth of a component
func (c *ComponentState) SetDepth(v int) *ComponentState {
	c.Depth = v
	return c
}

// Get the rect of the component. (X,Y,Width,Height)
func (c *ComponentState) GetRect() (int, int, int, int) {
	return c.GetX(), c.GetY(), c.GetWidth(), c.GetHeight()
}

// Get X coordinates
func (c *ComponentState) GetX() int {
	if c.Centered && c.Absolute {
		pw := c.GetParent().GetInnerWidth()
		w := c.GetWidth()
		return (pw - w) / 2
	}
	return c.x
}

// Get Y coordinates
func (c *ComponentState) GetY() int {
	if c.Centered && c.Absolute {
		pw := c.GetParent().GetInnerHeight()
		w := c.GetHeight()
		return (pw - w) / 2
	}
	return c.y
}

// Get coordinates of the component. (X,Y)
func (c *ComponentState) GetPos() (int, int) {
	return c.GetX(), c.GetY()
}

// Set X coordinate of the component
func (c *ComponentState) SetX(v int) *ComponentState {
	c.x = v
	return c
}

// Set Y coordinate of the component
func (c *ComponentState) SetY(v int) *ComponentState {
	c.y = v
	return c
}

// Set coordinates of the component
func (c *ComponentState) SetPos(x int, y int) *ComponentState {
	c.SetX(x)
	c.SetY(y)
	return c
}

// Set width of the component
func (c *ComponentState) SetW(v int) *ComponentState {
	c.Width = v
	return c
}

// Set height of the component
func (c *ComponentState) SetH(v int) *ComponentState {
	c.Height = v
	return c
}

// Set if component's height should be equal to the parent's inner height
func (c *ComponentState) SetHeightInherit(v bool) *ComponentState {
	c.InheritHeight = v
	return c
}

// Set if component's width should be equal to the parent's inner width
func (c *ComponentState) SetWidthInherit(v bool) *ComponentState {
	c.InheritWidth = v
	return c
}

// Set size of the component
func (c *ComponentState) SetSize(w int, h int) *ComponentState {
	c.SetW(w)
	c.SetH(h)
	return c
}

// Get width of a component. if InheritWidth is true, retrieve parent's inner width
func (c *ComponentState) GetWidth() int {
	if c.InheritWidth == true && c.GetParent() != nil {
		return c.GetParent().GetInnerWidth()
	}
	return c.Width
}

// Get height of a component. if InheritHeight is true, retrieve parent's inner height
func (c *ComponentState) GetHeight() int {
	if c.InheritHeight == true && c.GetParent() != nil {
		return c.GetParent().GetInnerHeight()
	}
	return c.Height
}

// Get inner width of a component. (width - side paddings).
func (c *ComponentState) GetInnerWidth() int {
	if c.ShowBorder {
		_, _, l, r := c.GetBorderPaddings()
		return c.GetWidth() - (l + r)
	}
	return c.GetWidth()
}

// Get inner height of a component. (height - top and bottom paddings)
func (c *ComponentState) GetInnerHeight() int {
	if c.ShowBorder {
		t, b, _, _ := c.GetBorderPaddings()
		return c.GetHeight() - (t + b)
	}
	return c.GetHeight()
}

// Check if element's position should be dictated by parent element
func (c *ComponentState) IsAbsolute() bool {
	return c.Absolute
}

// Get smallest area that can fit all children. (width,height)
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

// Get smallest height that can fit all children.
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

// Get the width of the padding
func (c *ComponentState) GetBorderPadding() int {
	if c.ShowBorder {
		return c.BorderPadWidth
	}
	return 0
}

// Get the width of the padding for each sides. (top,bottom,left,right)
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

// Set a flag to check if a component is a child of a flex component
func (c *ComponentState) SetIsFlexItem(flex bool) *ComponentState {
	c.IsFlexItem = flex
	return c
}

// Check if a component is a child of a flex component
func (c *ComponentState) GetIsFlexItem() bool {
	return c.IsFlexItem
}

func (c *ComponentState) SetAbsolute(v bool) *ComponentState {
	c.Absolute = v
	return c
}

func (c *ComponentState) SetCentered(v bool) *ComponentState {
	c.Centered = v
	return c
}

func (c *ComponentState) GetCentered() bool {
	return c.Centered
}

// #endregion Component graphical properties

// #region Event handler

// Perform bubbletea's update event on all elements that has focus
func (c *ComponentState) Update() {
	for _, child := range c.GetChildren() {
		child.Update()
	}

	if c.GetFocusState() {
		c.DispatchEvent("onUpdate")
	}
}

// Hook a callback to a specific event
func (c *ComponentState) AddEventListener(event string, cb func(tea.Msg, int), bubble bool) {
	list := c.GetEventCallbacks(event)
	list[uuid.New().String()] = EventCb{
		Cb:     cb,
		Bubble: bubble,
	}
	c.EventCallbacks[event] = list
}

// Queue all functions hooked to an event to be executed at the end of the frame
func (c *ComponentState) DispatchEvent(event string) {
	var bubble []Component
	pt := c.GetParent()

	for pt != nil {
		bubble = append(bubble, pt)
		pt = pt.GetParent()
	}

	for _, element := range bubble {
		for callbackUUID, cb := range element.GetEventCallbacks(event) {
			if cb.Bubble {
				Global.AddEventCallback(event, element.GetUUID(), callbackUUID, cb.Cb)
			}
		}
	}

	for callbackUUID, cb := range c.GetEventCallbacks(event) {
		Global.AddEventCallback(event, c.GetUUID(), callbackUUID, cb.Cb)
	}
}

// Called on all function right before rendering
func (c *ComponentState) Propagate() {
	for _, c := range c.GetChildren() {
		c.Propagate()
	}
}

// Get a list of callbacks hooked to a specific event
func (c *ComponentState) GetEventCallbacks(event string) map[string]EventCb {
	if c.EventCallbacks[event] == nil {
		c.EventCallbacks[event] = map[string]EventCb{}
	}
	return c.EventCallbacks[event]
}

//#endregion Event handler

// #region Component non graphical properties

// Set focus on a component.
// Only the component with focus receives bubbletea's event.
func (c *ComponentState) Focus() {
	Global.BlurAll()
	c.Focused = true
	c.DispatchEvent("onFocus")
	c.DispatchEvent("onFocusChange")
}

// Remove focus from a component
func (c *ComponentState) Blur() {
	for _, child := range c.GetChildren() {
		child.Blur()
	}
	c.Focused = false
	c.DispatchEvent("onBlur")
	c.DispatchEvent("onFocusChange")
}

// Check if a component is focused
func (c *ComponentState) GetFocusState() bool {
	return c.Focused
}

// Get a name of an individual component
func (c *ComponentState) GetName() string {
	return c.Name
}

// Get a name of the component type
func (c *ComponentState) GetComponentName() string {
	return c.ComponentName
}

// Set name of a component
func (c *ComponentState) SetName(n string) *ComponentState {
	c.Name = n
	return c
}

// Set name of a component type
func (c *ComponentState) SetComponentName(n string) *ComponentState {
	c.ComponentName = n
	return c
}

// Set title of a component.
// Title is displayed at the top of the component with border.
func (c *ComponentState) SetTitle(str string) *ComponentState {
	c.Title = str
	return c
}

// Get the title of a component
func (c *ComponentState) GetTitle() string {
	return c.Title
}

// Get UUID of a component
func (c *ComponentState) GetUUID() string {
	return c.UUID
}

// #endregion Component non graphical properties

// #region Rendering

// Perform rendering for a component and all its child components.
// Rendered result is written to the Canvas property
func (b *ComponentState) PrepareFrame() {
	var result, fg, bg = b.CreateCanvas()
	if !b.Visibility {
		b.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}

	top, _, left, _ := b.GetBorderPaddings()
	cursor := top

	innerWidth := b.GetInnerWidth() + 1

	for _, c := range b.GetChildren() {
		childHeight := c.GetHeight()
		childWidth := c.GetWidth()
		if !c.GetVisibility() {
			cursor += childHeight
			continue
		}
		c.PrepareFrame()
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
				// Loop through characters
				for x := range min(childWidth, innerWidth-left, len(line)) {
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
func (c *ComponentState) CreateCanvas() ([][]string, [][]string, [][]string) {
	height := c.GetContentsHeight() + 1
	width := c.GetWidth()

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
func (c *ComponentState) GetCanvas() ([][]string, [][]string, [][]string) {
	return c.Canvas, c.FGSheet, c.BGSheet
}

// Get the rendered canvas
func (c *ComponentState) SetCanvas(
	canvas [][]string,
	fg [][]string,
	bg [][]string,
) {
	c.Canvas = canvas
	c.FGSheet = fg
	c.BGSheet = bg
}

// Add border to a component if applicable
func (c *ComponentState) addBorder(arr [][]string) [][]string {
	if !c.ShowBorder || c.GetBorderPadding() == 0 || len(arr) < 3 || len(arr[0]) < 3 {
		return arr
	}

	side := (helper.Dictionary(helper.BorderSide))
	top := (helper.Dictionary(helper.BorderTop))
	tl := (helper.Dictionary(helper.BorderTopLeft))
	tr := (helper.Dictionary(helper.BorderTopRight))
	bl := (helper.Dictionary(helper.BorderBottomLeft))
	br := (helper.Dictionary(helper.BorderBottomRight))

	if c.GetFocusState() || c.GetDoubleBorder() {
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
				arr[0][i+1] = (string(char))
			}

		case "center":
			length := len(title)
			for i := range min(wid-1, runewidth.StringWidth(title)) {
				char := title[i]
				arr[0][i+max(1, (wid-length)/2)] = (string(char))
			}

		case "right":
			strWidth := len(title)
			for i := 0; i < min(wid-2, strWidth); i++ {
				char := title[strWidth-(i+1)]
				arr[0][wid-(i+2)] = (string(char))
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
					arr[0][i+1] = (string(char))
				}

			case "Top":
				length := len(str)
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[0][i+max(1, (wid-length)/2)] = (string(char))
				}

			case "TopRight":
				strWidth := len(str)
				for i := 0; i < min(wid-2, strWidth); i++ {
					char := str[strWidth-(i+1)]
					arr[0][wid-(i+2)] = (string(char))
				}

			case "BottomLeft":
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[hei-1][i+1] = (string(char))
				}

			case "Bottom":
				length := len(str)
				center := wid - length/2
				for i := range min(wid-1, runewidth.StringWidth(str)) {
					char := str[i]
					arr[hei-1][i+max(1, center)] = (string(char))
				}

			case "BottomRight":
				strWidth := len(str)
				for i := 0; i < min(wid-2, strWidth); i++ {
					char := str[strWidth-(i+1)]
					arr[hei-1][wid-(i+2)] = (string(char))
				}
			}

		}

	}

	return arr
}

// Set if border should be visible
func (c *ComponentState) SetBorder(show bool) *ComponentState {
	c.ShowBorder = show
	if show && c.BorderPadWidth == 0 {
		c.SetBorderPadding(1)
	}
	return c
}

// Set where should the title be rendered
func (c *ComponentState) SetTitleAlignment(str string) *ComponentState {
	c.Title = str
	return c
}

// Get where should the title be rendered
func (c *ComponentState) GetTitleAlignment() string {
	return c.TitleAlignment
}

// Set Visibility for each border edges
func (c *ComponentState) SetBorders(top bool, bottom bool, left bool, right bool) *ComponentState {
	c.ShowTopBorder = top
	c.ShowBottomBorder = bottom
	c.ShowLeftBorder = left
	c.ShowRightBorder = right
	return c
}

// Set if rounded border corder should be rendered
func (c *ComponentState) SetBorderCorner(show bool) *ComponentState {
	c.ShowBorderCorner = show
	return c
}

// Set padding width of the border
func (c *ComponentState) SetBorderPadding(v int) *ComponentState {
	c.BorderPadWidth = v
	return c
}

// Set if rendered border should be double border
func (c *ComponentState) SetDoubleBorder(v bool) *ComponentState {
	c.ShowDoubleBorder = v
	return c
}

// Get if rendered border should be double border
func (c *ComponentState) GetDoubleBorder() bool {
	return c.ShowDoubleBorder
}

// Set if a component should be rendered
func (c *ComponentState) SetVisibility(v bool) *ComponentState {
	c.Visibility = v
	return c
}

func (c *ComponentState) GetVisibility() bool {
	return c.Visibility
}

// Set a label to be displayed on corners of the border
func (c *ComponentState) SetBorderLabel(key string, str string) {
	c.BorderLabels[key] = str
}

// Recursively hide all elements invisible to the parent element
func (c *ComponentState) UpdateVisibility(ytop int, hei int) {
	top := 0
	y := ytop
	h := hei
	hidden := 0
	for _, child := range c.GetChildren() {
		childHeight := child.GetHeight()
		child.SetVisibility(!(top+childHeight < y || top > y+h))
		if top+childHeight < y || top > y+h {
			hidden++
		}
		top += childHeight
	}
}

func (c *ComponentState) SetBackgroundGradient(v []color.Color) *ComponentState {
	c.BackgroundGradient = v
	return c
}
func (c *ComponentState) SetForegroundGradient(v []color.Color) *ComponentState {
	c.ForegroundGradient = v
	return c
}
func (c *ComponentState) GetBackgroundGradient() []color.Color {
	return c.BackgroundGradient
}
func (c *ComponentState) GetForegroundGradient() []color.Color {
	return c.ForegroundGradient
}
func (c *ComponentState) ClearBackgroundGradient() {
	c.BackgroundGradient = []color.Color{}
}
func (c *ComponentState) ClearForegroundGradient() {
	c.ForegroundGradient = []color.Color{}
}
func (c *ComponentState) SetForeground(v string) {
	c.Foreground = v
}
func (c *ComponentState) SetBackground(v string) {
	c.Background = v
}
func (c *ComponentState) ClearForeground() {
	c.Foreground = ""
}
func (c *ComponentState) ClearBackground() {
	c.Background = ""
}

func (c *ComponentState) GetForeground() string {
	return c.Foreground
}
func (c *ComponentState) GetBackground() string {
	return c.Background
}

//#endregion Rendering

// #region Debugging

// Retrieve the list of parents
func (c *ComponentState) Trace(list []string) []string {

	if c.GetParent() != nil {
		list = append(list, c.GetParent().Trace(list)...)
	}

	list = append(list, strconv.Itoa(c.Depth)+":"+c.GetComponentName()+"("+c.GetName()+")")
	return list
}

// Retrieve the list of parents
func (c *ComponentState) GetTrace() []string {
	return c.Trace([]string{})
}

// String representation of a component for debugging purpose.
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

//#endregion Debugging
