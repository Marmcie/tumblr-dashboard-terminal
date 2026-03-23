package component

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Component that can scroll to show child elements
type Scrollable struct {
	ComponentState
	OffsetY     int
	OffsetX     int
	InnerHeight int
	InnerWidth  int
	ScrollX     bool
	ScrollY     bool
	Bottom      int
}

func NewScrollable(name string) *Scrollable {
	flex := &Scrollable{}
	flex.Initialize(name)
	flex.OffsetX = 0
	flex.OffsetY = 0
	flex.SetComponentName("Scrollable")
	flex.ScrollX = false
	flex.ScrollY = true
	flex.Bottom = 0

	flex.AddEventListener("onAddChild", func(msg tea.Msg, time int) {
		w, h := flex.GetContentsSize()
		flex.InnerHeight = h
		flex.InnerWidth = w
	})

	return flex
}

func (b *Scrollable) findBottom(canvas [][]string) {
	for i := len(canvas) - 1; i >= 0; i-- {
		line := canvas[i]
		if len(strings.ReplaceAll(strings.Join(line, ""), " ", "")) > 0 {
			b.Bottom = i
			return
		}
	}
}

func (c *Scrollable) CreateCanvas() ([][]string, [][]string, [][]string) {
	var arr [][]string
	var fg [][]string
	var bg [][]string
	height := c.GetHeight()
	width := c.GetWidth()

	// height := c.GetHeight()

	for range height {
		arr = append(arr, strings.Split(strings.Repeat(" ", width), ""))
		fg = append(fg, strings.Split(strings.Repeat(c.Foreground+",", width), ","))
		bg = append(bg, strings.Split(strings.Repeat(c.Background+",", width), ","))
	}

	return arr, fg, bg
}

// Returns Line per line contents,x,y
func (b *Scrollable) PrepareFrame() {

	if !b.Visibility {
		b.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}
	result, fg, bg := b.CreateCanvas()
	b.ComponentState.PrepareFrame()

	output, childFG, childBG := b.GetCanvas()
	boxHeight := b.GetInnerHeight()
	boxWidth := b.GetInnerWidth()

	b.findBottom(output)
	bottomEdge := b.OffsetY + boxHeight + 1

	for lineY := b.OffsetY; lineY < min(bottomEdge, len(output)); lineY++ {
		line := output[lineY]
		leftEdge := boxWidth + b.OffsetX + 1
		for lineX := b.OffsetX; lineX < min(len(line), leftEdge); lineX++ {
			char := line[lineX]
			result[lineY-b.OffsetY][lineX-b.OffsetX] = char
			if len(childFG[lineY][lineX]) > 0 {
				fg[lineY-b.OffsetY][lineX-b.OffsetX] = childFG[lineY][lineX]
			}
			if len(childBG[lineY][lineX]) > 0 &&lineY-b.OffsetY > 0 {
				bg[lineY-b.OffsetY][lineX-b.OffsetX] = childBG[lineY][lineX]
			}
		}

	}

	result = b.addBorder(result)

	b.SetCanvas(result, fg, bg)
}

func (c *Scrollable) Propagate() {
	hei := c.GetInnerHeight()
	c.UpdateVisibility(c.OffsetY, hei)
	pt := 0
	for _, child := range c.GetChildren() {
		child.UpdateVisibility(c.OffsetY-pt, hei)
		pt += child.GetHeight()
	}

	c.ComponentState.Propagate()
}
