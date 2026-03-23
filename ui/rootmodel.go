package ui

import (
	"time"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type RootModel struct {
	App *App
}

func TickCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg{
			Ms: t.UnixMilli(),
		}
	})
}

type TickMsg struct {
	Ms int64
}

func NewRootModel() RootModel {

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	app := NewApp()
	app.Width = s.Width
	app.Height = s.Height

	return RootModel{
		App: app,
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
	// TickCmd(),
	)
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}
	m.App.Width = s.Width
	m.App.Height = s.Height

	(*m.App.root).SetSize(s.Width, s.Height)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+d":
			modules.RemoveToken()
			return m, tea.Quit

		case "ctrl+l":
			component.Global.PrintLog()

		}

		return m, m.App.Update(msg)
	}
	return m, nil

}

func (m RootModel) View() tea.View {
	v := tea.NewView(m.App.Render())
	v.AltScreen = true
	return v
}
