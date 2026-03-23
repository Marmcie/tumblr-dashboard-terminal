package component

type Scrollable struct {
	ComponentState
	OffsetY     int
	OffsetX     int
	InnerHeight int
	InnerWidth  int
}

func NewScrollable() *Scrollable {
	flex := &Scrollable{}
	flex.Initialize()
	flex.OffsetX = 0
	flex.OffsetY = 0
	flex.SetComponentName("Scrollable")

	// flex.AddEventListener("onUpdate", func(msg tea.Msg, time int) {
	//
	// 	switch msg := msg.(type) {
	//
	// 	// Is it a key press?
	// 	case tea.KeyMsg:
	//
	// 		// Cool, what was the actual key pressed?
	// 		switch msg.String() {
	//
	// 		// These keys should exit the program.
	// 		case "j":
	// 			flex.OffsetY = min(flex.InnerHeight-1, flex.OffsetY+1)
	// 		case "k":
	// 			flex.OffsetY = max(0, flex.OffsetY-1)
	//
	// 		case "l":
	// 			flex.OffsetX = min(flex.InnerWidth-1, flex.OffsetX+1)
	// 		case "h":
	// 			flex.OffsetX = max(0, flex.OffsetX-1)
	// 		}
	// 	}
	// })
	// flex.AddEventListener("onAddChild", func(msg tea.Msg, time int) {
	// 	w, h := flex.GetContentsSize()
	// 	flex.InnerHeight = h
	// 	flex.InnerWidth = w
	// })

	return flex
}

// Returns Line per line contents,x,y
func (b *Scrollable) PrepareFrame() {

	var result = b.CreateCanvas()
	b.ComponentState.PrepareFrame()
	var output = b.GetCanvas()

	for lineY, line := range output {
		if lineY >= b.OffsetY {
			if lineY-b.OffsetY >= b.GetHeight()-1 {
				break
			}
			for lineX, char := range line {
				if lineX-b.OffsetX >= b.GetWidth()-1 {
					break
				}

				if lineX >= b.OffsetX {
					result[lineY-b.OffsetY][lineX-b.OffsetX] = char
				}
			}
		}
	}

	result = b.addBorder(result)

	b.Canvas = result
	b.DispatchEvent("onRenderReady")

}
