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

	l.Text = strings.ReplaceAll(text, "\n", "")

	if !l.InheritWidth {
		l.SetW(runewidth.StringWidth(text))
	}
	l.DispatchEvent("onChange")
	return l
}

// Returns Line per line contents,x,y
func (l *Line) PrepareFrame() {

	if !l.Visibility {
		l.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}
	w := l.GetInnerWidth()
	var result [][]string = [][]string{make([]string, w)}
	var fg [][]string = [][]string{make([]string, w)}
	var bg [][]string = [][]string{make([]string, w)}
	str := strings.Split(l.Text, "")
	for i := range w {
		if len(str) > i {
			result[0][i] = str[i]
		} else {
			result[0][i] = " "
		}
		fg[0][i] = l.Foreground
		bg[0][i] = l.Background
	}

	l.SetCanvas(result, fg, bg)
}
