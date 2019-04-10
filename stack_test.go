package terminfo

import "testing"

func TestStack(t *testing.T) {
	s := newStack()
	if !s.isEmpty() {
		t.Fail()
	}

	_, err := s.pop()
	if err != errStackEmpty {
		t.Fail()
	}

	s.push("e")

	_, err = s.popI()
	if err != errStackIllegalPopType {
		t.Fail()
	}

	e, err := s.popS()
	if err != nil {
		t.Fail()
	}

	if e != "e" {
		t.Fail()
	}
}
