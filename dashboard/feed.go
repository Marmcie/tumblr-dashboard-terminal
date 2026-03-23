package dashboard

import (
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Feed struct {
	listElem  *component.Selectlist
	dashboard *Dashboard
	posts     []modules.Post
}

func NewFeed(dashboard *Dashboard) *Feed {
	f := &Feed{}
	f.listElem = component.NewSelectlist("Feed")
	f.listElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetWidthInherit(true)
	f.dashboard = dashboard

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
			case "enter":
				f.dashboard.FocusContents()

			case "j":
				f.listElem.IncrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdateSelectedOptionBorder()
			case "k":
				f.listElem.DecrementCursor()
				f.listElem.RunSelectedOption()
				f.UpdateSelectedOptionBorder()

			case "o":
				post := f.GetSelectedPost()
				modules.OpenInBrowser(post.Short_url)
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
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00aaaa"))
		children[f.listElem.Cursor].SetBorderStyle(style)
		children[f.listElem.Cursor].SetDoubleBorder(true)
	}
	if f.listElem.Cursor > 0 {
		children[f.listElem.Cursor-1].ResetBorderStyle()
		children[f.listElem.Cursor-1].SetDoubleBorder(false)
	}

	if f.listElem.Cursor < len(children)-1 {
		children[f.listElem.Cursor+1].ResetBorderStyle()
		children[f.listElem.Cursor+1].SetDoubleBorder(false)
	}
}

func (f *Feed) GetSelectedPost() modules.Post {
	return f.posts[f.listElem.Cursor]
}

func (f *Feed) AddPosts(posts []modules.Post) {
	for _, post := range posts {
		f.posts = append(f.posts, post)
		item := component.NewBox("Feed post")
		item.SetBorder(true).
			SetBorderPadding(1).
			SetBorderCorner(true).
			SetH(4).
			SetWidthInherit(true)

		blogName := component.NewLine("User name : " + post.Blog.Name)
		blogName.SetText(post.Blog.Name)
		blogName.SetWidthInherit(true)
		blogName.SetStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a4fa")))

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
