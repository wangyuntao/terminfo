package terminfo

import (
	"bytes"
	"testing"
)

func newTestReader(n int) *Reader {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(i + 1)
	}
	return NewReader(b)
}

func TestLen(t *testing.T) {
	r := newTestReader(2)
	if r.Len() != 2 {
		t.Fail()
	}
}

func TestReadBytes(t *testing.T) {
	r := newTestReader(4)

	b, err := r.ReadBytes(2)
	if err != nil {
		t.Fail()
	}

	if !bytes.Equal(b, []byte{1, 2}) {
		t.Fail()
	}

	b, err = r.ReadBytes(2)
	if err != nil {
		t.Fail()
	}

	if !bytes.Equal(b, []byte{3, 4}) {
		t.Fail()
	}

	if r.Len() != 0 {
		t.Fail()
	}
}

func TestReadInts(t *testing.T) {
	r := newTestReader(6)

	a, err := r.ReadInt16s(1)
	if err != nil {
		t.Fail()
	}

	if a[0] != 2<<8|1 {
		t.Fail()
	}

	a, err = r.ReadInt32s(1)
	if err != nil {
		t.Fail()
	}

	if a[0] != 6<<24|5<<16|4<<8|3 {
		t.Fail()
	}

	if r.Len() != 0 {
		t.Fail()
	}
}
