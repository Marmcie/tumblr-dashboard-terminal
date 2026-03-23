package dashboard

import (
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"

	tea "github.com/charmbracelet/bubbletea"
)

type Feed struct {
	listElem  *component.Selectlist
	dashboard *Dashboard
}

func NewFeed(dashboard *Dashboard) *Feed {
	f := &Feed{}
	f.listElem = component.NewSelectlist("Feed")
	f.listElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetWidthInherit(true)
	f.dashboard=dashboard

	f.InitEvents()
	return f
}

func (f *Feed) InitEvents() {

	f.listElem.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "l":
				f.dashboard.FocusContents()
			}
		}
	})

	f.listElem.AddEventListener("onChange", func(msg tea.Msg, i int) {
		f.listElem.RunSelectedOption()
	})

}

func (f *Feed) AddPosts(posts []modules.Post) {
	for _, post := range posts {
		item := component.NewBox("Feed post")
		item.SetBorder(true).
			SetBorderPadding(1).
			SetBorderCorner(true).
			SetH(4).
			SetWidthInherit(true)

		blogName := component.NewLine("User name : " + post.Blog.Name)
		date := component.NewLine("User name : " + post.Date)
		blogName.SetText(post.Blog.Name)
		blogName.SetWidthInherit(true)
		date.SetText(post.Date)
		date.SetWidthInherit(true)

		item.AddChild(blogName)
		item.AddChild(date)
		f.listElem.AddOption(item, func() {
			f.dashboard.DisplayPost(post)
		})
	}
}
func (f *Feed) Focus(){
	f.listElem.Focus()
}
