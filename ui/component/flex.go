package component

// Flex box component
type Flex struct {
	BaseComponent
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
	c.BaseComponent.AddChild(child)
	c.Descriptors = append(c.Descriptors, NewFlexDescriptor(1, 1))
}

func (c *Flex) AddItem(child Component, desc FlexDescriptor) {
	child.SetIsFlexItem(true)
	c.BaseComponent.AddChild(child)
	c.Descriptors = append(c.Descriptors, desc)
}

func (f *Flex) GetProportionSum() int {
	res := 0
	children := f.GetChildren()
	for i := 0; i < len(f.Descriptors); i++ {
		if children[i].IsAbsolute() || !children[i].GetVisibility() {
			continue
		}
		res += f.Descriptors[i].Proportion
	}
	return res
}

func (f *Flex) GetFixedSizeSum() int {
	res := 0
	children := f.GetChildren()
	for i, p := range f.Descriptors {
		if children[i].GetVisibility() {
			res += p.FixedSize
		}
	}
	return res
}

func (b *Flex) SetDirection(dir int) *Flex {
	b.Direction = dir
	return b
}

func (b *Flex) UpdateChildSize() {

	gapSize := b.Gap * (len(b.GetChildren()) - 1)
	fixedSize := b.GetFixedSizeSum()
	if b.Direction == 0 {

		flexH := b.GetInnerHeight() - (gapSize + fixedSize)

		// Proportion should  fixed size
		proportionSum := b.GetProportionSum()

		for i, child := range b.GetChildren() {
			if child.IsAbsolute() || !child.GetVisibility() {
				continue
			}
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
		flexW := b.GetInnerWidth() - (gapSize + fixedSize)
		// Proportion should  fixed size

		proportionSum := b.GetProportionSum()

		for i, child := range b.GetChildren() {

			if child.IsAbsolute() || !child.GetVisibility() {
				continue
			}
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

func (c *Flex) BeforeRender() {
	c.UpdateChildSize()
	c.BaseComponent.BeforeRender()
}

func (b *Flex) RenderToCanvas() {
	var result, fg, bg = b.CreateCanvas()

	top, _, left, _ := b.GetPaddings()
	cursor := top
	sideOffset := left

	for _, c := range b.GetChildren() {
		if !c.GetVisibility() {
			continue
		}
		c.RenderToCanvas()
		output, childFG, childBG := c.GetCanvas()
		if c.IsAbsolute() == true {
			childX, childY := c.GetPos()
			globalX := left + childX

			for ind, line := range output {
				posY := ind + b.GetY() + childY + top
				for index, char := range line {
					result[posY][globalX+index] = char
					if len(childFG[ind][index]) > 0 {
						fg[posY][globalX+index] = childFG[ind][index]
					}
					if len(childBG[ind][index]) > 0 {
						bg[posY][globalX+index] = childBG[ind][index]
					}
				}
			}
		} else {
			for i := 0; i < min(len(result), len(output)); i++ {
				line := output[i]
				if len(result) <= cursor {
					break
				}
				for a := 0; a < min(len(line), len(result[cursor])-sideOffset); a++ {
					char := line[a]
					result[cursor][a+sideOffset] = char
					if len(childFG[i][a]) > 0 {
						fg[cursor][a+sideOffset] = childFG[i][a]
					}
					if len(childBG[i][a]) > 0 {
						bg[cursor][a+sideOffset] = childBG[i][a]
					}
				}
				cursor++
			}
			if b.Direction == 1 {
				cursor = top
				sideOffset += c.GetWidth()
				sideOffset += b.Gap
			} else {
				cursor += b.Gap
			}
		}

	}

	result, fg, bg = b.addBorder(result, fg, bg)
	b.SetCanvas(result, fg, bg)
}
