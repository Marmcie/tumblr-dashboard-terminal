package main

import (
	"fmt"
	"os"
	"tumblr-dt/dashboard"
	"tumblr-dt/modules"

	tea "charm.land/bubbletea/v2"
)

func main() {

	modules.HandleArgs()
	config := modules.GetConfig()
	dashboard := dashboard.NewDashboard(config)

	fmt.Printf("\033]0;%s\a", "tumblr-dt")
	p := tea.NewProgram(dashboard.GetRootModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bubbletea event loop error: %v", err)
		os.Exit(1)
	}
}
