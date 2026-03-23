package dashboard

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"tumblr-dt/modules"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

type Dashboard struct {
	config    modules.Config
	rootModel ui.RootModel
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
	option    string
}

func NewDashboard(config modules.Config) *Dashboard {
	d := &Dashboard{}
	d.config = config

	d.root = component.NewFlex("Root")
	d.root.SetDirection(1)
	d.root.SetBorder(true)
	d.root.SetBackground(ui.GetColorStr(ui.ColorBG))
	d.root.SetForeground(ui.GetColorStr(ui.ColorWhite))

	if !d.config.Testing {

		var s tsize.Size
		s, err := tsize.GetSize()
		if err != nil {
			panic(err)
		}
		d.root.SetSize(s.Width, s.Height)
	} else {

		d.root.SetSize(120, 60)
	}

	d.rootModel = ui.NewRootModel()
	d.timestamp = time.Now().Local().UnixMilli() / 1000

	d.left = component.NewFlex("Left")
	d.left.SetHeightInherit(true)
	d.left.Direction = 0

	d.right = component.NewFlex("Right")
	d.right.SetHeightInherit(true)
	d.right.Direction = 0

	d.info = component.NewText("Info")
	d.info.SetWidthInherit(true).SetBorder(true)

	d.control = component.NewText("Control")
	d.control.SetBorder(true).
		SetSize(40, 10).
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
	d.rootModel.App.SetRoot(d.root)

	d.client = modules.NewTumblrClient(d.config)
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
	d.root.AddEventListener("onUpdate", func(msg tea.Msg) {
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

			case "b":
				post := d.feed.GetSelectedPost()
				d.SwitchMode("blog", post.Blog.Name)

			case "?":
				d.toggleControl()
			}

		}
	}, true)
}

func (d *Dashboard) SwitchMode(mode string, option string) {
	d.switcher.Window.SetVisibility(false)
	d.feed.listElem.Focus()
	d.mode = mode
	switch d.mode {
	case "dashboard":
		d.feed.listElem.SetTitle("Dashboard")
	case "tag":
		d.option = option
		d.feed.listElem.SetTitle("Tagged posts : " + d.option)

	case "blog":
		d.option = option
		d.feed.listElem.SetTitle("Posts from : " + d.option)
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
		posts = d.client.GetTaggedPosts(int(d.timestamp), d.option)

	case "blog":
		posts = d.client.GetBlogPosts(int(d.timestamp), d.option)
	}
	if len(posts) == 0 {
		d.SwitchMode("dashboard", "")
		d.root.SetBorderLabel("BottomLeft", "Could not retrieve posts")
		return
	}

	d.timestamp = posts[len(posts)-1].Timestamp
	d.feed.AddPosts(posts)
	d.offset++

}

func (d *Dashboard) GetRootModel() ui.RootModel {
	return d.rootModel
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
		str += "b        :  Open blog feed    \n"
		str += "o        :  Open post in browser    \n"
		str += "Ctrl+c   :  Exit the program  \n"
		str += "Ctrl+d   :  Log out of the account  \n"
	} else {
		str += "j/k      :  Scroll post contents  \n"
		str += "h        :  Focus feed  \n"
		str += "r        :  Load more posts     \n"
		str += "]        :  Open feed switcher    \n"
		str += "b        :  Open blog feed    \n"
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
	config := d.config
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

	var str = bytes.Buffer{}
	str.WriteString("Date      :  " + t.Format("2006-01-02 15:04:05 MST") + " (" + diffStr + " ago)" + "\n")
	str.WriteString("URL       :  " + post.Short_url + "\n")
	str.WriteString("Blog name :  " + post.Blog_name + "\n")
	str.WriteString("Tags      :  ")
	if len(post.Tags) > 0 {
		str.WriteString("#")
		str.WriteString(strings.Join(post.Tags, " #"))
	}

	d.info.SetText(str.String())
}
