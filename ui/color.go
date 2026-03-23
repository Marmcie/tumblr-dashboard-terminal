package ui

import (
	"tumblr-dt/modules"
)

const ColorBG = "ColorBG"
const ColorFocus = "ColorFocus"
const ColorFocusBorder = "ColorFocusBorder"
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
		"ColorBG":          config.Colors.Bg,
		"ColorFocus":       config.Colors.Focus,
		"ColorFocusBorder": config.Colors.Focus_border,
		"ColorWhite":       config.Colors.White,
		"ColorGrey":        config.Colors.Grey,
		"ColorH1":          config.Colors.H1,
		"ColorH2":          config.Colors.H2,
		"ColorImage":       config.Colors.Image,
		"ColorQuote":       config.Colors.Quote,
		"ColorFiltered":    config.Colors.Filtered,
		"ColorBlacklisted": config.Colors.Blacklisted,
	}[key]
}
