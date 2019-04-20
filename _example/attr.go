package main

import (
	"fmt"
	"log"

	"github.com/wangyuntao/terminfo"
)

func main() {
	ti, err := terminfo.LoadEnv()
	failIfErr(err)

	err = standout(ti)
	failIfErr(err)

	err = underline(ti)
	failIfErr(err)

	err = reverse(ti)
	failIfErr(err)

	err = blink(ti)
	failIfErr(err)

	err = dim(ti)
	failIfErr(err)

	err = bold(ti)
	failIfErr(err)

	err = italics(ti)
	failIfErr(err)
}

func standout(ti *terminfo.Terminfo) error {
	err := ti.EnterStandoutMode()
	if err != nil {
		return err
	}
	fmt.Println("standout")
	return ti.ExitAttributeMode()
}

func underline(ti *terminfo.Terminfo) error {
	err := ti.EnterUnderlineMode()
	if err != nil {
		return err
	}
	fmt.Println("underline")
	return ti.ExitAttributeMode()
}

func reverse(ti *terminfo.Terminfo) error {
	err := ti.EnterReverseMode()
	if err != nil {
		return err
	}
	fmt.Println("reverse")
	return ti.ExitAttributeMode()
}

func blink(ti *terminfo.Terminfo) error {
	err := ti.EnterBlinkMode()
	if err != nil {
		return err
	}
	fmt.Println("blink")
	return ti.ExitAttributeMode()
}

func dim(ti *terminfo.Terminfo) error {
	err := ti.EnterDimMode()
	if err != nil {
		return err
	}
	fmt.Println("dim")
	return ti.ExitAttributeMode()
}

func bold(ti *terminfo.Terminfo) error {
	err := ti.EnterBoldMode()
	if err != nil {
		return err
	}
	fmt.Println("bold")
	return ti.ExitAttributeMode()
}

func italics(ti *terminfo.Terminfo) error {
	err := ti.EnterItalicsMode()
	if err != nil {
		return err
	}
	fmt.Println("italics")
	return ti.ExitAttributeMode()
}

func failIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
