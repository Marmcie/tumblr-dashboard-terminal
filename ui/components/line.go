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

	str := l.Text
	innerWidth := l.GetInnerWidth()

	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "'", "")

	var res []string
	style := l.GetStyle()
	ct := 0
	for _, c := range str {
		if ct >= innerWidth {
			break
		}
		if c == '\u200b' {
			innerWidth--
		}
		res = append(res, style.Render(string(c)))
		ct++
	}
	result = append(result, res)

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
