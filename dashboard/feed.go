package dashboard

import (
	"fmt"
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/component"

	tea "charm.land/bubbletea/v2"
)

type Feed struct {
	listElem         *component.Selectlist
	dashboard        *Dashboard
	posts            []*npf.Post
	prev             string
	showFilteredPost bool
}

func NewFeed(dashboard *Dashboard) *Feed {
	f := &Feed{}
	f.listElem = component.NewSelectlist("Feed")
	f.listElem.SetBorder(true).SetWidthInherit(true)
	f.dashboard = dashboard
	f.listElem.SetBorderLabel("BottomRight", "? For keybind")
	f.listElem.SetSelectedOptionBackground(ui.GetColorStr(ui.ColorFocus))
	f.listElem.SetSelectedOptionForeground(ui.GetColorStr(ui.ColorWhite))
	f.listElem.SetBorderFocusForeground(ui.GetColorStr(ui.ColorFocusBorder))
	f.showFilteredPost = false

	f.InitEvents()
	return f
}

func (f *Feed) UpdatePostCounter() {
	f.listElem.SetBorderLabel("BottomLeft", fmt.Sprintf("%d/%d", f.listElem.Cursor+1, len(f.listElem.GetChildren())))
}

func (f *Feed) InitEvents() {

	f.listElem.AddEventListener("onUpdate", func(msg tea.Msg) {
		switch msg := msg.(type) {

		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter", "l":
				if f.listElem.Cursor < len(f.listElem.GetChildren()) {
					f.showFilteredPost = true
					f.dashboard.FocusContents()
					f.listElem.RunSelectedOption()
				}
			case "j":
				if f.listElem.Cursor == len(f.listElem.GetChildren())-1 {
					done := make(chan bool)
					go f.dashboard.LoadPosts(done)
				}
				f.showFilteredPost = false
				f.listElem.IncrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdatePostCounter()
			case "k":
				f.showFilteredPost = false
				f.listElem.DecrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdatePostCounter()
			case "G":
				f.showFilteredPost = false
				f.listElem.SetCursor(len(f.posts) - 1)
				f.listElem.RunSelectedOption()
				f.UpdatePostCounter()
			case "g":
				if f.prev == "g" {
					f.showFilteredPost = false
					f.listElem.SetCursor(0)
					f.listElem.RunSelectedOption()
					f.UpdatePostCounter()
				}
			}
			f.prev = msg.String()
		}
	}, true)

}

func (f *Feed) GetSelectedPost() *npf.Post {
	return f.posts[f.listElem.Cursor]
}

func (f *Feed) ClearPosts() {

	for _, p := range f.posts {
		p.RemoveRenderResult()
	}

	f.posts = []*npf.Post{}
}

func (f *Feed) AddPosts(posts []*npf.Post) {
	for _, post := range posts {
		f.posts = append(f.posts, post)
		item := component.NewBox("Feed post")
		item.SetBorder(true).
			SetBorders(false, true, false, false).
			SetBorderCorner(false).
			SetH(3).
			SetWidthInherit(true)

		blogName := component.NewLine("User name : " + post.Blog.Name)
		blogName.SetText(post.Blog.Name)

		blogName.SetWidthInherit(true)
		if f.dashboard.config.Use_blog_avatar_color {
			blogName.SetForeground(post.Blog.GetBlogColor())
		} else {
			blogName.SetForeground(ui.GetColorStr(ui.ColorH1))
		}

		summary := component.NewLine("Post summary")
		summary.SetText(post.GetSummary())
		summary.SetWidthInherit(true)

		if post.IsFiltered {
			summary.SetText("Filtered")
			summary.SetForeground(ui.GetColorStr(ui.ColorFiltered))
		}

		item.AddChild(blogName)
		item.AddChild(summary)
		f.listElem.AddOption(item, func() {
			f.dashboard.DisplayPost(post, f.showFilteredPost)
		})
	}
	f.listElem.SetCursor(f.listElem.Cursor)
	f.UpdatePostCounter()
}

func (f *Feed) Focus() {
	f.listElem.Focus()
}
