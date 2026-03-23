package dashboard

import (
	"tumblr-dt/modules"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	tsize "github.com/kopoli/go-terminal-size"
)

type Dashboard struct {
	core     ui.RootModel
	root     *component.Flex
	left     *component.Flex
	right    *component.Flex
	feed     Feed
	contents *Contents

	client modules.TumblrClient

	offset int
}

func NewDashboard() *Dashboard {
	d := &Dashboard{}

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	d.core = ui.NewRootModel()

	d.root = component.NewFlex("Root")
	d.root.SetDirection(1)
	d.root.SetSize(s.Width, s.Height)
	d.root.SetBorder(true).SetBorderPadding(1)

	d.left = component.NewFlex("Left")
	d.left.SetHeightInherit(true)
	d.left.Direction = 0

	d.right = component.NewFlex("Right")
	d.right.SetHeightInherit(true)
	d.right.Direction = 0

	d.feed = NewFeed()
	d.contents = NewContents()

	d.contents.contentElem.SetText("aaa")
	d.feed.contents = *d.contents

	d.left.AddItem(d.feed.listElem, component.NewFlexDescriptor(0, 1))
	d.right.AddItem(d.contents.contentElem, component.NewFlexDescriptor(0, 1))

	d.root.AddItem(d.left, component.NewFlexDescriptor(0, 1))
	d.root.AddItem(d.right, component.NewFlexDescriptor(0, 3))

	d.feed.listElem.Focus()
	d.core.App.SetRoot(d.root)

	d.client = modules.NewTumblrClient()
	d.offset = 0
	d.initEvents()

	return d
}

func (d *Dashboard) initEvents() {
	d.root.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "r":
				posts := d.client.GetDashboard(d.offset)
				d.feed.AddPosts(posts)
				d.offset++
			}
		}
	})
}

func (d *Dashboard) GetCore() ui.RootModel {
	return d.core
}
