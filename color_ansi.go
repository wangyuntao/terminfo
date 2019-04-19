package terminfo

import (
	"errors"
	"fmt"
)

const (
	rgbFgFmt = "\x1b[38;2;%d;%d;%dm"
	rgbBgFmt = "\x1b[48;2;%d;%d;%dm"
)

var (
	errColorIllegalRgbVal = errors.New("color: illegal rgb value")
)

// The functions in this file may not work in some terminals

func isValidRgb(r, g, b int) bool {
	return r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255
}

func (ti *Terminfo) ColorRgbFg(r, g, b int) error {
	if !isValidRgb(r, g, b) {
		return errColorIllegalRgbVal
	}
	_, err := fmt.Printf(rgbFgFmt, r, g, b)
	return err
}

func (ti *Terminfo) ColorRgbBg(r, g, b int) error {
	if !isValidRgb(r, g, b) {
		return errColorIllegalRgbVal
	}
	_, err := fmt.Printf(rgbBgFmt, r, g, b)
	return err
}
