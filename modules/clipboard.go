package modules

import (
	"context"

	"golang.design/x/clipboard"
)

var clipboardInitialized = false

func CopyToClipboard(str string) {
	if !clipboardInitialized {
		initClipboard()
	}
	if clipboardInitialized {
		ctx := context.Background()
		clipboard.Write(ctx, clipboard.FmtText, []byte(str))
		ctx.Done()
	}

}
func initClipboard() {
	err := clipboard.Init()
	if err == nil {
		clipboardInitialized = true
	}
}
