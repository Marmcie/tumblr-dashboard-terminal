package main

import (
	"fmt"
	"os"
	"tumblr-dt/dashboard"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	dashboard := dashboard.NewDashboard()

	p := tea.NewProgram(dashboard.GetCore())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bubbletea event loop error: %v", err)
		os.Exit(1)
	}
}
