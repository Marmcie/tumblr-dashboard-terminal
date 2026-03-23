package main

import (
	"fmt"
	"os"
	"tumblr-dt/modules/ui"
	component "tumblr-dt/modules/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

// import (
// 	"os"
// 	"tumblr-dt/modules"
//
// 	"github.com/rivo/tview"
// )
//
// func main() {
// 	client := modules.NewTumblrClient()
//
// 	app := tview.NewApplication()
//
// 	dashboard := modules.NewDashboard(&client, app)
//
// 	if len(os.Args) > 1 {
// 		// Black and white dashboard
// 		dashboard.BWMode = os.Args[1] == "BW"
// 	}
//
// 	// Load first sets of posts.
// 	dashboard.Update()
// 	dashboard.RenderPost()
//
// 	if err := app.SetRoot(dashboard.Root, true).Run(); err != nil {
// 		print("Error in tview loop\n")
// 		panic(err)
// 	}
// }

func main() {
	root := ui.NewRootModel()

	scroll := component.NewScrollable()
	scroll.ShowBorder = true
	scroll.BorderPadWidth = 1
	scroll.SetSize(100, 30)
	scroll.SetPos(0, 0)

	flex := component.NewFlex()
	flex.ShowBorder = true
	flex.BorderPadWidth = 1
	flex.SetPos(0, 0)
	flex.InheritWidth = true
	flex.InheritHeight = true

	for range 2 {
		var box = component.NewBox()
		box.Width = 10
		box.Height = 10
		box.ShowBorder = true
		box.BorderPadWidth = 1
		box.InheritWidth = true
		box.SetPos(0, 0)

		// var line = component.NewLine()
		// line.Text = "aaa"
		// box.AddChild(line)
		flex.AddItem(box, component.NewFlexDescriptor(0, 1))
	}

	scroll.AddChild(flex)
	scroll.Focus()
	root.App.SetRoot(scroll)

	p := tea.NewProgram(root)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
