package ui

import (
	"runtime"
	"time"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type RootModel struct {
	App       *App
	isWindows bool
	config    modules.Config
}

func NewRootModel(config modules.Config) RootModel {

	app := NewApp()

	return RootModel{
		App:       app,
		isWindows: runtime.GOOS == "windows",
		config:    config,
	}
}
func TickCmd(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return TickMsg(t.UnixMilli())
	})
}

type TickMsg int64

func (m RootModel) Init() tea.Cmd {
	return tea.Batch(
		TickCmd(component.Global.TickInterval),
	)
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	commands := []tea.Cmd{}
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
			commands = append(commands, tea.ClearScreen)
		}
	}
	switch msg := msg.(type) {
	case TickMsg:
		m.App.Update(msg)
		commands = append(commands, TickCmd(component.Global.TickInterval))

	case tea.KeyPressMsg:
		switch msg.String() {
		case m.config.Keymaps.Quit, "ctrl+c":
			commands = append(commands, tea.Quit)

		case m.config.Keymaps.LogOut:
			modules.RemoveToken()
			commands = append(commands, tea.Quit)

		case m.config.Keymaps.Log:
			component.Global.PrintLog()

		}

		commands = append(commands, m.App.Update(msg))
	case tea.WindowSizeMsg:
		m.App.UpdateSize(msg.Width, msg.Height)
		commands = append(commands, tea.ClearScreen)
	}

	return m, tea.Batch(
		commands...,
	)

}

func (m RootModel) View() tea.View {
	v := tea.NewView(m.App.Render())
	v.AltScreen = true
	return v
}
