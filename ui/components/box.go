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
	var result = b.GetCanvas()
	result = b.addBorder(result)
	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}

