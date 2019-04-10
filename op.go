package terminfo

import (
	"io"
	"os"
)

func (ti *Terminfo) Do(cap int, a ...interface{}) error {
	return ti.FmtTo(os.Stdout, cap, a...)
}

func (ti *Terminfo) FmtTo(w io.Writer, cap int, a ...interface{}) error {
	b, err := ti.Fmt(cap, a...)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}