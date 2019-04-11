package terminfo

import "testing"
import "fmt"

func TestTermSize(t *testing.T) {
	r, c, err := TermSize()
	if err != nil {
		fmt.Println("TestTermSize error", err)
		return
	}
	fmt.Println("TestTermSize", r, c)
}

func TestTermFileIn(t *testing.T) {
	f, err := TermFileIn()
	if err != nil {
		fmt.Println("TestTermFileIn error", err)
		return
	}
	defer TermFileClose(f)

	fmt.Println("TestTermFileIn", f.Name(), f.Fd())
}

func TestTermFileOut(t *testing.T) {
	f, err := TermFileOut()
	if err != nil {
		fmt.Println("TestTermFileOut error", err)
		return
	}
	defer TermFileClose(f)

	fmt.Println("TestTermFileOut", f.Name(), f.Fd())
}
