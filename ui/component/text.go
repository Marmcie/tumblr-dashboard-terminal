package component

import (
	"bytes"
	"strings"

	"github.com/rivo/uniseg"
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
	itr := uniseg.NewGraphemes(t.Text)
	for itr.Next() {
		r := itr.Str()
		x += uniseg.StringWidth(r)
		if r == string('\n') {
			x = 0
			res = append(res, strings.Split(strings.Trim(str.String(), " "), ""))
			str.Reset()
		} else {
			str.WriteString(r)
			for range uniseg.StringWidth(r) - 1 {
				str.WriteRune('\u200b')
			}
			if x >= width {
				res = append(res, strings.Split(strings.Trim(str.String(), " "), ""))
				x = 0
				str.Reset()
			}
		}
	}

	res = append(res, strings.Split(strings.Trim(str.String(), " "), ""))
	return res
}

// Returns Text per line contents,x,y
func (l *Text) RenderToCanvas() {
	if !l.Visibility {
		l.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		l.DispatchEvent("onRenderReady")
		return
	}
	arr := l.GetStringArray()
	top, bottom, left, _ := l.GetPaddings()

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

	result, fg, bg = l.addBorder(result, fg, bg)

	l.SetCanvas(result, fg, bg)
	l.DispatchEvent("onRenderReady")
}
