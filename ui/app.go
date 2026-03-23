package ui

import (
	"bytes"
	"strings"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type App struct {
	root   *component.Component
	Height int
	Width  int
	Time   int
}

func NewApp() *App {
	return &App{}
}
func (m *App) UpdateSize(w int, h int) {
	m.Width = w
	m.Height = h
	(*m.root).SetSize(w, h)
}
func (m *App) GetSize() (int, int) {
	return m.Width, m.Height
}

func (m *App) SetRoot(child component.Component) {
	m.root = &child
	initializeDepth((*m.root), 0)
}

func initializeDepth(comp component.Component, depth int) {
	comp.SetDepth(depth)
	for _, c := range comp.GetChildren() {
		initializeDepth(c, depth+1)
	}
}

func (m *App) Render() string {
	(*m.root).Propagate()

	(*m.root).PrepareFrame()
	result, foreground, background := (*m.root).GetCanvas()
	var res bytes.Buffer
	style := lipgloss.NewStyle()
	for i := 0; i < min(len(result), m.Height); i++ {
		line := result[i]
		lineFG := foreground[i]
		lineBG := background[i]
		currentFG := lineBG[0]
		currentBG := lineBG[0]
		left := 0

		str := strings.Join(line, "")
		if len(lineFG) > 0 {
			ranges := []lipgloss.Range{}
			for i := range lineFG {
				fg := lineFG[i]
				bg := lineBG[i]
				if currentFG != fg || currentBG != bg {
					ranges = append(ranges, lipgloss.NewRange(left, i, style.Foreground(lipgloss.Color(currentFG)).Background(lipgloss.Color(currentBG))))
					currentFG = fg
					currentBG = bg
					left = i
				}
			}
			ranges = append(ranges, lipgloss.NewRange(left, len(line), style.Foreground(lipgloss.Color(currentFG)).Background(lipgloss.Color(currentBG))))
			str = lipgloss.StyleRanges(str, ranges...)
		}
		res.WriteString(str + "\n")
	}
	return res.String()
}

func (m *App) Update(msg tea.Msg) tea.Cmd {
	component.UpdateGlobalValues(msg)
	(*m.root).Update()
	component.Global.CallEvents()
	m.Time++
	cmd := component.Global.Command
	component.Global.Command = nil
	return cmd
}
