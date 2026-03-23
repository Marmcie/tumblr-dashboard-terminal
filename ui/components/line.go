package component

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

// Component that displays a single line of text
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

	if !l.InheritWidth {
		l.SetW(runewidth.StringWidth(l.Text))
	}
	l.DispatchEvent("onChange")
	return l
}

// Returns Line per line contents,x,y
func (l *Line) PrepareFrame() {
	if !l.Visibility {
		l.Canvas = [][]string{{""}}
		l.DispatchEvent("onRenderReady")
		return
	}

	str := l.Text
	innerWidth := l.GetInnerWidth()

	var result [][]string
	str = strings.ReplaceAll(str, "\n", "")

	var res []string
	parts := strings.Split(str, "")
	res = strings.Split(strings.Repeat(" ", max(innerWidth, len(parts))), "")

	style := l.GetStyle()
	ct := 0
	for i, c := range parts {
		if ct >= innerWidth {
			break
		}

		if runewidth.StringWidth(c) > 0 {
			ct++
		}
		res[i] = style.Render(c)
	}
	result = append(result, res)

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
