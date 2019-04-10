package terminfo

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type Terminfo struct {
	name string
	path string
	desc string

	// Predefined Capabilities
	bools []byte
	nums  []int
	strs  [][]byte

	// User-Defined Capabilities
	extBools     []byte
	extBoolNames map[string]int

	extNums     []int
	extNumNames map[string]int

	extStrs     [][]byte
	extStrNames map[string]int
}

func (ti *Terminfo) Name() string {
	return ti.name
}

func (ti *Terminfo) Path() string {
	return ti.path
}

func (ti *Terminfo) Desc() string {
	return ti.desc
}

func (ti *Terminfo) GetFlag(i int) bool {
	if i < len(ti.bools) {
		return ti.bools[i] == 1
	}
	return false
}

func (ti *Terminfo) GetNum(i int) (int, bool) {
	if i < len(ti.nums) {
		n := ti.nums[i]
		return n, n >= 0
	}
	return -1, false
}

func (ti *Terminfo) GetStr(i int) ([]byte, bool) {
	if i < len(ti.strs) {
		s := ti.strs[i]
		return s, s != nil
	}
	return nil, false
}

func (ti *Terminfo) GetExtFlag(s string) bool {
	if i, ok := ti.extBoolNames[s]; ok {
		return ti.extBools[i] == 1
	}
	return false
}

func (ti *Terminfo) GetExtNum(s string) (int, bool) {
	if i, ok := ti.extNumNames[s]; ok {
		n := ti.extNums[i]
		return n, n >= 0
	}
	return -1, false
}

func (ti *Terminfo) GetExtStr(s string) ([]byte, bool) {
	if i, ok := ti.extStrNames[s]; ok {
		s := ti.extStrs[i]
		return s, s != nil
	}
	return nil, false
}

func (ti *Terminfo) Sput(cap int, a ...interface{}) ([]byte, error) {
	b, ok := ti.GetStr(cap)
	if !ok {
		return nil, fmt.Errorf("str cap absent: %d", cap)
	}

	b, err := Fmt(b, a...)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return b, err
}

func (ti *Terminfo) Fput(w io.Writer, cap int, a ...interface{}) error {
	b, err := ti.Sput(cap, a...)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func (ti *Terminfo) Put(cap int, a ...interface{}) error {
	return ti.Fput(os.Stdout, cap, a...)
}

func (ti *Terminfo) ClearScreen() error {
	return ti.Put(ClearScreen)
}

func (ti *Terminfo) Cursor(row, col int) error {
	return ti.Put(CursorAddress, row, col)
}

func (ti *Terminfo) Color(fg, bg int) error {
	mc, ok := ti.GetNum(MaxColors)
	if !ok {
		return errors.New("no MaxColors cap")
	}

	if fg != ColorDefault {
		if fg < 0 || fg >= mc {
			return errors.New("illegal fg color")
		}
		err := ti.Put(SetAForeground, fg)
		if err != nil {
			return err
		}
	}

	if bg != ColorDefault {
		if bg < 0 || bg >= mc {
			return errors.New("illegal bg color")
		}
		err := ti.Put(SetABackground, bg)
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
	return ti.Put(OrigPair)
}
