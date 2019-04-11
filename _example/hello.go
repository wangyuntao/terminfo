package main

import (
	"fmt"
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	err = ti.ClearScreen()
	failIfErr(err)

	rainbow(ti, "HELLO", terminfo.ColorRed)
	err = ti.CursorAddress(1, 5)
	failIfErr(err)

	rainbow(ti, "WORLD", terminfo.ColorRed)
	fmt.Println()
}

func rainbow(ti *terminfo.Terminfo, s string, c int) {
	for _, r := range s {
		err := ti.ColorFg(c)
		failIfErr(err)

		fmt.Printf("%c", r)
		c++
	}
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
