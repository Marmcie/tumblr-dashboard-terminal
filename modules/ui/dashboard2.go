package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	tsize "github.com/kopoli/go-terminal-size"
)

type RootModel struct {
	App *App
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
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}
	m.App.Width = s.Width
	m.App.Height = s.Height

	(*m.App.root).SetSize(s.Width,s.Height)

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, m.App.Update(msg)
}

func (m RootModel) View() string {
	return m.App.Render()
}
