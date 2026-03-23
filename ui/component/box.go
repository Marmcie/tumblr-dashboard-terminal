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

func (b *Box) RenderToCanvas() {
	b.BaseComponent.RenderToCanvas()
	var result, fg, bg = b.GetCanvas()
	result, fg, bg = b.addBorder(result, fg, bg)
	b.SetCanvas(result, fg, bg)
}
