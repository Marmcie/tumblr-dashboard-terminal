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

// import (
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"tumblr-dt/ui"
// 	component "tumblr-dt/ui/components"
//
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// 	tsize "github.com/kopoli/go-terminal-size"
// )
//
// func main() {
// 	root := ui.NewRootModel()
//
// 	var s tsize.Size
//
// 	s, err := tsize.GetSize()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	main := component.NewFlex("Main")
// 	main.SetSize(s.Width, s.Height-1)
// 	main.Direction = 1
// 	main.Gap = 0
// 	main.SetBorder(true).SetBorderPadding(1)
//
// 	left := component.NewFlex("Left")
// 	left.InheritHeight = true
// 	left.SetDirection(0)
// 	style := lipgloss.NewStyle().Background(lipgloss.Color("#000000"))
// 	left.SetStyle(style)
//
// 	right := component.NewFlex("Right")
// 	right.InheritHeight = true
// 	right.SetDirection(0)
// 	right.SetStyle(style)
//
// 	slist := component.NewSelectlist("Select")
// 	slist.ShowBorder = true
// 	slist.BorderPadWidth = 1
// 	slist.SetPos(0, 0)
// 	slist.SetWidthInherit(true)
//
// 	scroll := component.NewScrollable("Content scroll")
// 	scroll.SetHeightInherit(true).SetWidthInherit(true)
// 	scroll.SetBorder(true).SetBorderPadding(1)
// 	scroll.AddEventListener("onChange", func(m tea.Msg, i int) {
// 		scroll.OffsetY = 0
// 	})
//
// 	content := component.NewText("Content")
// 	// content.ShowBorder = true
// 	// content.BorderPadWidth = 1
// 	content.InheritWidth = true
// 	content.InheritHeight = true
//
// 	for i := range 3 {
// 		option := component.NewBox("Option :" + strconv.Itoa(i))
// 		option.
// 			SetBorder(true).
// 			SetBorderCorner(true).
// 			SetBorderPadding(1).
// 			SetH(4).
// 			SetWidthInherit(true)
//
// 		postTitle := component.NewLine("Post title").
// 			SetText(strconv.Itoa(i))
// 		postTitle.SetWidthInherit(true)
//
// 		postDate := component.NewLine("Post date").
// 			SetText("2026-02-20")
//
// 		postDate.SetWidthInherit(true)
//
// 		option.AddChild(postTitle)
// 		option.AddChild(postDate)
// 		slist.AddOption(option, func() {
// 			content.SetText(postTitle.Text + "\n" + postDate.Text)
// 		})
// 	}
//
// 	slist.AddEventListener("onChange", func(m tea.Msg, i int) {
// 		slist.OptionCallbacks[slist.Cursor]()
// 	})
//
// 	scroll.AddChild(content)
// 	left.AddItem(slist, component.NewFlexDescriptor(0, 3))
// 	right.AddItem(scroll, component.NewFlexDescriptor(0, 3))
//
// 	main.AddItem(left, component.NewFlexDescriptor(0, 1))
// 	main.AddItem(right, component.NewFlexDescriptor(0, 3))
//
// 	main.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
// 		switch msg := msg.(type) {
// 		case tea.KeyMsg:
// 			switch msg.String() {
// 			case "tab":
// 				if slist.GetFocusState() {
// 					scroll.Focus()
// 				} else {
// 					slist.Focus()
// 				}
// 			}
// 		}
// 	})
//
// 	slist.Focus()
// 	root.App.SetRoot(main)
//
// 	p := tea.NewProgram(root)
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Bubbletea event loop error: %v", err)
// 		os.Exit(1)
// 	}
// }
