package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	for {
		err = refresh(ti)
		failIfErr(err)
		time.Sleep(20 * time.Millisecond)
	}
}

func refresh(ti *terminfo.Terminfo) error {
	row, col, err := terminfo.TermSize()
	if err != nil {
		return err
	}

	maxColors, ok := ti.GetNum(terminfo.MaxColors)
	if !ok {
		return terminfo.ErrCapAbsent
	}

	for r := 0; r < row; r++ {
		err = ti.CursorAddress(r, 0)
		if err != nil {
			return err
		}

		for c := 0; c < col; c++ {
			err = ti.ColorFg(rand.Intn(maxColors))
			if err != nil {
				return err
			}
			fmt.Printf("%c", 33+rand.Intn(127-33))
		}
	}
	return nil
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
