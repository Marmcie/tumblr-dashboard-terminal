package main

import (
	"fmt"
	"os"
	"tumblr-dt/dashboard"

	tea "charm.land/bubbletea/v2"
)

func main() {

	dashboard := dashboard.NewDashboard()

	p := tea.NewProgram(dashboard.GetCore())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Bubbletea event loop error: %v", err)
		os.Exit(1)
	}
}

	//
	// r := "👨🏻‍🦰"
	// println("Char : " + string(r))
	// println("len  : " + strconv.Itoa(len(string(r))))
	// println("stringWidth  : " + strconv.Itoa(runewidth.StringWidth(string(r))))
	// println("\n")
	// println("================")
	//
	// println("Rune rendering")
	// for _, v := range "👨🏻‍🦰" {
	// 	println("Char : " + string(v))
	// 	println("len  : " + strconv.Itoa(len(string(v))))
	// 	println("runeWidth  : " + strconv.Itoa(runewidth.RuneWidth(v)))
	// 	println("stringWidth  : " + strconv.Itoa(runewidth.StringWidth(string(v))))
	// }
	//
	// println("\n")
	// println("================")
	// println("String split rendering")
	// for v := range strings.SplitSeq("👨🏻‍🦰", "") {
	// 	println("Char : " + string(v))
	// 	println("len  : " + strconv.Itoa(len(string(v))))
	// 	// println("runeWidth  : " + strconv.Itoa(runewidth.RuneWidth(v)))
	// 	println("stringWidth  : " + strconv.Itoa(runewidth.StringWidth(string(v))))
	// }
	//
	// println("\n")
	// println("================")
	// println("String rendering")
	// for _, v := range "👨🏻‍🦰" {
	// 	println("Char : " + string(v))
	// 	println("len  : " + strconv.Itoa(len(string(v))))
	// 	println("runeWidth  : " + strconv.Itoa(runewidth.RuneWidth(v)))
	// 	println("stringWidth  : " + strconv.Itoa(runewidth.StringWidth(string(v))))
	// }
	//
	// panic(1)
