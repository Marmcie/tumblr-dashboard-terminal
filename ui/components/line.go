package component

import (
	"strings"
)

type Line struct {
	ComponentState
	Text string
}

func NewLine(name string) *Line {
	l := &Line{}
	l.Initialize(name)
	l.SetComponentName("Line")
	l.SetH(1)
	l.SetW(0)
	l.SetPos(0, 0)
	return l
}
func (l *Line) SetText(text string) *Line {
	l.Text = text
	return l
}

func (l *Line) GetRect() (int, int, int, int) {
	return l.x, max(0, l.y), len(l.Text), l.Height
}

// Returns Line per line contents,x,y
func (l *Line) PrepareFrame() {
	var result [][]string

	str := l.Text[:min(l.GetInnerWidth(), len(l.Text))]
	str = strings.ReplaceAll(str, "\n", "")
	result = append(result, strings.Split(str, ""))

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
