package ui

import (
	"strings"
	component "tumblr-dt/modules/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct {
	children []*component.Component
	Height   int
	Width    int
	Time     int
}

func NewApp() *App {
	return &App{}
}

func (m *App) AddChild(child component.Component) {
	m.children = append(m.children, &child)
}

func (m *App) Render() string {
	var result = m.CreateCanvas()

	for _, c := range m.children {
		child := (*c)
		child.PrepareFrame()
		output := child.GetCanvas()
		x, y, _, _ := child.GetRect()

		for ind, line := range output {
			posY := ind + y
			if posY >= m.Height {
				break
			}

			for index, char := range line {
				result[posY][x+index] = char
			}
		}
	}

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
	component.UpdateGlobalValues(msg,m.Time)
	for _, c := range m.children {
		if (*c).GetFocusState() {
			(*c).Update()
		}
	}

	m.Time++

	return nil
}

func (a *App) CreateCanvas() [][]string {
	var arr [][]string

	for range a.Height {
		arr = append(arr, strings.Split(strings.Repeat(" ", a.Width), ""))
	}

	return arr
}
