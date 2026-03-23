package component

import (
	"bytes"
	"strings"
)

type Text struct {
	ComponentState
	Text string
}

func NewText(name string) *Text {
	l := &Text{}
	l.Initialize(name)
	l.SetComponentName("Text")
	l.SetPos(0, 0)
	return l
}

func (t *Text) SetText(v string) *Text {
	t.Text = v
	t.DispatchEvent("onChange")
	return t
}

func (t *Text) GetStringArray() [][]string {
	var res [][]string
	var str bytes.Buffer
	x := 0
	width := t.GetInnerWidth()
	for _, r := range t.Text {
		x++
		if r == '\n' {
			x = 0
			res = append(res, strings.Split(str.String(), ""))
			str.Reset()
		} else {
			str.WriteRune(r)
			if x >= width {
				res = append(res, strings.Split(str.String(), ""))
				x = 0
				str.Reset()
			}
		}
	}

	res = append(res, strings.Split(str.String(), ""))
	return res
}

// Returns Text per line contents,x,y
func (l *Text) PrepareFrame() {
	arr := l.GetStringArray()
	top, _, left, _ := l.GetBorderPaddings()

	innerHeight := l.GetInnerHeight()
	innerWidth := l.GetInnerWidth()
	var result = l.CreateCanvas()
	for i, line := range arr {
		if i+top >= innerHeight {
			break
		}
		for a, char := range line {
			if a+left >= innerWidth {
				break
			}
			result[i+top][a+left] = char
		}
	}

	result = l.addBorder(result)

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
