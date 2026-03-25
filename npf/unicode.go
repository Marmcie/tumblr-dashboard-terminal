package npf

import (
	"bytes"
	"strings"

	"github.com/mattn/go-runewidth"
	"golang.org/x/text/width"
)

// INFO: SUS EMOJIS LIST : 👨‍👩‍👧✌️
// String width
// ✌️ = 1 2
// 👨‍👩‍👧 = 2 0 2 0 2
func RenderUnicode(str string) string {
	con := runewidth.Condition{
		EastAsianWidth:     true,
		StrictEmojiNeutral: true,
	}
	var result bytes.Buffer
	for v := range strings.SplitSeq(str, "") {
		info, _ := width.LookupString(v)

		if info.Kind() == width.EastAsianFullwidth || info.Kind() == width.EastAsianWide {
			/**INFO:
					Use StringWidth instead of RuneWidth because sometimes
			  	rune count and actual string width are different
			*/
			for range con.StringWidth(v) - 1 {
				// INFO: Output 0 width character to account for full width chars
				result.WriteRune('\u200b')
			}
		}

		result.WriteString(v)
	}
	return result.String()
}
