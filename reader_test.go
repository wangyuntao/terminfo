package terminfo

import (
	"bytes"
	"testing"
)

func newTestReader(n int) *reader {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(i + 1)
	}
	return newReader(b)
}

func TestLen(t *testing.T) {
	r := newTestReader(2)
	if r.len() != 2 {
		t.Fail()
	}
}

func TestReadBytes(t *testing.T) {
	r := newTestReader(4)

	b, err := r.readBytes(2)
	if err != nil {
		t.Fail()
	}

	if !bytes.Equal(b, []byte{1, 2}) {
		t.Fail()
	}

	b, err = r.readBytes(2)
	if err != nil {
		t.Fail()
	}

	if !bytes.Equal(b, []byte{3, 4}) {
		t.Fail()
	}

	if r.len() != 0 {
		t.Fail()
	}
}

func TestReadInts(t *testing.T) {
	r := newTestReader(6)

	a, err := r.readInt16s(1)
	if err != nil {
		t.Fail()
	}

	if a[0] != 2<<8|1 {
		t.Fail()
	}

	a, err = r.readInt32s(1)
	if err != nil {
		t.Fail()
	}

	if a[0] != 6<<24|5<<16|4<<8|3 {
		t.Fail()
	}

	if r.len() != 0 {
		t.Fail()
	}
}
