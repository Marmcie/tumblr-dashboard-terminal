package main

import (
	"fmt"
	"os"
	"strconv"
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
	main.SetSize(s.Width, s.Height)
	main.Direction = 1

	left := component.NewFlex()
	left.ShowBorder = true
	left.BorderPadWidth = 1
	left.ShowBorderCorner = true
	left.InheritHeight = true
	left.SetDirection(0)

	right := component.NewFlex()
	right.ShowBorder = true
	right.BorderPadWidth = 1
	right.ShowBorderCorner = true
	right.InheritHeight = true
	right.SetDirection(0)

	slist := component.NewSelectlist()
	slist.ShowBorder = true
	slist.BorderPadWidth = 1
	slist.SetPos(0, 0)
	slist.SetWidthInherit(true)

	content := component.NewText()
	content.ShowBorder = true
	content.BorderPadWidth = 1
	content.InheritWidth = true
	content.InheritHeight = true

	for i := range 3 {
		option := component.NewBox()
		option.
			SetBorder(true).
			SetBorderCorner(true).
			SetBorderPadding(1).
			SetH(4).
			SetWidthInherit(true)

		postTitle := component.NewLine().
			SetText(strconv.Itoa(i))
		postTitle.SetWidthInherit(true)

		postDate := component.NewLine().
			SetText("2026-02-20")

		postDate.SetWidthInherit(true)

		option.AddChild(postTitle)
		option.AddChild(postDate)
		slist.AddOption(option, func() {
			content.SetText(postTitle.Text + "\n" + postDate.Text)
		})
	}

	slist.AddEventListener("onChange", func(m tea.Msg, i int) {
		slist.OptionCallbacks[slist.Cursor]()
	})

	left.AddItem(slist, component.NewFlexDescriptor(0, 3))
	right.AddItem(content, component.NewFlexDescriptor(0, 3))

	main.AddItem(left, component.NewFlexDescriptor(0, 1))
	main.AddItem(right, component.NewFlexDescriptor(0, 3))

	content.SetName("content")
	right.SetName("right")
	main.SetName("main")
	left.SetName("left")
	slist.SetName("list")
	slist.SetTitle("0")
	main.SetTitle("0")

	main.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				if slist.GetFocusState() {
					content.Focus()
				} else {
					slist.Focus()
				}
			}
		}
	})

	slist.Focus()
	root.App.SetRoot(main)

	p := tea.NewProgram(root)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
