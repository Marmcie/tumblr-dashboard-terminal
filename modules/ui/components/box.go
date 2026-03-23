package component

import ()

type Box struct {
	ComponentState
}

func NewBox() *Box {
	b := &Box{}
	b.Initialize()
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
