package component

type Flex struct {
	ComponentState
	Direction   int
	InnerHeight int
	Descriptors []FlexDescriptor
}

type FlexDescriptor struct {
	Proportion int
	FixedSize  int
}

func NewFlexDescriptor(fixedsize int, proportion int) FlexDescriptor {
	return FlexDescriptor{
		FixedSize:  fixedsize,
		Proportion: proportion,
	}
}

func NewFlex() *Flex {
	flex := &Flex{}
	flex.Initialize()
	flex.Direction = 0
	flex.SetComponentName("Flex")
	return flex
}

func (c *Flex) AddChild(child Component) {
	c.ComponentState.AddChild(child)
	c.Descriptors = append(c.Descriptors, NewFlexDescriptor(1, 1))
}

func (c *Flex) AddItem(child Component, desc FlexDescriptor) {
	c.ComponentState.AddChild(child)
	c.Descriptors = append(c.Descriptors, desc)
}

func (f *Flex) GetProportionSum() int {
	res := 0
	for _, p := range f.Descriptors {
		res += p.Proportion
	}
	return res
}

func (b *Flex) UpdateChildSize() {
	flexH := b.GetInnerHeight()
	// Proportion should  fixed size

	proportionSum := b.GetProportionSum()

	for i, child := range b.GetChildren() {
		descriptor := b.Descriptors[i]

		if descriptor.Proportion > 0 {
			ratio := float64(descriptor.Proportion) / float64(proportionSum)
			childSize := int(float64(flexH) * ratio)
			child.SetH(childSize)
		} else {
			child.SetH(descriptor.FixedSize)
		}
	}
}

func (c *Flex) Propagate() {
	c.UpdateChildSize()
	c.ComponentState.Propagate()
}

// Returns Line per line contents,x,y
func (b *Flex) PrepareFrame() {
	b.ComponentState.PrepareFrame()
	var result = b.ComponentState.GetCanvas()
	result = b.addBorder(result)
	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}
