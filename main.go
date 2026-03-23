package main

import (
	"fmt"
	"os"
	"tumblr-dt/dashboard"

	tea "charm.land/bubbletea/v2"
)

func main() {

	dashboard := dashboard.NewDashboard()

	p := tea.NewProgram(dashboard.GetCore())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bubbletea event loop error: %v", err)
		os.Exit(1)
	}
}

