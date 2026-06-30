package dashboard

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"strings"
	"time"
	"tumblr-dt/modules"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"
	"tumblr-dt/ui/helper"

	tea "charm.land/bubbletea/v2"
	mapset "github.com/deckarep/golang-set/v2"
	tsize "github.com/kopoli/go-terminal-size"
)

var girlNo []string = []string{
	"Girl no",
	"It's not worth it",
	"Bad girl",
	"Awawawawa",
	"Bonk",
}

type Dashboard struct {
	config           modules.Config
	rootModel        ui.RootModel
	root             *component.Flex
	left             *component.Flex
	right            *component.Flex
	switcher         *Switcher
	feed             *Feed
	info             *component.Text
	control          *component.Text
	contents         *Contents
	mode             string
	client           modules.TumblrClient
	offset           int
	next             string
	timestamp        int64
	option           string
	TagTrie          *helper.Trie
	BlogTrie         *helper.Trie
	FilteredTags     mapset.Set[string]
	FilteredContents mapset.Set[string]
	IsLoading        bool
	LinkWindow       *LinkWindow
}

func NewDashboard(config modules.Config) *Dashboard {
	d := &Dashboard{}
	d.config = config
	d.IsLoading = false

	d.FilteredTags = mapset.NewSet[string]()
	d.FilteredContents = mapset.NewSet[string]()
	d.TagTrie = helper.NewTrie()
	d.BlogTrie = helper.NewTrie()

	d.initComponents(d.config)
	d.initEvents()
	d.UpdateControlText()

	if !d.config.Initialized {
		d.SwitchMode("tutorial", "")
	} else {
		d.initClient(d.config)

		d.SwitchMode("dashboard", "")
	}

	return d
}

func (d *Dashboard) initClient(config modules.Config) {
	d.client = modules.NewTumblrClient(config)
	d.offset = 0

	filteredTagsChan := make(chan []string)
	filteredContentsChan := make(chan []string)

	go d.client.GetFilteredTags(filteredTagsChan)
	go d.client.GetFilteredContents(filteredContentsChan)

	d.FilteredTags = mapset.NewSet[string](<-filteredTagsChan...)
	d.FilteredContents = mapset.NewSet[string](<-filteredContentsChan...)
}

func (d *Dashboard) initComponents(config modules.Config) {
	d.root = component.NewFlex("Root")
	d.root.SetDirection(1)
	d.root.SetBorder(true)
	d.root.SetBackground(ui.GetColorStr(ui.ColorBG))
	d.root.SetForeground(ui.GetColorStr(ui.ColorWhite))

	if !config.Testing {
		var s tsize.Size
		s, err := tsize.GetSize()
		if err != nil {
			panic(err)
		}
		d.root.SetSize(s.Width, s.Height)
	} else {

		d.root.SetSize(120, 60)
	}

	d.rootModel = ui.NewRootModel(d.config)
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
		SetSize(40, 12).
		SetTitle("Control").
		SetPos(0, 0).
		SetVisibility(false).
		SetAbsolute(true).
		SetCentered(true).
		SetForeground(ui.GetColorStr(ui.ColorWhite))
	d.control.SetBorderForeground(ui.GetColorStr(ui.ColorWhite))

	d.switcher = NewSwitcher(d)
	d.switcher.TagInput.SetSuggestions(d.TagTrie)
	d.switcher.BlogInput.SetSuggestions(d.BlogTrie)
	d.switcher.Window.SetVisibility(false)

	d.feed = NewFeed(d)
	d.contents = NewContents(d)
	d.contents.contentElem.SetW(d.root.GetWidth() / 2)

	d.LinkWindow = NewLinkWindow(d)

	d.left.AddItem(d.feed.listElem, 0, 1)

	d.right.AddItem(d.contents.contentElem, 0, 1)
	d.right.AddItem(d.info, 9, 0)

	d.root.AddItem(d.left, 0, 1)
	d.root.AddItem(d.right, 0, 3)
	d.root.AddItem(d.control, 0, 3)
	d.root.AddItem(d.switcher.Window, 0, 3)
	d.root.AddItem(d.LinkWindow.Window, 0, 3)

	d.feed.listElem.Focus()
	d.rootModel.App.SetRoot(d.root)

}

func (d *Dashboard) toggleControl() {
	d.control.SetVisibility(!d.control.Visibility)
}
func (d *Dashboard) toggleSwitcher() {
	state := !d.switcher.Window.GetVisibility()
	d.switcher.Window.SetVisibility(state)
	if state {
		d.ShowFeed()
		d.switcher.DashOption.Focus()
	} else {
		d.feed.listElem.Focus()
	}
}
func (d *Dashboard) toggleLinkWindow() {
	state := !d.LinkWindow.Window.GetVisibility()
	d.switcher.Window.SetVisibility(state)
	if state {
		d.LinkWindow.Focus()
	} else {
		d.feed.listElem.Focus()
		d.LinkWindow.Blur()
	}
}

func (d *Dashboard) initEvents() {

	d.root.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case d.config.Keymaps.LoadMore:
				done := make(chan bool)
				go d.LoadPosts(done)

			case d.config.Keymaps.IncreaseSize:
				// When feed is visible
				if d.left.GetVisibility() {
					d.left.SetFlexProportion(d.left.GetFlexProportion() + 0.1)
					d.feed.listElem.RunSelectedOption()
				} else {
					// When feed is hidden
					d.contents.contentElem.SetW(d.contents.contentElem.GetWidth() + 1)
				}

			case d.config.Keymaps.DecreaseSize:
				// When feed is visible
				if d.left.GetVisibility() {
					proportion := d.left.GetFlexProportion()
					d.left.SetFlexProportion(max(0.1, proportion-0.1))
					d.feed.listElem.RunSelectedOption()
				} else {
					// When feed is hidden
					d.contents.contentElem.SetW(d.contents.contentElem.GetWidth() - 1)
				}

			case d.config.Keymaps.Quit:
				component.Global.SetCmd(tea.Quit)

			case d.config.Keymaps.OpenLink:
				post := d.feed.GetSelectedPost()
				modules.OpenInBrowser(post.Short_url)
				component.Global.SetCmd(tea.ClearScreen)

			case d.config.Keymaps.Switcher.Open:
				d.toggleSwitcher()

			case d.config.Keymaps.LoadBlog:
				post := d.feed.GetSelectedPost()
				d.SwitchMode("blog", post.Blog.Name)

			case d.config.Keymaps.ControlHelp:
				d.toggleControl()
			case d.config.Keymaps.Links.Open:
				d.LinkWindow.Focus()
			case d.config.Keymaps.ToggleFeed:
				d.feed.showFilteredPost = true
				d.ToggleFeed()
			}

		}
	}, true)
}

func (d *Dashboard) SwitchMode(mode string, option string) {
	d.switcher.Window.SetVisibility(false)
	d.feed.listElem.Focus()

	if mode != "dashboard" && slices.Contains(d.config.Blacklist, option) {
		d.root.SetBorderLabel("BottomLeft", girlNo[rand.Intn(len(girlNo))])
		d.root.SetBorderLabelColor("BottomLeft", ui.GetColorStr(ui.ColorBlacklisted))
		return
	}

	d.mode = mode
	d.timestamp = time.Now().Local().UnixMilli() / 1000
	switch d.mode {
	case "dashboard":
		d.feed.listElem.SetTitle("Dashboard")
	case "tag":
		d.option = option
		d.feed.listElem.SetTitle(fmt.Sprintf("Tagged posts : %s", d.option))

	case "search":
		d.option = option
		d.feed.listElem.SetTitle(fmt.Sprintf("Searched posts : %s", d.option))
		d.next = ""

	case "blog":
		d.option = option
		d.feed.listElem.SetTitle(fmt.Sprintf("Posts from : %s", d.option))
	case "tutorial":
		d.feed.listElem.SetTitle("Tutorial")
	}

	d.root.SetBorderLabelColor("BottomLeft", "")

	d.feed.ClearPosts()
	d.feed.listElem.ClearChildren()
	d.offset = 0
	d.contents.contentElem.OffsetY = 0
	done := make(chan bool)
	go d.LoadPosts(done)
	<-done
}

func (d *Dashboard) filterPosts(posts []npf.Post) []*npf.Post {
	result := []*npf.Post{}
	filteredContents := d.FilteredContents.ToSlice()
	for _, postObject := range posts {
		post := &postObject
		if post.Type != "blocks" {
			continue
		}
		result = append(result, post)
		post.FilteredContents = mapset.NewSet[string]()
		post.FilteredTags = mapset.NewSet[string]()

		if len(post.Tags) > 0 {
			for _, tag := range post.Tags {
				if d.FilteredTags.Contains(tag) {
					post.IsFiltered = true
					post.FilteredTags.Add(tag)
				}
			}
		}
		reblogs := post.Render()
		for i := range reblogs {
			reblog := reblogs[i]
			for a := range filteredContents {
				filteredWord := filteredContents[a]
				if strings.Contains(reblog.Blog.Name, filteredWord) {
					post.IsFiltered = true
					post.FilteredContents.Add(filteredWord)
				}
				for b := 0; b < len(reblog.Contents); b++ {
					content := reblog.Contents[b]
					if strings.Contains(content.Str, filteredWord) {
						post.IsFiltered = true
						post.FilteredContents.Add(filteredWord)
					}
				}
			}
		}
	}
	return result
}

func (d *Dashboard) LoadPosts(ch chan bool) {
	defer func() { ch <- true }()

	if d.IsLoading {
		return
	}
	d.feed.listElem.SetBorderLabel("Bottom", "Loading...")
	d.IsLoading = true
	var posts []npf.Post
	switch d.mode {
	case "dashboard":
		posts = d.client.GetDashboard(d.offset)
	case "tag":
		posts = d.client.GetTaggedPosts(int(d.timestamp), d.option)
	case "blog":
		posts = d.client.GetBlogPosts(int(d.timestamp), d.option)

	case "search":
		posts, d.next = d.client.GetSearchedPosts(int(d.timestamp), d.option, d.next)

	case "tutorial":
		posts = d.client.GetTutorial()
	}
	if len(posts) == 0 {
		defer func() {
			d.feed.listElem.SetBorderLabel("Bottom", "")
			d.root.SetBorderLabel("BottomLeft", "Could not retrieve posts")
		}()

		d.IsLoading = false
		return
	}

	sort.Sort(npf.SortPostByTimestamp(posts))
	d.root.SetBorderLabel("BottomLeft", "")

	if d.mode != "tutorial" {
		d.timestamp = posts[len(posts)-1].Timestamp
	}

	result := d.filterPosts(posts)
	d.feed.AddPosts(result)
	d.offset++

	if d.mode != "tutorial" {
		for _, p := range posts {
			for _, tag := range p.Tags {
				d.TagTrie.Insert(tag)
			}
			d.BlogTrie.Insert(p.Blog.GetName())
			for _, t := range p.Trail {
				d.BlogTrie.Insert(t.Blog.Name)
			}
		}
	}
	d.IsLoading = false
	component.Global.TickInterval = time.Second / 15

	d.feed.listElem.SetBorderLabel("Bottom", "")
}

func (d *Dashboard) GetRootModel() ui.RootModel {
	return d.rootModel
}

func (d *Dashboard) GetSelectedPost() *npf.Post {
	return d.feed.GetSelectedPost()
}

func (d *Dashboard) FocusContents() {
	d.contents.Focus()
	d.UpdateControlText()
}

func (d *Dashboard) FocusFeed() {
	d.contents.contentElem.OffsetY = 0
	d.feed.Focus()
	d.ShowFeed()
	d.UpdateControlText()
}

func (d *Dashboard) ToggleFeed() {
	if d.left.GetVisibility() {
		d.HideFeed()
	} else {
		d.ShowFeed()
	}
}
func (d *Dashboard) ShowFeed() {
	d.left.SetVisibility(true)
	d.contents.contentElem.SetAbsolute(false)
	d.contents.contentElem.SetWidthInherit(true)
	d.info.SetVisibility(true)
	d.feed.Focus()
	d.contents.contentElem.SetBorderLabel("BottomLeft", "")
}
func (d *Dashboard) HideFeed() {
	d.left.SetVisibility(false)
	d.contents.contentElem.SetAbsolute(true)
	d.contents.contentElem.SetWidthInherit(false)
	if d.contents.contentElem.GetWidth() > d.root.GetWidth() {
		d.contents.contentElem.SetW(d.root.GetWidth() / 2)
	}
	d.info.SetVisibility(false)
	d.contents.Focus()
	d.contents.contentElem.SetBorderLabel("BottomLeft", fmt.Sprintf("[%s] to restore feed", d.config.Keymaps.ToggleFeed))
}

func (d *Dashboard) UpdateControlText() {
	str := ""
	if d.feed.listElem.GetFocusState() {

		str += fmt.Sprintf("Scroll post on feed     : %s/%s\n", d.config.Keymaps.Navigation.Up, d.config.Keymaps.Navigation.Down)
		str += fmt.Sprintf("Scroll to top or bottom : %s%s/%s\n", d.config.Keymaps.Navigation.JumpTop, d.config.Keymaps.Navigation.JumpTop, d.config.Keymaps.Navigation.JumpBottom)
		str += fmt.Sprintf("Focus post window       : %s\n", d.config.Keymaps.Navigation.Right)
		str += fmt.Sprintf("Load more posts         : %s\n", d.config.Keymaps.LoadMore)
		str += fmt.Sprintf("Open feed switcher      : %s\n", d.config.Keymaps.Switcher.Open)
		str += fmt.Sprintf("Open post links         : %s\n", d.config.Keymaps.Links.Open)
		str += fmt.Sprintf("Open blog feed          : %s\n", d.config.Keymaps.LoadBlog)
		str += fmt.Sprintf("Open post in browser    : %s\n", d.config.Keymaps.OpenLink)
		str += fmt.Sprintf("Adjust feed width       : %s/%s\n", d.config.Keymaps.IncreaseSize, d.config.Keymaps.DecreaseSize)
		str += fmt.Sprintf("Toggle feed visibility  : %s\n", d.config.Keymaps.ToggleFeed)
		str += fmt.Sprintf("Exit the program        : Ctrl+c/%s\n", d.config.Keymaps.Quit)
		str += fmt.Sprintf("Log out of the account  : %s\n", d.config.Keymaps.LogOut)
		d.control.SetH(14)
	} else {
		str += fmt.Sprintf("Scroll post contents    : %s/%s %s/%s\n", d.config.Keymaps.Navigation.Up, d.config.Keymaps.Navigation.Down, strings.ToUpper(d.config.Keymaps.Navigation.Up), strings.ToUpper(d.config.Keymaps.Navigation.Down))
		str += fmt.Sprintf("Scroll to next reblog   : %s/%s\n", d.config.Keymaps.Navigation.JumpNext, d.config.Keymaps.Navigation.JumpPrev)
		str += fmt.Sprintf("Scroll to top or bottom : %s%s/%s\n", d.config.Keymaps.Navigation.JumpTop, d.config.Keymaps.Navigation.JumpTop, d.config.Keymaps.Navigation.JumpBottom)
		str += fmt.Sprintf("Focus feed              : Ctrl+c/%s\n", d.config.Keymaps.Navigation.Left)
		str += fmt.Sprintf("Load more posts         : %s\n", d.config.Keymaps.LoadMore)
		str += fmt.Sprintf("Open feed switcher      : %s\n", d.config.Keymaps.Switcher.Open)
		str += fmt.Sprintf("Open post links         : %s\n", d.config.Keymaps.Links.Open)
		str += fmt.Sprintf("Open blog feed          : %s\n", d.config.Keymaps.LoadBlog)
		str += fmt.Sprintf("Open post in browser    : %s\n", d.config.Keymaps.OpenLink)
		str += fmt.Sprintf("Adjust feed width       : %s/%s\n", d.config.Keymaps.IncreaseSize, d.config.Keymaps.DecreaseSize)
		str += fmt.Sprintf("Toggle feed visibility  : %s\n", d.config.Keymaps.ToggleFeed)
		str += fmt.Sprintf("Exit the program        : Ctrl+c/%s\n", d.config.Keymaps.Quit)
		str += fmt.Sprintf("Log out of the account  : %s\n", d.config.Keymaps.LogOut)
		d.control.SetH(15)
	}

	d.control.SetText(str)
}

func (d *Dashboard) DisplayPost(post *npf.Post, showFiltered bool) {
	d.contents.DisplayPost(post, showFiltered)
	d.LinkWindow.SetLinks(post.GetLinks())

	d.contents.contentElem.OffsetY = 0

	d.UpdateInfo(post)
}

func (d *Dashboard) UpdateInfo(post *npf.Post) {
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
	days := hours / 24

	if days > 0 {
		diffStr = fmt.Sprintf("%d days", days)
	} else {
		if hours > 0 {
			diffStr = fmt.Sprintf("%d hours", hours)
		} else {
			if minutes > 0 {
				diffStr = fmt.Sprintf("%d minutes", minutes)
			} else {
				diffStr = fmt.Sprintf("%d seconds", seconds)
			}
		}
	}

	//TODO: Implement this better
	//Perhaps use line objects?
	// var str = bytes.Buffer{}
	var b strings.Builder
	if !post.IsFiltered {
		fmt.Fprintf(&b, "Date      :  %s (%s ago)\n", t.Format("2006-01-02 15:04:05 MST"), diffStr)
		fmt.Fprintf(&b, "URL       :  %s\n", post.Short_url)
		fmt.Fprintf(&b, "Blog name :  %s\n", post.Blog_name)
		fmt.Fprintf(&b, "Tags      :  ")
		if len(post.Tags) > 0 {
			fmt.Fprintf(&b, "#%s", strings.Join(post.Tags, " #"))
		}
	} else {

		filteredContents := post.FilteredContents.ToSlice()
		filteredTags := post.FilteredTags.ToSlice()

		fmt.Fprintf(&b, "Date      :  %s (%s ago)\n", t.Format("2006-01-02 15:04:05 MST"), diffStr)
		fmt.Fprintf(&b, "URL       :  %s\n", post.Short_url)
		fmt.Fprintf(&b, "Blog name :  %s\n", post.Blog_name)
		if len(filteredContents) > 0 {
			fmt.Fprintf(&b, "Filtered contents :  %s\n", strings.Join(filteredContents, ", "))
		}
		if len(filteredTags) > 0 {
			fmt.Fprintf(&b, "Filtered tags     :  #%s\n", strings.Join(filteredTags, " #"))
		}
		fmt.Fprintf(&b, "Tags      :  ")
		if len(post.Tags) > 0 {
			fmt.Fprintf(&b, "#%s", strings.Join(post.Tags, " #"))
		}
	}

	d.info.SetText(b.String())
}
