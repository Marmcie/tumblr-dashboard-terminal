package npf

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"math"
	"net/http"

	"github.com/mattn/go-runewidth"
)

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

var blogColors map[string]string = map[string]string{}

func (b *Blog) GetName() string {
	if runewidth.StringWidth(b.Name) > 0 {
		return b.Name
	}

	return b.Title
}

// Fetch blog color based on profile picture.
//
// Only available if use_blog_avatar_color option is true in config file
func (b *Blog) GetBlogColor() string {
	if col, ok := blogColors[b.Name]; ok {
		return col
	}

	// Check if no Avatar is detected
	if len(b.Avatar) == 0 {
		blogColors[b.Name] = "#ffffff"
		return blogColors[b.Name]
	}

	avatar := b.Avatar[len(b.Avatar)-1]
	resp, err := http.Get(avatar.Url)

	// Check if HTTP error occurs
	if err != nil {
		blogColors[b.Name] = "#ffffff"
		return blogColors[b.Name]
	}
	defer resp.Body.Close()

	// Check if HTTP error occurs
	if resp.StatusCode != 200 {
		blogColors[b.Name] = "#ffffff"
		return blogColors[b.Name]
	}

	image.RegisterFormat("jpg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	img, _, err := image.Decode(resp.Body)

	// Check if file format is unsupported
	if err != nil {
		blogColors[b.Name] = "#ffffff"
		return blogColors[b.Name]
	}

	// Get the color at the center of the screen
	red, green, blue, _ := img.At(img.Bounds().Max.X/2, img.Bounds().Max.Y/2).RGBA()

	// Map each color into 0~1
	fRed := float64(float64(red) / 65535)
	fGreen := float64(float64(green) / 65535)
	fBlue := float64(float64(blue) / 65535)

	// Square all colors for brighter color
	fRed *= fRed
	fGreen *= fGreen
	fBlue *= fBlue

	// Normalize all colors
	mag := math.Sqrt(float64((fRed * fRed) + (fGreen * fGreen) + (fBlue * fBlue)))
	fRed *= 1 / mag
	fGreen *= 1 / mag
	fBlue *= 1 / mag

	blogColors[b.Name] = fmt.Sprintf("#%02x%02x%02x", int(fRed*255), int(fGreen*255), int(fBlue*255))
	return blogColors[b.Name]
}
