package terminfo

import (
	"testing"
)

func TestColorValues(t *testing.T) {
	if ColorDefault >= 0 {
		t.Fail()
	}

	cs := []int{
		ColorBlack,
		ColorRed,
		ColorGreen,
		ColorYellow,
		ColorBlue,
		ColorMagenta,
		ColorCyan,
		ColorWhite,
		ColorBrightBlack,
		ColorBrightRed,
		ColorBrightGreen,
		ColorBrightYellow,
		ColorBrightBlue,
		ColorBrightMagenta,
		ColorBrightCyan,
		ColorBrightWhite,
	}

	for i, c := range cs {
		if i != c {
			t.Fail()
		}
	}
}
