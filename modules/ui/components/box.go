package component

import (
	helper "tumblr-dt/modules/ui/helper"
)

type Box struct {
	ComponentState
	ShowBorder     bool
	BorderPadWidth int
}

func NewBox() *Box {
	return &Box{}
}

// Returns x,y,width,height
func (c *Box) GetRect() (int, int, int, int) {
	return c.x, c.y, c.Width, c.Height
}

// Returns x,y,width,height
func (c *Box) GetRenderArea() (int, int, int, int) {
	return c.x + c.BorderPadWidth, c.y + c.BorderPadWidth, c.Width - c.BorderPadWidth, c.Height - c.BorderPadWidth
}

func (c *Box) addBorder(arr [][]string) [][]string {
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
func (b *Box) PrepareFrame() {
	var result = b.CreateCanvas()

	result = b.addBorder(result)

	for _, c := range b.Children {
		c.PrepareFrame()
		output := c.GetCanvas()
		childX, childY := c.GetPos()
		globalX := b.BorderPadWidth + childX

		for ind, line := range output {
			posY := ind + b.GetY() + childY + b.BorderPadWidth
			for index, char := range line {
				result[posY][globalX+index] = char
			}
		}
	}

	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}
