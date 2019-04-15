package terminfo

import "io"

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

	w io.Writer
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

func (ti *Terminfo) DefaultWriter(w io.Writer) {
	ti.w = w
}
