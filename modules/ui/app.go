package ui

import (
	"strings"
	component "tumblr-dt/modules/ui/components"

	tea "github.com/charmbracelet/bubbletea"
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

func (m *App) SetRoot(child component.Component) {
	m.root = &child
}

func (m *App) Render() string {
	(*m.root).Propagate()
	
	(*m.root).PrepareFrame()
	result := (*m.root).GetCanvas()
	str := ""
	for i, line := range result {
		if i >= m.Height-1 {
			break
		}
		str += strings.Join(line, "")
		str += "\n"
	}
	return str
}

func (m *App) Update(msg tea.Msg) tea.Cmd {
	component.UpdateGlobalValues(msg, m.Time)
	(*m.root).Update()
	m.Time++
	return nil
}
