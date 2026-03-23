package dashboard

import (
	"fmt"
	"strings"
	"time"
	"tumblr-dt/modules"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type Dashboard struct {
	core      ui.RootModel
	root      *component.Flex
	left      *component.Flex
	right     *component.Flex
	switcher  *Switcher
	feed      *Feed
	info      *component.Text
	control   *component.Text
	contents  *Contents
	mode      string
	client    modules.TumblrClient
	offset    int
	timestamp int64
	tag       string
}

func NewDashboard() *Dashboard {
	d := &Dashboard{}

	var s tsize.Size

	s, err := tsize.GetSize()
	if err != nil {
		panic(err)
	}

	d.core = ui.NewRootModel()
	d.timestamp = time.Now().Local().UnixMilli() / 1000

	d.root = component.NewFlex("Root")
	d.root.SetDirection(1)
	d.root.SetSize(s.Width, s.Height)
	d.root.SetBorder(true).SetBorderPadding(1)
	d.root.SetBackground(ui.GetColorStr(ui.ColorBG))

	d.left = component.NewFlex("Left")
	d.left.SetHeightInherit(true)
	d.left.Direction = 0

	d.right = component.NewFlex("Right")
	d.right.SetHeightInherit(true)
	d.right.Direction = 0

	d.info = component.NewText("Info")
	d.info.SetWidthInherit(true).SetBorder(true).SetBorderPadding(1)

	d.control = component.NewText("Control")
	d.control.SetBorder(true).
		SetBorderPadding(1).
		SetSize(40, 8).
		SetTitle("Control").
		SetPos(0, 0).
		SetVisibility(false).
		SetAbsolute(true).
		SetCentered(true)

	d.switcher = NewSwitcher(d)

	d.feed = NewFeed(d)
	d.contents = NewContents(d)

	d.left.AddItem(d.feed.listElem, component.NewFlexDescriptor(0, 1))

	d.right.AddItem(d.contents.contentElem, component.NewFlexDescriptor(0, 1))
	d.right.AddItem(d.info, component.NewFlexDescriptor(9, 0))

	d.root.AddItem(d.left, component.NewFlexDescriptor(0, 1))
	d.root.AddItem(d.right, component.NewFlexDescriptor(0, 3))
	d.root.AddItem(d.control, component.NewFlexDescriptor(0, 3))
	d.root.AddItem(d.switcher.Window, component.NewFlexDescriptor(0, 3))

	d.feed.listElem.Focus()
	d.switcher.Window.Focus()
	d.core.App.SetRoot(d.root)

	d.client = modules.NewTumblrClient()
	d.offset = 0
	d.initEvents()
	d.SwitchMode("dashboard", "")
	d.UpdateControlText()

	return d
}
func (d *Dashboard) toggleControl() {
	d.control.SetVisibility(!d.control.Visibility)
}
func (d *Dashboard) toggleSwitcher() {
	state := !d.switcher.Window.GetVisibility()
	d.switcher.Window.SetVisibility(state)
	if state {
		d.switcher.DashOption.Focus()
	} else {
		d.feed.listElem.Focus()
	}
}

func (d *Dashboard) initEvents() {
	d.root.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "r":
				d.LoadPosts()

			case "o":
				post := d.feed.GetSelectedPost()
				modules.OpenInBrowser(post.Short_url)
				component.Global.SetCmd(tea.ClearScreen)

			case "]":
				d.toggleSwitcher()

			case "?":
				d.toggleControl()
			}

		}
	}, true)
}

func (d *Dashboard) SwitchMode(mode string, tag string) {
	d.switcher.Window.SetVisibility(false)
	d.feed.listElem.Focus()
	d.mode = mode
	switch d.mode {
	case "dashboard":
		d.feed.listElem.SetTitle("Dashboard")
	case "tag":
		d.tag = tag
		d.feed.listElem.SetTitle("Tagged posts : " + d.tag)
	}

	d.feed.ClearPosts()
	d.feed.listElem.ClearChildren()
	d.offset = 0
	d.LoadPosts()
}

func (d *Dashboard) LoadPosts() {
	var posts []npf.Post
	switch d.mode {
	case "dashboard":
		posts = d.client.GetDashboard(d.offset)
	case "tag":
		posts = d.client.GetTaggedPosts(int(d.timestamp), d.tag)
		if len(posts) == 0 {
			d.SwitchMode("dashboard", "")
			return
		}
		d.timestamp = posts[len(posts)-1].Timestamp
	}
	d.feed.AddPosts(posts)
	d.offset++

}

func (d *Dashboard) GetCore() ui.RootModel {
	return d.core
}

func (d *Dashboard) GetSelectedPost() npf.Post {
	return d.feed.GetSelectedPost()
}

func (d *Dashboard) FocusContents() {
	d.contents.Focus()
	d.UpdateControlText()
}

func (d *Dashboard) FocusFeed() {
	d.contents.contentElem.OffsetY = 0
	d.feed.Focus()
	d.UpdateControlText()
}

func (d *Dashboard) UpdateControlText() {
	str := ""
	if d.feed.listElem.GetFocusState() {
		str += "j/k      :  Scroll post on feed  \n"
		str += "l/Enter  :  Focus post window   \n"
		str += "r        :  Load more posts    \n"
		str += "]        :  Open feed switcher    \n"
		str += "o        :  Open post in browser    \n"
		str += "Ctrl+c   :  Exit the program  \n"
		str += "Ctrl+d   :  Log out of the account  \n"
	} else {
		str += "j/k      :  Scroll post contents  \n"
		str += "h        :  Focus feed  \n"
		str += "r        :  Load more posts     \n"
		str += "]        :  Open feed switcher    \n"
		str += "o        :  Open post in browser    \n"
		str += "Ctrl+c   :  Exit the program   \n"
		str += "Ctrl+d   :  Log out of the account  \n"
	}

	d.control.SetText(str)
}

func (d *Dashboard) DisplayPost(post npf.Post) {
	d.contents.DisplayPost(post)
	d.UpdateInfo(post)
}

func (d *Dashboard) UpdateInfo(post npf.Post) {
	config := modules.GetConfig()
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		panic(err)
	}
	timestamp := post.Timestamp
	t := time.Unix(timestamp, 0).In(loc)
	now := time.Now()
	diff := now.Sub(t)
	diffStr := ""
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) - (hours * 60)
	seconds := int(diff.Seconds()) - ((hours * 60 * 60) + (minutes * 60))

	diffStr = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	str := ""
	str += "Date      :  " + t.Format("2006-01-02 15:04:05 MST") + " (" + diffStr + " ago)" + "\n"
	str += "URL       :  " + post.Short_url + "\n"
	str += "Blog name :  " + post.Blog_name + "\n"
	str += "Tags      :  "

	if len(post.Tags) > 0 {
		str += "#"
		str += strings.Join(post.Tags, " #")
	}

	d.info.SetText(str)
}
