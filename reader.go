package terminfo

import (
	"encoding/binary"
	"errors"
	"io"
)

type Reader struct {
	b []byte
	i int
}

func NewReader(b []byte) *Reader {
	return &Reader{b, 0}
}

func (r *Reader) Len() int {
	return len(r.b) - r.i
}

func (r *Reader) isEmpty() bool {
	return r.Len() == 0
}

func (r *Reader) Align() {
	if r.i%2 == 1 {
		r.i += 1
	}
}

func (r *Reader) SliceUnsafe(off int) []byte {
	return r.b[r.i+off:]
}

func (r *Reader) Slice(off int) ([]byte, error) {
	if r.Len() < off {
		return nil, io.ErrUnexpectedEOF
	}
	return r.SliceUnsafe(off), nil
}

func (r *Reader) ReadBytes(n int) ([]byte, error) {
	if r.Len() < n {
		return nil, io.ErrUnexpectedEOF
	}
	r.i += n
	return r.b[r.i-n : r.i], nil
}

func (r *Reader) ReadInts(n, byteSize int) ([]int, error) {
	r.Align()

	if r.Len() < n*byteSize {
		return nil, io.ErrUnexpectedEOF
	}
	rst := make([]int, n)
	switch byteSize {
	case 2:
		for i := 0; i < n; i++ {
			rst[i] = int(int16(binary.LittleEndian.Uint16(r.SliceUnsafe(i * 2))))
		}
	case 4:
		for i := 0; i < n; i++ {
			rst[i] = int(int32(binary.LittleEndian.Uint32(r.SliceUnsafe(i * 4))))
		}
	default:
		return nil, errors.New("Illegal byteSize")
	}
	r.i += n * byteSize
	return rst, nil
}

func (r *Reader) ReadInt16s(n int) ([]int, error) {
	return r.ReadInts(n, 2)
}

func (r *Reader) ReadInt32s(n int) ([]int, error) {
	return r.ReadInts(n, 4)
}
