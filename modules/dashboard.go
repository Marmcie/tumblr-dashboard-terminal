package modules

import (
	"math/rand/v2"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tumblr/tumblr.go"
)

type Dashboard struct {
	offset        int
	viewOffset    int
	selectedIndex int
	client        *TumblrClient

	Root        *tview.Flex
	postWrapper *tview.Flex
	postContent *tview.TextView

	leftSide  *tview.Flex
	rightSide *tview.Flex

	control  *tview.TextView
	postInfo *tview.TextView

	rows  []*tview.TextView
	posts []tumblr.Post

	colors map[string]tcell.Color
}

func (d *Dashboard) initEvents(app *tview.Application) *Dashboard {

	d.postWrapper.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			d.selectedIndex = fit(d.selectedIndex+1, len(d.posts))
			d.UpdateView()
			d.RenderPost()
		}
		if event.Rune() == 'k' {
			d.selectedIndex = fit(d.selectedIndex-1, len(d.posts))
			d.UpdateView()
			d.RenderPost()
		}

		if event.Key() == tcell.KeyEnter {
			app.SetFocus(d.postContent)
			d.postContent.SetBorderColor(tcell.ColorBlue)
			d.postWrapper.SetBorderColor(tcell.ColorWhite)
		}

		if event.Rune() == 'l' {
			app.SetFocus(d.postContent)
			d.postContent.SetBorderColor(tcell.ColorBlue)
			d.postWrapper.SetBorderColor(tcell.ColorWhite)
		}

		if event.Rune() == 'r' {
			d.Update()
		}
		return event
	})

	d.postContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			app.SetFocus(d.postWrapper)
			d.postContent.SetBorderColor(tcell.ColorWhite)
			d.postWrapper.SetBorderColor(tcell.ColorBlue)
		}

		return event
	})

	d.Root.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'o' {
			post := d.GetSelectedPost()
			openInBrowser(post.ShortUrl)
		}
		if event.Rune() == 'q' {
			app.Stop()
		}

		return event
	})

	return d
}

func (d *Dashboard) Update() {
	d.UpdateFeed()
	d.UpdateView()
	d.offset++
}

func (d *Dashboard) UpdateFeed() {
	res := d.client.GetDashboard(d.offset)

	colors := [...]tcell.Color{
		tcell.ColorYellowGreen,
		tcell.ColorLawnGreen,
		tcell.ColorDarkOrange,
		tcell.ColorFloralWhite,
		tcell.ColorAqua,
		tcell.ColorBlue,
		tcell.ColorDodgerBlue,
	}

	for n := range res {
		post := res[n]
		d.posts = append(d.posts, post)

		blogName := post.BlogName
		_, ok := d.colors[blogName]

		if !ok {
			d.colors[blogName] = colors[rand.Int()%len(colors)]
		}

		contents, err := htmltomarkdown.ConvertString(post.Body)
		length := min(len(contents), 120)
		postName := contents[:length]
		if err != nil {
		}
		row := tview.NewTextView().SetText(blogName + ":" + postName)
		row.SetTextColor(d.colors[blogName])
		d.rows = append(d.rows, row)
	}
}

// Render the post list window
func (d *Dashboard) UpdateView() {

	count := len(d.rows)

	_, _, _, h := d.postWrapper.GetInnerRect()

	// INFO:Moving the window down by cursor
	if d.selectedIndex >= h+d.viewOffset {
		d.viewOffset = d.selectedIndex - h + 1
	}

	// INFO:Reset the window when cursor is back up
	if d.selectedIndex <= 0 {
		d.viewOffset = 0
	}

	// INFO:Move the window up
	if d.selectedIndex < d.viewOffset {
		d.viewOffset = d.selectedIndex
	}

	d.postWrapper.Clear()

	for v := range max(min(h, count), 20) {
		d.postWrapper.AddItem(d.rows[v+d.viewOffset], 1, 0, false)
	}

	prevColor := d.colors[d.posts[fit(d.selectedIndex-1, count)].BlogName]
	nextColor := d.colors[d.posts[fit(d.selectedIndex+1, count)].BlogName]

	d.rows[d.selectedIndex].SetBackgroundColor(tcell.ColorWhiteSmoke)
	d.rows[d.selectedIndex].SetTextColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex+1, count)].SetBackgroundColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex+1, count)].SetTextColor(nextColor)
	d.rows[fit(d.selectedIndex-1, count)].SetBackgroundColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex-1, count)].SetTextColor(prevColor)
}

// Render the selected post into the post window
func (d *Dashboard) RenderPost() {
	post := d.posts[d.selectedIndex]
	contents := post.Body
	if post.Format == "html" {
		contents, _ = htmltomarkdown.ConvertString(post.Body)
	}

	d.postContent.SetText(contents)

	info := ""
	info += "Date      :  " + post.Date + "\n"
	info += "URL       :  " + post.ShortUrl + "\n"
	info += "Blog name :  " + post.BlogName + "\n"
	info += "Tags      :  "
	if len(post.Tags) > 0 {
		info += "#"
		info += strings.Join(post.Tags, " #")
	}

	d.postInfo.SetText(info)
}

func (d *Dashboard) GetSelectedPost() tumblr.Post {
	return d.posts[d.selectedIndex]
}

func NewDashboard(client *TumblrClient, app *tview.Application) *Dashboard {
	d := Dashboard{}
	d.offset = 0
	d.viewOffset = 0
	d.selectedIndex = 0
	d.client = client

	d.Root = tview.NewFlex()

	//INFO: Left side
	d.leftSide = tview.NewFlex().SetDirection(tview.FlexRow)

	d.postWrapper = tview.NewFlex().SetDirection(tview.FlexRow)
	d.postWrapper.SetTitle("Post list")
	d.postWrapper.SetBorder(true)
	d.postWrapper.SetBorderColor(tcell.ColorBlue)

	d.leftSide.AddItem(d.postWrapper, 0, 4, true)

	d.control = tview.NewTextView()
	d.control.SetBorder(true)
	d.control.SetTitle("Control")
	d.control.SetTextColor(tcell.ColorDarkOrange)
	d.control.SetText(controlText())

	d.leftSide.AddItem(d.control, 0, 1, false)

	//INFO: Right side
	d.rightSide = tview.NewFlex().SetDirection(tview.FlexRow)

	d.postContent = tview.NewTextView().SetScrollable(true)
	d.postContent.SetBorder(true)
	d.postContent.SetTitle("Post")
	d.rightSide.AddItem(d.postContent, 0, 4, false)

	d.postInfo = tview.NewTextView()
	d.postInfo.SetTitle("Post information")
	d.postInfo.SetBorder(true)
	d.rightSide.AddItem(d.postInfo, 0, 1, false)

	d.Root.AddItem(d.leftSide, 0, 1, true)
	d.Root.AddItem(d.rightSide, 0, 3, false)

	d.colors = map[string]tcell.Color{}
	d.initEvents(app)
	return &d
}

func abs(val int) int {
	if val >= 0 {
		return val
	} else {
		return val * -1
	}
}

func fit(val int, limit int) int {
	if val >= 0 {
		return val % limit
	} else {
		return limit + val
	}
}

func controlText() string {
	str := ""
	str += "j/k      :  Up/Down\n"
	str += "r        :  Load more posts\n"
	str += "Enter/l  :  Focus post window\n"
	str += "h        :  Focus post list window\n"
	str += "q        :  Quit\n"

	return str
}
