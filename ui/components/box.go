package component

// Basic component that displays its child.
type Box struct {
	ComponentState
}

func NewBox(name string) *Box {
	b := &Box{}
	b.Initialize(name)
	b.SetComponentName("Box")
	return b
}

func (b *Box) PrepareFrame() {
	b.ComponentState.PrepareFrame()
	var result, fg, bg = b.GetCanvas()
	result = b.addBorder(result)
	b.SetCanvas(result, fg, bg)
}
