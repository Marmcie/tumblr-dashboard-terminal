package npf

import "github.com/mattn/go-runewidth"

type Blog struct {
	Name                   string
	Title                  string
	Description            string
	Url                    string
	Uuid                   string
	Updated                int64
	Avatar                 []Avatar
	Tumblrmart_accessories struct {
		Badges               []Badge
		Blue_checkmark_count int64
	}
	Can_show_badges    bool
	Active             bool
	Show_follow_action bool
}

func (b *Blog) GetName() string {
	if runewidth.StringWidth(b.Name) > 0 {
		return b.Name
	}

	return b.Title
}
