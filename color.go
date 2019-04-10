package terminfo

import "errors"

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
	errColorMaxColorsCapAbsent = errors.New("capability colors absent")
	errColorIllegalFgColor     = errors.New("illegal fg color")
	errColorIllegalBgColor     = errors.New("illegal bg color")
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
