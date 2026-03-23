package component

import (
	"strings"
)

type Line struct {
	ComponentState
	Text string
}

func NewLine() *Line {
	l := &Line{}
	l.Initialize()
	l.SetComponentName("Line")
	l.SetH(1)
	l.SetW(0)
	l.SetPos(0,0)
	return l
}

func (l *Line) GetRect() (int, int, int, int) {
	return l.x, max(0, l.y), len(l.Text), l.Height
}

// Returns Line per line contents,x,y
func (l *Line) PrepareFrame() {
	var result [][]string
	l.Width = len(l.Text)

	result = append(result, strings.Split(l.Text, ""))

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
