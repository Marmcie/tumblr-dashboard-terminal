package ui

import (
	"runtime"
	"time"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type RootModel struct {
	App       *App
	isWindows bool
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

	app := NewApp()

	return RootModel{
		App:       app,
		isWindows: runtime.GOOS == "windows",
	}
}

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
	// TickCmd(),
	)
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
	case tea.WindowSizeMsg:
		m.App.UpdateSize(msg.Width, msg.Height-1)
		return m, tea.ClearScreen
	}

	if m.isWindows {
		var s tsize.Size
		s, err := tsize.GetSize()
		if err != nil {
			println("Could not determine the size of the terminal window")
			panic(err)
		}
		w, h := (*m.App).GetSize()
		if w != s.Width || h != s.Height-1 {
			(*m.App).UpdateSize(s.Width, s.Height-1)
			return m, tea.ClearScreen
		}
	}
	return m, nil

}

func (m RootModel) View() tea.View {
	v := tea.NewView(m.App.Render())
	v.AltScreen = true
	return v
}
