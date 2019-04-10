package main

import (
	"fmt"
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = ti.ClearScreen()
	if err != nil {
		log.Fatal(err)
	}

	err = ti.ColorFg(terminfo.ColorYellow)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("hello ")

	err = ti.ColorFg(terminfo.ColorBlue)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("world")
}
