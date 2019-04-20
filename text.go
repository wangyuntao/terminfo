package terminfo

import "fmt"

type Text struct {
	ti *Terminfo

	standout  bool
	underline bool
	reverse   bool
	dim       bool
	bold      bool
	italic    bool

	fg, bg int
}

func (ti *Terminfo) Text() *Text {
	return &Text{
		ti: ti,
		fg: ColorDefault,
		bg: ColorDefault,
	}
}

func (t *Text) Standout() *Text {
	t.standout = true
	return t
}

func (t *Text) Underline() *Text {
	t.underline = true
	return t
}

func (t *Text) Reverse() *Text {
	t.reverse = true
	return t
}

func (t *Text) Dim() *Text {
	t.dim = true
	return t
}

func (t *Text) Bold() *Text {
	t.bold = true
	return t
}

func (t *Text) Italic() *Text {
	t.italic = true
	return t
}

func (t *Text) ColorFg(fg int) *Text {
	t.fg = fg
	return t
}
func (t *Text) ColorBg(bg int) *Text {
	t.bg = bg
	return t
}

func (t *Text) Do() error {
	if t.standout {
		err := t.ti.EnterStandoutMode()
		if err != nil {
			return err
		}
	}

	if t.underline {
		err := t.ti.EnterUnderlineMode()
		if err != nil {
			return err
		}
	}

	if t.reverse {
		err := t.ti.EnterReverseMode()
		if err != nil {
			return err
		}
	}

	if t.dim {
		err := t.ti.EnterDimMode()
		if err != nil {
			return err
		}
	}

	if t.bold {
		err := t.ti.EnterBoldMode()
		if err != nil {
			return err
		}
	}

	if t.italic {
		err := t.ti.EnterItalicsMode()
		if err != nil {
			return err
		}
	}

	if t.fg != ColorDefault || t.bg != ColorDefault {
		err := t.ti.Color(t.fg, t.bg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Text) Print(a ...interface{}) error {
	err := t.Do()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(t.ti.w, a...)
	if err != nil {
		return err
	}
	return t.ti.ExitAttributeMode()
}

func (t *Text) Println(a ...interface{}) error {
	err := t.Do()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(t.ti.w, a...)
	if err != nil {
		return err
	}
	return t.ti.ExitAttributeMode()
}

func (t *Text) Printf(format string, a ...interface{}) error {
	err := t.Do()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(t.ti.w, format, a...)
	if err != nil {
		return err
	}
	return t.ti.ExitAttributeMode()
}
