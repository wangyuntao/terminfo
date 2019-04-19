package terminfo

import (
	"encoding/binary"
	"errors"
	"io"
)

type reader struct {
	b []byte
	i int
}

func newReader(b []byte) *reader {
	return &reader{b, 0}
}

func (r *reader) len() int {
	return len(r.b) - r.i
}

func (r *reader) isEmpty() bool {
	return r.len() == 0
}

func (r *reader) align() {
	if r.i%2 == 1 {
		r.i += 1
	}
}

func (r *reader) sliceUnsafe(off int) []byte {
	return r.b[r.i+off:]
}

func (r *reader) slice(off int) ([]byte, error) {
	if r.len() < off {
		return nil, io.ErrUnexpectedEOF
	}
	return r.sliceUnsafe(off), nil
}

func (r *reader) readBytes(n int) ([]byte, error) {
	if r.len() < n {
		return nil, io.ErrUnexpectedEOF
	}
	r.i += n
	return r.b[r.i-n : r.i], nil
}

func (r *reader) readInts(n, byteSize int) ([]int, error) {
	r.align()

	if r.len() < n*byteSize {
		return nil, io.ErrUnexpectedEOF
	}
	rst := make([]int, n)
	switch byteSize {
	case 2:
		for i := 0; i < n; i++ {
			rst[i] = int(int16(binary.LittleEndian.Uint16(r.sliceUnsafe(i * 2))))
		}
	case 4:
		for i := 0; i < n; i++ {
			rst[i] = int(int32(binary.LittleEndian.Uint32(r.sliceUnsafe(i * 4))))
		}
	default:
		return nil, errors.New("Illegal byteSize")
	}
	r.i += n * byteSize
	return rst, nil
}

func (r *reader) readInt16s(n int) ([]int, error) {
	return r.readInts(n, 2)
}

func (r *reader) readInt32s(n int) ([]int, error) {
	return r.readInts(n, 4)
}
