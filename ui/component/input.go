package component

import (
	"tumblr-dt/ui/helper"
)

// Component that displays a single line of text
type Input struct {
	Line
	Value            string
	Placeholder      string
	EmptyForeground  string
	ActiveForeground string
	Suggestions      *helper.Trie
}

func NewInput(name string) *Input {
	l := &Input{}
	l.Initialize(name)
	l.SetComponentName("Input")
	l.SetH(1)
	l.SetW(0)
	l.SetPos(0, 0)
	l.EmptyForeground = "#aaaaaa"
	l.ActiveForeground = "#ffffff"

	return l
}

func (i *Input) ApplyTopSuggestion() {
	if len(i.Value) > 0 {
		suggestions := i.Suggestions.Search(i.Value)
		if len(suggestions) > 0 {
			i.Value = suggestions[0]
			i.UpdateText()
		}
	}
}

func (i *Input) SetSuggestions(trie *helper.Trie) {
	i.Suggestions = trie
}
func (i *Input) SetPlaceholder(s string) *Input {
	i.Placeholder = s
	i.UpdateText()
	return i
}

func (i *Input) UpdateText() {
	text := i.Value
	if len(text) == 0 {
		i.Line.SetForeground(i.EmptyForeground)
		i.Line.SetText(i.Placeholder)
		return
	}
	i.Line.SetForeground(i.ActiveForeground)
	i.Line.SetText(text)
}

func (i *Input) AppendChar(ch string) {
	i.Value += ch
	i.UpdateText()
}

func (i *Input) DeleteChar() {
	if len(i.Value) == 1 {
		i.Value = ""
	} else {
		i.Value = i.Value[:len(i.Value)-1]
	}
	i.UpdateText()
}

func (i *Input) ClearInput() {
	i.Value = ""
	i.UpdateText()
}

func (l *Input) RenderToCanvas() {

	if !l.Visibility {
		l.SetCanvas([][]string{{""}}, [][]string{{""}}, [][]string{{""}})
		return
	}
	l.Line.RenderToCanvas()
	canvas, fg, bg := l.GetCanvas()
	suggestion := ""
	if len(l.Value) > 0 {
		suggestions := l.Suggestions.Search(l.Value)
		if len(suggestions) > 0 {
			suggestion = suggestions[0]
		}
	}
	for i := len(l.Value); i < min(len(suggestion), len(canvas[0])); i++ {
		canvas[0][i] = string(suggestion[i])
		fg[0][i] = l.EmptyForeground
	}

	l.SetCanvas(canvas, fg, bg)
}
