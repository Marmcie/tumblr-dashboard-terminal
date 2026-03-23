package helper

import "github.com/charmbracelet/lipgloss"

func ApplyStyleToLine(text []string, style lipgloss.Style) []string {
	var res []string

	for i := range text {
		res = append(res, style.SetString(text[i]).String())
	}
	return res
}

func ApplyStyleToCanvas(canvas [][]string, style lipgloss.Style) [][]string {
	var res [][]string

	for _, line := range canvas {
		res = append(res, ApplyStyleToLine(line, style))
	}
	return res
}
