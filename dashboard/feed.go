package dashboard

import (
	"tumblr-dt/npf"
	"tumblr-dt/ui"
	component "tumblr-dt/ui/components"

	tea "charm.land/bubbletea/v2"
)

type Feed struct {
	listElem  *component.Selectlist
	dashboard *Dashboard
	posts     []npf.Post
}

func NewFeed(dashboard *Dashboard) *Feed {
	f := &Feed{}
	f.listElem = component.NewSelectlist("Feed")
	f.listElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetWidthInherit(true)
	f.dashboard = dashboard
	f.listElem.SetBorderLabel("BottomRight", "? For keybind")
	

	f.InitEvents()
	return f
}

func (f *Feed) InitEvents() {

	f.listElem.AddEventListener("onUpdate", func(msg tea.Msg, i int) {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter", "l":
				f.dashboard.FocusContents()
				f.listElem.RunSelectedOption()
			case "j":
				f.listElem.IncrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdateSelectedOptionBorder()
			case "k":
				f.listElem.DecrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdateSelectedOptionBorder()
			}
		}
	})

}

func (f *Feed) UpdateSelectedOptionBorder() {
	children := f.listElem.GetChildren()
	if len(children) == 0 {
		return
	}
	if children[f.listElem.Cursor] != nil {
		children[f.listElem.Cursor].SetBackground(ui.GetColorStr(ui.ColorFocus))
		children[f.listElem.Cursor].SetForeground(ui.GetColorStr(ui.ColorWhite))
	}
	if f.listElem.Cursor > 0 {
		children[f.listElem.Cursor-1].ClearBackground()
		children[f.listElem.Cursor-1].ClearForeground()
	}

	if f.listElem.Cursor < len(children)-1 {
		children[f.listElem.Cursor+1].ClearBackground()
		children[f.listElem.Cursor+1].ClearForeground()
	}
}

func (f *Feed) GetSelectedPost() npf.Post {
	return f.posts[f.listElem.Cursor]
}

func (f *Feed) AddPosts(posts []npf.Post) {
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
		blogName.SetForeground(ui.GetColorStr(ui.ColorH1))

		summary := component.NewLine("Post summary")
		summary.SetText(post.GetSummary())
		summary.SetWidthInherit(true)

		item.AddChild(blogName)
		item.AddChild(summary)
		f.listElem.AddOption(item, func() {
			f.dashboard.DisplayPost(post)
		})
	}
	f.UpdateSelectedOptionBorder()
}

func (f *Feed) Focus() {
	f.listElem.Focus()
}
