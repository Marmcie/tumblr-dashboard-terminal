package main

import (
	"fmt"
	"os"
	"tumblr-dt/modules/ui"
	component "tumblr-dt/modules/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	tsize "github.com/kopoli/go-terminal-size"
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

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	main := component.NewFlex()
	main.ShowBorder = true
	main.BorderPadWidth = 1
	main.SetPos(0, 0)
	main.SetSize(s.Width, s.Height)
	main.Direction = 1

	slist := component.NewSelectlist()
	slist.ShowBorder = true
	slist.BorderPadWidth = 1
	slist.SetSize(100, 30)
	slist.SetPos(0, 0)

	b := component.NewBox()
	b.InheritHeight = true
	b.ShowBorder = true
	b.BorderPadWidth = 1
	bt := component.NewLine()
	bt.Text = "123"
	b.AddChild(bt)

	for i := range 2 {
		var box = component.NewBox()
		box.Width = 10
		box.Height = 2
		box.ShowBorder = true
		box.BorderPadWidth = 1
		box.InheritWidth = true
		box.SetBorders(false, true, false, false)
		box.SetBorderCorner(false)
		box.SetPos(0, 0)

		var line = component.NewLine()
		line.Text = "aaa"
		box.AddChild(line)
		slist.AddOption(box, func() {
			bt.Text = line.Text
		})
	}

	slist.AddEventListener("onChange", func(m tea.Msg, i int) {
		slist.OptionCallbacks[slist.Cursor]()
	})

	main.AddItem(slist, component.NewFlexDescriptor(0, 1))
	main.AddItem(b, component.NewFlexDescriptor(0, 2))

	slist.Focus()
	root.App.SetRoot(main)

	p := tea.NewProgram(root)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
