package component

import (
	"github.com/rivo/uniseg"
	"strings"
)

// Component that displays a single line of text
type Line struct {
	BaseComponent
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

	l.Text = strings.ReplaceAll(text, "\n", "")

	if !l.InheritWidth {
		l.SetW(uniseg.StringWidth(text))
	}
	l.DispatchEvent("onChange")
	return l
}

// Returns Line per line contents,x,y
func (l *Line) RenderToCanvas() {

	w := l.GetInnerWidth()
	if !l.Visibility || w <= 0 {
		l.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}
	var result [][]string = [][]string{make([]string, w)}
	var fg [][]string = [][]string{make([]string, w)}
	var bg [][]string = [][]string{make([]string, w)}

	str := l.Text

	var c string
	state := -1
	i := 0
	var wid int
	for i < w {
		if len(str) > 0 {
			c, str, wid, state = uniseg.FirstGraphemeClusterInString(str, state)

			if i+wid <= w {
				result[0][i] = c
			} else {

				wid = 1
				result[0][i] = " "
			}
		} else {
			wid = 1
			result[0][i] = " "
		}
		fg[0][i] = l.Foreground
		bg[0][i] = l.Background
		i += wid
	}
	l.SetCanvas(result, fg, bg)
}
