package component

import ()

type Box struct {
	ComponentState
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
