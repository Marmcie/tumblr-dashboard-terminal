package component

import (
	"strings"

	"github.com/mattn/go-runewidth"
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
	var result [][]string

	str := l.Text
	innerWidth := l.GetInnerWidth()

	str = strings.ReplaceAll(str, "\n", "")

	var res []string
	style := l.GetStyle()
	ct := 0
	for c := range strings.SplitSeq(str, "") {
		if ct >= innerWidth {
			break
		}

		if runewidth.StringWidth(c) > 0 {
			ct++
		}
		res = append(res, style.Render(c))
	}
	result = append(result, res)

	l.Canvas = result
	l.DispatchEvent("onRenderReady")
}
