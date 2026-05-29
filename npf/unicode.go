package npf

import (
	"bytes"
	"github.com/rivo/uniseg"
)

func RenderUnicode(str string) string {
	var result bytes.Buffer
	itr := uniseg.NewGraphemes(str)
	for itr.Next() {
		char := itr.Str()
		result.WriteString(char)
		for range uniseg.StringWidth(char) - 1 {
			result.WriteRune('\u200b')
		}
	}
	return result.String()
}
