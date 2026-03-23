package component

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

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

func (c *Scrollable) CreateCanvas() [][]string {
	var arr [][]string
	height := c.GetHeight()
	width := c.GetWidth()

	// height := c.GetHeight()

	for range height {
		arr = append(arr, strings.Split(strings.Repeat(" ", width), ""))
	}

	return arr
}

// Returns Line per line contents,x,y
func (b *Scrollable) PrepareFrame() {
	var result = b.CreateCanvas()
	b.ComponentState.PrepareFrame()

	var output = b.GetCanvas()
	boxHeight := b.GetInnerHeight()
	boxWidth := b.GetInnerWidth()
	style:=b.GetStyle()

	b.findBottom(output)
	for lineY, line := range output {
		if lineY < b.OffsetY {
			continue
		}
		if lineY-b.OffsetY > boxHeight {
			break
		}
		for lineX, char := range line {
			if lineX-b.OffsetX > boxWidth {
				break
			}

			if lineX >= b.OffsetX {
				result[lineY-b.OffsetY][lineX-b.OffsetX] = style.Render(char)
			}
		}
	}

	result = b.addBorder(result)

	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}
