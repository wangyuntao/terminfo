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

	err = ti.ColorFg(terminfo.Red)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("hello")

	err = ti.ColorReset()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("world")
}
