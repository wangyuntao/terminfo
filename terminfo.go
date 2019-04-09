package terminfo

import (
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

func (ti *Terminfo) Sprintf(scap int, a ...interface{}) ([]byte, error) {
	s, ok := ti.GetStr(scap)
	if !ok {
		return nil, fmt.Errorf("str cap absent: %d", scap)
	}

	b, err := Sprintf(s, a...)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return b, err
}

func (ti *Terminfo) Fprintf(w io.Writer, scap int, a ...interface{}) (int, error) {
	b, err := ti.Sprintf(scap, a...)
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}

func (ti *Terminfo) Printf(scap int, a ...interface{}) (int, error) {
	return ti.Fprintf(os.Stdout, scap, a...)
}
