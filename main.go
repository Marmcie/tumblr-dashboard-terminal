package main

import (
	"fmt"
	"os"
	"tumblr-dt/dashboard"
	"tumblr-dt/modules"

	tea "charm.land/bubbletea/v2"
)

func main() {
	config := modules.GetConfig()
	dashboard := dashboard.NewDashboard(config)

	p := tea.NewProgram(dashboard.GetRootModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bubbletea event loop error: %v", err)
		os.Exit(1)
	}
}
