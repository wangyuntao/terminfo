package terminfo

import "testing"

import "bytes"

func TestSprintf(t *testing.T) {
	b, err := Sprintf([]byte("0x1b[%i%p1%d;%p2%dH"), 10, 20)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, []byte("0x1b[11;21H")) {
		t.Fail()
	}
}
