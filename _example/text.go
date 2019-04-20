package main

import (
	"fmt"
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	err = ti.Text().Underline().Italic().ColorFg(terminfo.ColorRed).Print("hello")
	failIfErr(err)

	fmt.Println(" world")
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
