package component

// Basic component that displays its child.
type Box struct {
	BaseComponent
}

func NewBox(name string) *Box {
	b := &Box{}
	b.Initialize(name)
	b.SetComponentName("Box")
	return b
}

func (b *Box) PrepareFrame() {
	b.BaseComponent.PrepareFrame()
	var result, fg, bg = b.GetCanvas()
	result = b.addBorder(result)
	b.SetCanvas(result, fg, bg)
}
