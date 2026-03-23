package helper

const BorderSide string = "borderSide"
const BorderTop string = "borderTop"
const BorderTopRight string = "borderTopRight"
const BorderTopLeft string = "borderTopLeft"
const BorderBottomLeft string = "borderBottomLeft"
const BorderBottomRight string = "borderBottomRight"

const BorderSideDouble string = "borderSideDouble"
const BorderTopDouble string = "borderTopDouble"
const BorderTopRightDouble string = "borderTopRightDouble"
const BorderTopLeftDouble string = "borderTopLeftDouble"
const BorderBottomLeftDouble string = "borderBottomLeftDouble"
const BorderBottomRightDouble string = "borderBottomRightDouble"

// INFO:Reference : https://en.wikipedia.org/wiki/Box-drawing_characters
// borderchars = { "─", "│", "─", "│", "╭", "╮", "╯", "╰" },
func Dictionary(key string) string {
	return map[string]string{
		"borderSide":        "│",
		"borderTop":         "─",
		"borderTopLeft":     "╭",
		"borderTopRight":    "╮",
		"borderBottomLeft":  "╰",
		"borderBottomRight": "╯",

		"borderSideDouble":        "║",
		"borderTopDouble":         "═",
		"borderTopLeftDouble":     "╔",
		"borderTopRightDouble":    "╗",
		"borderBottomLeftDouble":  "╚",
		"borderBottomRightDouble": "╝",
	}[key]
}

// Side,Top,TopLeft,TopRight,BottomLeft,BottomRight
func GetBorders() (string, string, string, string, string, string) {
	return Dictionary(BorderSide), Dictionary(BorderTop), Dictionary(BorderTopLeft), Dictionary(BorderTopRight), Dictionary(BorderBottomLeft), Dictionary(BorderBottomRight)
}
