package component

import (
	"bytes"
	"strings"
)

// Component that displays multi line text
type Text struct {
	BaseComponent
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
	if !l.Visibility {
		l.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		l.DispatchEvent("onRenderReady")
		return
	}
	arr := l.GetStringArray()
	top, bottom, left, _ := l.GetBorderPaddings()

	innerWidth := l.GetInnerWidth()
	result, fg, bg := l.CreateCanvas()
	for i := range len(result) - (top + bottom) {
		for a := range result[i] {
			if a+left > innerWidth {
				break
			}
			if i >= len(arr) || a >= len(arr[i]) {
				result[i+top][a+left] = " "
			} else {
				result[i+top][a+left] = arr[i][a]
				fg[i+top][a+left] = l.Foreground
				bg[i+top][a+left] = l.Background
			}
		}
	}

	result = l.addBorder(result)

	l.SetCanvas(result, fg, bg)
	l.DispatchEvent("onRenderReady")
}
