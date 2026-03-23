package ui

import ()

const ColorBG = "ColorBG"
const ColorFocus = "ColorFocus"
const ColorWhite = "ColorWhite"
const ColorH1 = "ColorH1"
const ColorH2 = "ColorH2"
const ColorImage = "ColorImage"
const ColorQuote = "ColorQuote"

func GetColorStr(key string) string {
	return map[string]string{
		"ColorBG": "#060616",
		"ColorFocus": "#135366",
		"ColorWhite": "#ffffff",
		"ColorH1":    "#40f0f0",
		"ColorH2":    "#a0f000",
		"ColorImage": "#40a0f0",
		"ColorQuote": "#f0f000",
	}[key]
}
