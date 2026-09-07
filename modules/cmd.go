package modules

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

func HandleArgs() {
	//#region Version flag initialization
	version := flag.Bool("v", false, "Print current version")
	versionAlt := flag.Bool("version", false, "Print current version")
	//#endregion

	flag.Parse()

	//#region Version flag handling
	if *version || *versionAlt {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			panic("Failed to read build info")
		}
		versionStr := info.Main.Version
		if len(versionStr) == 0 {
			fmt.Println("Version missing")
		} else {
			fmt.Println(versionStr)
		}
		os.Exit(0)
	}
	//#endregion
}
