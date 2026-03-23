package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

func GetColorStr(key string) string {
	return map[string]string{
		"ColorFocus": "#a0a4fa",
		"ColorWhite": "#ffffff",
		"ColorH1":    "#40f0f0",
		"ColorH2":    "#a0f000",
		"ColorImage": "#40a0f0",
		"ColorQuote": "#f0f000",
	}[key]
}
func GetColor(key string) color.Color {
	return lipgloss.Color(GetColorStr(key))
}
