package main

import (
	"math/rand/v2"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/tumblr/tumblr.go"
	"github.com/tumblr/tumblrclient.go"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dashboard struct {
	offset        int
	viewOffset    int
	selectedIndex int
	client        *tumblrclient.Client

	flex        *tview.Flex
	postWrapper *tview.Flex
	postContent *tview.TextView

	rows  []*tview.TextView
	posts []tumblr.PostInterface

	colors map[string]tcell.Color
}

func newDashboard(client *tumblrclient.Client, app *tview.Application) *Dashboard {
	d := Dashboard{}
	d.offset = 0
	d.viewOffset = 0
	d.selectedIndex = 0
	d.client = client

	d.flex = tview.NewFlex()
	d.postWrapper = tview.NewFlex().SetDirection(tview.FlexRow)
	d.postContent = tview.NewTextView().SetScrollable(true)

	d.postWrapper.SetBorder(true)
	d.postContent.SetBorder(true)

	d.flex.AddItem(d.postWrapper, 0, 1, true)
	d.flex.AddItem(d.postContent, 0, 3, false)

	d.colors = map[string]tcell.Color{}
	d.initEvents(app)
	return &d
}

func (d *Dashboard) initEvents(app *tview.Application) *Dashboard {

	d.postWrapper.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			d.selectedIndex = fit(d.selectedIndex+1, len(d.posts))
			d.updateView()
		}
		if event.Rune() == 'k' {
			d.selectedIndex = fit(d.selectedIndex-1, len(d.posts))
			d.updateView()
		}

		if event.Key() == tcell.KeyEnter {
			d.renderPost()
		}

		if event.Rune() == 'l' {
			d.renderPost()
			app.SetFocus(d.postContent)
		}

		if event.Rune() == 'r' {
			d.updateFeed()
			d.addPostsToList()
			d.offset = d.offset + 1
		}
		if event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	d.postContent.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'h' {
			app.SetFocus(d.postWrapper)
		}
		return event
	})

	return d
}

func main() {

	godotenv.Load("D:/Active projects/tumblr2/.env")
	var consumerkey = os.Getenv("consumer_key")
	var secretkey = os.Getenv("secret_key")
	var oauthtoken = os.Getenv("oauth_token")
	var oauthsecret = os.Getenv("oauth_secret")
	client := tumblrclient.NewClientWithToken(consumerkey, secretkey, oauthtoken, oauthsecret)
	app := tview.NewApplication()

	dashboard := newDashboard(client, app)

	dashboard.updateFeed()
	dashboard.addPostsToList()
	dashboard.updateView()

	if err := app.SetRoot(dashboard.flex, true).Run(); err != nil {
		panic(err)
	}

}

func (d *Dashboard) updateFeed() {
	var params = url.Values{}
	params.Set("offset", strconv.Itoa(d.offset*20))
	res, err := d.client.GetDashboardWithParams(params)
	if err != nil {
	}

	colors := [...]tcell.Color{
		tcell.ColorYellowGreen,
		tcell.ColorBlueViolet,
		tcell.ColorLawnGreen,
		tcell.ColorDarkOrange,
		tcell.ColorFloralWhite,
	}

	for n := range res.Posts {
		post := res.Posts[n]
		d.posts = append(d.posts, post)

		blogName := post.GetSelf().BlogName
		_, ok := d.colors[blogName]

		if !ok {
			d.colors[blogName] = colors[rand.Int()%len(colors)]
		}

	}
}

func (d *Dashboard) addPostsToList() {
	for i := range 20 {
		index := (d.offset * 20) + i
		post := d.posts[index]

		contents, err := htmltomarkdown.ConvertString(post.GetSelf().Body)
		length := min(len(contents), 120)
		postName := contents[:length]
		blogName := post.GetSelf().GetSelf().BlogName
		if err != nil {
		}
		row := tview.NewTextView().SetText(blogName + ":" + postName)
		row.SetTextColor(d.colors[blogName])
		d.rows = append(d.rows, row)
		d.postWrapper.AddItem(row, 1, 0, false)
	}
}

func (d *Dashboard) updateView() {

	count := len(d.rows)

	_, _, _, h := d.postWrapper.GetInnerRect()

	if d.selectedIndex > h+d.viewOffset-1 {
		d.viewOffset = d.selectedIndex - h
	}

	if d.selectedIndex <= 0 {
		d.viewOffset = 0
	}

	if d.selectedIndex < d.viewOffset {
		d.viewOffset = d.selectedIndex
	}

	d.postWrapper.Clear()

	for v := range max(min(h, count), 20) {
		d.postWrapper.AddItem(d.rows[v+d.viewOffset], 1, 0, false)
	}

	prevColor := d.colors[d.posts[fit(d.selectedIndex-1, count)].GetSelf().BlogName]
	nextColor := d.colors[d.posts[fit(d.selectedIndex+1, count)].GetSelf().BlogName]

	d.rows[d.selectedIndex].SetBackgroundColor(tcell.ColorWhiteSmoke)
	d.rows[d.selectedIndex].SetTextColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex+1, count)].SetBackgroundColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex+1, count)].SetTextColor(nextColor)
	d.rows[fit(d.selectedIndex-1, count)].SetBackgroundColor(tcell.ColorBlack)
	d.rows[fit(d.selectedIndex-1, count)].SetTextColor(prevColor)
}

func (d *Dashboard) renderPost() {

	contents, err := htmltomarkdown.ConvertString(d.posts[d.selectedIndex].GetSelf().Body)
	if err != nil {
	}
	d.postContent.SetText(contents)
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
