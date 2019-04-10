package terminfo

import (
	"bytes"
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	b, err := Fmt([]byte("0x1b[%i%p1%d;%p2%dH"), 10, 20)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b, []byte("0x1b[11;21H")) {
		t.Fail()
	}
}

func TestFmtColor(t *testing.T) {
	eseq := []byte("\\E[%?%p1%{8}%<%t3%p1%d%e%p1%{16}%<%t9%p1%{8}%-%d%e38;5;%p1%d%;m")

	for i := 0; i < 8; i++ {
		b, err := Fmt(eseq, i)
		if err != nil {
			t.Error(err)
		}

		r := fmt.Sprintf("\\E[3%dm", i)
		if !bytes.Equal(b, []byte(r)) {
			t.Fail()
		}
	}

	for i := 8; i < 16; i++ {
		b, err := Fmt(eseq, i)
		if err != nil {
			t.Error(err)
		}

		r := fmt.Sprintf("\\E[9%dm", i-8)
		if !bytes.Equal(b, []byte(r)) {
			t.Fail()
		}
	}

	for i := 16; i < 256; i++ {
		b, err := Fmt(eseq, i)
		if err != nil {
			t.Error(err)
		}

		r := fmt.Sprintf("\\E[38;5;%dm", i)
		if !bytes.Equal(b, []byte(r)) {
			t.Fail()
		}
	}

}
