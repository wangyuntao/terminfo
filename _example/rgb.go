package main

import (
	"fmt"
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	n := 26
	for r := 0; r < 256; r += n {
		for g := 0; g < 256; g += n {
			for b := 0; b < 256; b += n {
				rgb(r, g, b, ti)
			}
		}
		err = ti.CursorDown()
		failIfErr(err)
	}
}

func rgb(r, g, b int, ti *terminfo.Terminfo) {
	err := ti.ColorRgbBg(r, g, b)
	failIfErr(err)
	fmt.Print(" ")
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
