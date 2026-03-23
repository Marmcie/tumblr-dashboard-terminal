package ui

import (
	"bytes"
	"strings"
	component "tumblr-dt/ui/components"

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
	result, foreground, _ := (*m.root).GetCanvas()
	var res bytes.Buffer
	for i := 0; i < min(len(result), m.Height); i++ {
		line := result[i]
		lineFG := foreground[i]
		currentFG := lineFG[0]
		left := 0
		ranges := []lipgloss.Range{}
		for ind, fg := range lineFG {
			if currentFG != fg {
				ranges = append(ranges, lipgloss.NewRange(left, ind-1, lipgloss.NewStyle().Foreground(lipgloss.Color(currentFG))))
				currentFG = fg
				left = ind
			}
		}
		ranges = append(ranges, lipgloss.NewRange(left, len(lineFG)-1, lipgloss.NewStyle().Foreground(lipgloss.Color(currentFG))))

		str := lipgloss.StyleRanges(strings.Join(line, ""), ranges...)
		res.WriteString(str + "\n")
	}
	return res.String()
}

func (m *App) Update(msg tea.Msg) tea.Cmd {
	component.UpdateGlobalValues(msg, m.Time)
	(*m.root).Update()
	component.Global.CallEvents()
	m.Time++
	cmd := component.Global.Command
	component.Global.Command = nil
	return cmd
}
