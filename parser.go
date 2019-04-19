package terminfo

import (
	"bytes"
	"errors"
	"fmt"
	"os"
)

const (
	magic1      = 0432
	magic2      = 01036
	maxFileSize = 32768
)

// https://invisible-island.net/ncurses/man/term.5.html
func Parse(bs []byte, term, filepath string) (*Terminfo, error) {
	if len(bs) > maxFileSize {
		return nil, fmt.Errorf("exceed file size limit %d", len(bs))
	}

	ti := &Terminfo{name: term, path: filepath}
	r := NewReader(bs)

	hdr, err := parseHeader(r)
	if err != nil {
		return nil, err
	}

	numByteSize := 2
	if hdr[0] == magic2 {
		numByteSize = 4
	}

	err = parseName(r, hdr[1], ti)
	if err != nil {
		return nil, err
	}

	err = parseBool(r, hdr[2], ti)
	if err != nil {
		return nil, err
	}

	err = parseNum(r, hdr[3], numByteSize, ti)
	if err != nil {
		return nil, err
	}

	err = parseStr(r, hdr[4], hdr[5], ti)
	if err != nil {
		return nil, err
	}

	err = parseExt(r, ti, numByteSize)
	if err != nil {
		return nil, err
	}

	ti.w = os.Stdout
	return ti, nil
}

func parseHeader(r *Reader) ([]int, error) {
	hdr, err := r.ReadInt16s(6)
	if err != nil {
		return nil, err
	}

	magic := hdr[0]
	if magic != magic1 && magic != magic2 {
		return nil, fmt.Errorf("Illegal magic number: %#o", magic)
	}

	return hdr, nil
}

func parseName(r *Reader, size int, ti *Terminfo) error {
	if size <= 0 {
		return errors.New("Illegal name size")
	}

	bs, err := r.ReadBytes(size)
	if err != nil {
		return err
	}

	li := len(bs) - 1
	if bs[li] != 0 {
		return errors.New("Illegal name section")
	}

	ti.desc = string(bs[:li])
	return nil
}

func parseBool(r *Reader, size int, ti *Terminfo) error {
	if size < 0 {
		return fmt.Errorf("illegal bool size: %d", size)
	}

	bs, err := r.ReadBytes(size)
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if bs[i] != 1 && bs[i] != 0 {
			return errors.New("illegal bool value")
		}
	}
	ti.bools = bs
	return nil
}

func parseNum(r *Reader, size, byteSize int, ti *Terminfo) error {
	if size < 0 {
		return fmt.Errorf("illegal num size: %d", size)
	}

	nums, err := r.ReadInts(size, byteSize)
	if err != nil {
		return err
	}

	for i := 0; i < size; i++ {
		if nums[i] < -1 {
			return fmt.Errorf("Illegal number %d", nums[i])
		}
	}
	ti.nums = nums
	return nil
}

func parseStr(r *Reader, nOffset, stSize int, ti *Terminfo) error {
	if nOffset < 0 {
		return fmt.Errorf("illegal string offset size: %d", nOffset)
	}

	if stSize < 0 {
		return errors.New("Illegal string table size")
	}

	offsets, err := r.ReadInt16s(nOffset)
	if err != nil {
		return err
	}

	st, err := r.ReadBytes(stSize)
	if err != nil {
		return err
	}

	ti.strs = make([][]byte, nOffset)
	for i := 0; i < nOffset; i++ {
		offset := offsets[i]
		if offset < -1 {
			return fmt.Errorf("illegal string offset: %d", offset)
		}
		if offset == -1 {
			continue
		}

		idx := bytes.IndexByte(st[offset:], 0)
		if idx == -1 {
			return errors.New("illegal string val")
		}

		if idx == 0 {
			// return errors.New("str val is empty")
			continue // TODO return error?
		}

		ti.strs[i] = st[offset : offset+idx]
	}
	return nil
}

func parseExt(r *Reader, ti *Terminfo, numByteSize int) error {
	if r.isEmpty() {
		return nil
	}

	hdr, err := r.ReadInt16s(5)
	if err != nil {
		return err
	}

	cbool, cnum, cstr, citem, stSize := hdr[0], hdr[1], hdr[2], hdr[3], hdr[4]
	if cbool+cnum+cstr*2 != citem {
		return errors.New("Illegal ext header")
	}

	err = parseExtBool(r, cbool, ti)
	if err != nil {
		return err
	}

	err = parseExtNum(r, cnum, numByteSize, ti)
	if err != nil {
		return err
	}

	err = parseExtStrAndNames(r, cbool, cnum, cstr, stSize, ti)
	if err != nil {
		return err
	}

	if !r.isEmpty() {
		return errors.New("Illegal ext data size")
	}

	return nil
}

func parseExtBool(r *Reader, size int, ti *Terminfo) error {
	if size < 0 {
		return errors.New("Illegal ext bool size")
	}

	bs, err := r.ReadBytes(size)
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		if bs[i] != 1 && bs[i] != 0 {
			return errors.New("illegal ext bool value")
		}
	}
	ti.extBools = bs
	return nil
}

func parseExtNum(r *Reader, size, byteSize int, ti *Terminfo) error {
	if size < 0 {
		return errors.New("Illegal ext num size")
	}

	nums, err := r.ReadInts(size, byteSize)
	if err != nil {
		return err
	}

	for i := 0; i < size; i++ {
		if nums[i] < -1 {
			return fmt.Errorf("Illegal ext number %d", nums[i])
		}
	}
	ti.extNums = nums
	return nil
}

func parseExtStrAndNames(r *Reader, cbool, cnum, cstr, stSize int, ti *Terminfo) error {
	if cstr < 0 {
		return errors.New("Illegal ext string offset size")
	}

	if stSize < 0 {
		return errors.New("Illegal ext string table size")
	}

	strOffsets, err := r.ReadInt16s(cstr)
	if err != nil {
		return err
	}

	boolNameOffsets, err := r.ReadInt16s(cbool)
	if err != nil {
		return err
	}

	numNameOffsets, err := r.ReadInt16s(cnum)
	if err != nil {
		return err
	}

	strNameOffsets, err := r.ReadInt16s(cstr)
	if err != nil {
		return err
	}

	st, err := r.ReadBytes(stSize)
	if err != nil {
		return err
	}

	idx, err := parseExtStr(strOffsets, st, ti)
	if err != nil {
		return err
	}

	// ext names

	st = st[idx+1:]

	names, err := parseExtNames(boolNameOffsets, st)
	if err != nil {
		return err
	}
	ti.extBoolNames = names

	names, err = parseExtNames(numNameOffsets, st)
	if err != nil {
		return err
	}
	ti.extNumNames = names

	names, err = parseExtNames(strNameOffsets, st)
	if err != nil {
		return err
	}
	ti.extStrNames = names

	return nil
}

func parseExtStr(offsets []int, st []byte, ti *Terminfo) (int, error) {
	ti.extStrs = make([][]byte, len(offsets))
	idx := 0

	for i := 0; i < len(offsets); i++ {
		offset := offsets[i]
		if offset < -1 {
			return 0, fmt.Errorf("illegal ext string offset: %d", offset)
		}
		if offset == -1 {
			continue
		}

		idx = bytes.IndexByte(st[offset:], 0)
		if idx == -1 {
			return 0, errors.New("illegal ext string")
		}

		if idx == 0 {
			return 0, errors.New("ext str val is empty")
		}

		idx += offset
		ti.extStrs[i] = st[offset:idx]
	}
	return idx, nil
}

func parseExtNames(offsets []int, st []byte) (map[string]int, error) {
	names := make(map[string]int)
	for i := 0; i < len(offsets); i++ {
		offset := offsets[i]
		if offset < 0 {
			return nil, fmt.Errorf("illegal ext name offset: %d", offset)
		}

		idx := bytes.IndexByte(st[offset:], 0)
		if idx == -1 {
			return nil, errors.New("illegal ext names")
		}

		if idx == 0 {
			return nil, errors.New("ext name is empty")
		}

		b := st[offset : offset+idx]
		names[string(b)] = i
	}
	return names, nil
}
