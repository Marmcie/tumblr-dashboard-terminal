package ui

import (
	"runtime"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type RootModel struct {
	App       *App
	isWindows bool
}

func NewRootModel() RootModel {

	app := NewApp()

	return RootModel{
		App:       app,
		isWindows: runtime.GOOS == "windows",
	}
}

func (m RootModel) Init() tea.Cmd {
	return nil
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if m.isWindows {
		var s tsize.Size
		s, err := tsize.GetSize()
		if err != nil {
			println("Could not determine the size of the terminal window")
			panic(err)
		}
		w, h := (*m.App).GetSize()
		if w != s.Width || h != s.Height {
			(*m.App).UpdateSize(s.Width, s.Height)
			return m, tea.ClearScreen
		}
	}
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		case "ctrl+d":
			modules.RemoveToken()
			return m, tea.Quit

		case "ctrl+l":
			component.Global.PrintLog()

		}

		return m, m.App.Update(msg)
	case tea.WindowSizeMsg:
		m.App.UpdateSize(msg.Width, msg.Height)
		return m, tea.ClearScreen
	}

	return m, nil

}

func (m RootModel) View() tea.View {
	v := tea.NewView(m.App.Render())
	v.AltScreen = true
	return v
}
