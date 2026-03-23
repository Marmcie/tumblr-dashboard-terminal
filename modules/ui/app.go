package ui

import (
	"strings"
	component "tumblr-dt/modules/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	tsize "github.com/kopoli/go-terminal-size"
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
func (m *App) UpdateSize() {
	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}
	(*m.root).SetSize(s.Width, s.Height)
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
	m.UpdateSize()
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
	component.Global.CallEvents()
	m.Time++
	return nil
}
