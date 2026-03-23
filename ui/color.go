package ui

import (
	"tumblr-dt/modules"
)

const ColorBG = "ColorBG"
const ColorFocus = "ColorFocus"
const ColorWhite = "ColorWhite"
const ColorGrey = "ColorGrey"
const ColorH1 = "ColorH1"
const ColorH2 = "ColorH2"
const ColorImage = "ColorImage"
const ColorQuote = "ColorQuote"
const ColorFiltered = "ColorFiltered"
const ColorBlacklisted = "ColorBlacklisted"

func GetColorStr(key string) string {
	config := modules.GetConfig()
	return map[string]string{
		"ColorBG":       config.Colors.Bg,
		"ColorFocus":    config.Colors.Focus,
		"ColorWhite":    config.Colors.White,
		"ColorGrey":     config.Colors.Grey,
		"ColorH1":       config.Colors.H1,
		"ColorH2":       config.Colors.H2,
		"ColorImage":    config.Colors.Image,
		"ColorQuote":    config.Colors.Quote,
		"ColorFiltered": config.Colors.Filtered,
	}[key]
}
