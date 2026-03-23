package component

type Flex struct {
	ComponentState
	Direction   int
	InnerHeight int
	Descriptors []FlexDescriptor
	Gap         int
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

func NewFlex(name string) *Flex {
	flex := &Flex{}
	flex.Initialize(name)
	flex.Direction = 0
	flex.SetComponentName("Flex")
	flex.Gap = 0
	return flex
}

func (c *Flex) AddChild(child Component) {
	child.SetIsFlexItem(true)
	c.ComponentState.AddChild(child)
	c.Descriptors = append(c.Descriptors, NewFlexDescriptor(1, 1))
}

func (c *Flex) AddItem(child Component, desc FlexDescriptor) {
	child.SetIsFlexItem(true)
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

func (b *Flex) SetDirection(dir int) *Flex {
	b.Direction = dir
	return b
}

func (b *Flex) UpdateChildSize() {

	gapSize := b.Gap * (len(b.GetChildren()) - 1)
	if b.Direction == 0 {

		flexH := b.GetInnerHeight() - gapSize

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
	} else {
		flexW := b.GetInnerWidth() - gapSize
		// Proportion should  fixed size

		proportionSum := b.GetProportionSum()

		for i, child := range b.GetChildren() {
			descriptor := b.Descriptors[i]

			if descriptor.Proportion > 0 {
				ratio := float64(descriptor.Proportion) / float64(proportionSum)
				childSize := int(float64(flexW) * ratio)
				child.SetW(childSize)

			} else {
				child.SetW(descriptor.FixedSize)
			}
		}
	}
}

func (c *Flex) Propagate() {
	c.UpdateChildSize()
	c.ComponentState.Propagate()
}

func (b *Flex) PrepareFrame() {
	var result = b.CreateCanvas()

	top, _, left, _ := b.GetBorderPaddings()
	cursor := top
	sideOffset := left

	for _, c := range b.GetChildren() {
		c.PrepareFrame()
		output := c.GetCanvas()
		style := c.GetStyle()
		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := left + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + top
				for index, char := range line {
					result[posY][globalX+index] = char
				}
			}
		} else {
			for _, line := range output {
				if cursor >= len(result) {
					break
				}
				for i, char := range line {
					index := i
					if index >= len(result[cursor]) {
						break
					}
					result[cursor][index+sideOffset] = style.Render(char)
				}
				cursor++
			}
		}
		if b.Direction == 1 {
			cursor = top
			sideOffset += c.GetWidth()
			sideOffset += b.Gap
		} else {
			cursor += b.Gap
		}
	}

	result = b.addBorder(result)
	b.Canvas = result
	b.DispatchEvent("onRenderReady")
}
