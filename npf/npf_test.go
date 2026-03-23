package npf_test

import (
	"testing"
	"tumblr-dt/npf"

	"github.com/forPelevin/gomoji"
)

func TestEmojistrip(t *testing.T) {
	res := npf.RenderUnicode("👨‍👩‍👧✌️")
	if gomoji.ContainsEmoji(res) {
		t.Errorf("Failed to strip emoji. Got : %s", res)
	}
}

