package terminfo

import "errors"
import "fmt"

const (
	ColorDefault = -1
)

const (
	ColorBlack = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite

	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

var (
	errColorMaxColorsCapAbsent = errors.New("color: capability colors absent")
	errColorIllegalFgColor     = errors.New("color: illegal fg color")
	errColorIllegalBgColor     = errors.New("color: illegal bg color")
	errColorIllegalRgbVal      = errors.New("color: illegal rgb value")
)

func (ti *Terminfo) Color(fg, bg int) error {
	mc, ok := ti.GetNum(MaxColors)
	if !ok {
		return errColorMaxColorsCapAbsent
	}

	if fg != ColorDefault {
		if fg < 0 || fg >= mc {
			return errColorIllegalFgColor
		}
		err := ti.Do(SetAForeground, fg)
		if err != nil {
			return err
		}
	}

	if bg != ColorDefault {
		if bg < 0 || bg >= mc {
			return errColorIllegalBgColor
		}
		err := ti.Do(SetABackground, bg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ti *Terminfo) ColorFg(fg int) error {
	return ti.Color(fg, ColorDefault)
}

func (ti *Terminfo) ColorBg(bg int) error {
	return ti.Color(ColorDefault, bg)
}

func (ti *Terminfo) ColorReset() error {
	return ti.Do(OrigPair)
}

// RGB TODO portable

const (
	rgbFgFmt = "\x1b[38;2;%d;%d;%dm"
	rgbBgFmt = "\x1b[48;2;%d;%d;%dm"
)

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
