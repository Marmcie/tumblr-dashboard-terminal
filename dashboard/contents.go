package dashboard

import (
	"strings"
	"tumblr-dt/modules"
	component "tumblr-dt/ui/components"
)

type Contents struct {
	contentElem *component.Text
}

func NewContents() *Contents {
	f := &Contents{}
	f.contentElem = component.NewText("Contents")
	f.contentElem.SetBorder(true).SetBorderPadding(1).SetBorderCorner(true).SetHeightInherit(true).SetWidthInherit(true)

	return f
}

func (f *Contents) DisplayPost(post modules.Post) {
	f.contentElem.SetText(strings.Join(post.Render(), "\n"))
}
