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

	fmt.Println("name:", ti.Name())
	fmt.Println("path:", ti.Path())
	fmt.Println("desc:", ti.Desc())

	if s, ok := ti.GetStr(terminfo.ClearScreen); ok {
		fmt.Print(string(s))
	} else {
		log.Fatal("Absent")
	}

	_, err = ti.Printf(terminfo.CursorAddress, 20, 50)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello, World!")

}
