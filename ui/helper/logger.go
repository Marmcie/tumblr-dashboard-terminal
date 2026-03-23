package helper

import "os"

func Log(text string) {
	filename:="./debug.log"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text+"\n"); err != nil {
		panic(err)
	}
}
