package component_test

import (
	"testing"
	component "tumblr-dt/ui/component"
)

func BenchmarkAddChildrenToBox(b *testing.B) {
	box := component.NewBox("Box")
	for b.Loop() {
		line := component.NewText("t")
		box.AddChild(line)
	}
}

func BenchmarkBoxRenderChildren(b *testing.B) {
	elem := component.NewBox("elem")
	elem.SetSize(50, 50)
	for range 10 {
		line := component.NewText("t")
		line.SetWidthInherit(true)
		line.SetText("1")
		elem.AddChild(line)
	}
	for b.Loop() {
		elem.RenderToCanvas()
	}
}

func BenchmarkScrollableRenderChildren(b *testing.B) {
	elem := component.NewScrollable("elem")
	elem.SetSize(50, 50)
	for range 10 {
		line := component.NewText("t")
		line.SetH(1)
		line.SetWidthInherit(true)
		line.SetText("1")
		elem.AddChild(line)
	}

	for b.Loop() {
		elem.RenderToCanvas()
	}
}

func BenchmarkFlexRenderChildren(b *testing.B) {
	flex := component.NewFlex("Flex")
	flex.SetSize(50, 50)
	for range 10 {
		line := component.NewText("t")
		line.SetWidthInherit(true)
		line.SetText("1")
		flex.AddItem(line, component.NewFlexDescriptor(0, 1))
	}
	for b.Loop() {
		flex.RenderToCanvas()
	}
}

func TestRenderChildren(t *testing.T) {
	initItem := func(c component.Component) {
		c.SetSize(50, 50)
		for range 10 {
			line := component.NewText("t")
			line.SetWidthInherit(true)
			line.SetText("1")
			(c).AddChild(line)
		}
	}
	t.Run("Box", func(t *testing.T) {
		elem := component.NewBox("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("Flex", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("Scrollable", func(t *testing.T) {
		elem := component.NewScrollable("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
}

func TestRenderChildrenOverflow(t *testing.T) {
	initItem := func(c component.Component) {
		c.SetSize(50, 50)
		for range 10 {
			line := component.NewBox("t")
			line.SetSize(60, 60)
			(c).AddChild(line)
		}
	}
	t.Run("Box", func(t *testing.T) {
		elem := component.NewBox("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("Flex", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("FlexRow", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.Direction = 1
		elem.RenderToCanvas()
	})
	t.Run("Scrollable", func(t *testing.T) {
		elem := component.NewScrollable("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
}

func TestRenderChildrenEmpty(t *testing.T) {
	initItem := func(c component.Component) {
		c.SetSize(50, 50)
		for range 10 {
			line := component.NewBox("t")
			line.SetVisibility(false)
			line.SetSize(10, 10)
			(c).AddChild(line)
		}
	}
	t.Run("Box", func(t *testing.T) {
		elem := component.NewBox("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("Flex", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("FlexRow", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.Direction = 1
		elem.RenderToCanvas()
	})
	t.Run("Scrollable", func(t *testing.T) {
		elem := component.NewScrollable("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
}

func TestRenderChildrenTypes(t *testing.T) {
	initItem := func(c component.Component) {
		c.SetSize(50, 50)
		box := component.NewBox("t")
		box.SetSize(10, 10)
		(c).AddChild(box)

		scroll := component.NewScrollable("t")
		scroll.SetSize(10, 10)
		(c).AddChild(scroll)

		flex := component.NewFlex("t")
		flex.SetSize(10, 10)
		(c).AddChild(flex)

		line := component.NewLine("t")
		line.SetText("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		line.SetSize(10, 10)
		(c).AddChild(flex)

		text := component.NewText("t")
		text.SetText("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
		text.SetSize(10, 10)
		(c).AddChild(text)
		
	}
	t.Run("Box", func(t *testing.T) {
		elem := component.NewBox("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("Flex", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
	t.Run("FlexRow", func(t *testing.T) {
		elem := component.NewFlex("b")
		initItem(elem)
		elem.Direction = 1
		elem.RenderToCanvas()
	})
	t.Run("Scrollable", func(t *testing.T) {
		elem := component.NewScrollable("b")
		initItem(elem)
		elem.RenderToCanvas()
	})
}
