package helper

const BorderSide string = "borderSide"
const BorderTop string = "borderTop"
const BorderTopRight string = "borderTopRight"
const BorderTopLeft string = "borderTopLeft"
const BorderBottomLeft string = "borderBottomLeft"
const BorderBottomRight string = "borderBottomRight"

// borderchars = { "─", "│", "─", "│", "╭", "╮", "╯", "╰" },
func Dictionary(key string) string {
	return map[string]string{
		"borderSide":        "│",
		"borderTop":         "─",
		"borderTopLeft":     "╭",
		"borderTopRight":    "╮",
		"borderBottomLeft":  "╰",
		"borderBottomRight": "╯",
	}[key]
}

// Side,Top,TopLeft,TopRight,BottomLeft,BottomRight
func GetBorders() (string, string, string, string, string, string) {
	return Dictionary(BorderSide), Dictionary(BorderTop), Dictionary(BorderTopLeft), Dictionary(BorderTopRight), Dictionary(BorderBottomLeft), Dictionary(BorderBottomRight)
}
