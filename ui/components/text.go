package component

import (
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
	str := ""
	x := 0
	width := t.GetInnerWidth()
	for c := range strings.SplitSeq(t.Text, "") {
		x++
		if c == "\n" {
			x = 0
			res = append(res, strings.Split(str, ""))
			str = ""
		} else {
			str += c
			if x >= width {
				res = append(res, strings.Split(str, ""))
				x = 0
				str = ""
			}
		}
	}

	res = append(res, strings.Split(str, ""))
	return res
}

// Returns Text per line contents,x,y
func (l *Text) PrepareFrame() {

	arr := l.GetStringArray()
	top, bottom, left, _ := l.GetBorderPaddings()
	l.SetH(len(arr) + top + bottom)

	var result = l.CreateCanvas()
	for i, line := range arr {
		for a, char := range line {
			result[i+top][a+left] = char
		}
	}

	result = l.addBorder(result)

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
