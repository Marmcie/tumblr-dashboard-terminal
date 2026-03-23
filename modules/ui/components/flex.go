package component

import (
	helper "tumblr-dt/modules/ui/helper"

	tea "github.com/charmbracelet/bubbletea"
)

type Flex struct {
	ComponentState
	ShowBorder     bool
	BorderPadWidth int
	FitHeight      bool
	FitWidth       bool
	Direction      int
	OffsetX        int
	OffsetY        int
}

func NewFlex() *Flex {
	flex := &Flex{}
	flex.FitHeight = true
	flex.FitWidth = true
	flex.Direction = 0
	flex.OffsetX = 0
	flex.OffsetY = 0
	return flex
}

// Returns x,y,width,height
func (c *Flex) GetRect() (int, int, int, int) {
	return c.x, c.y, c.Width, c.Height
}

func (c *Flex) addBorder(arr [][]string) [][]string {
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

// Returns Line per line contents,x,y
func (b *Flex) PrepareFrame() {

	// 0 col
	// 1 row

	var result = b.CreateCanvas()

	pos := 0
	renderPos := b.BorderPadWidth
	for _, c := range b.Children {
		c.PrepareFrame()
		output := c.GetCanvas()
		_, _, _, childH := c.GetRect()
		if pos+childH < b.OffsetY {
			pos += childH
			continue
		}

		for _, line := range output {
			pos++
			renderPos++
			for i, char := range line {
				index := i + b.BorderPadWidth
				if index >= b.Width-b.BorderPadWidth {
					break
				}
				result[renderPos][index] = char
			}
		}

	}

	result = b.addBorder(result)

	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}

func (c *Flex) OnUpdate(msg tea.Msg, time int) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "j":
			c.OffsetY += 1
		case "k":
			c.OffsetY = max(0, c.OffsetY-1)
		}
	}
}
