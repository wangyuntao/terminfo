package main

import (
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	err = ti.Text().Underline().Italic().ColorFg(terminfo.ColorRed).Println("hello")
	failIfErr(err)
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
