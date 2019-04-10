package terminfo

import (
	"errors"
)

var (
	ErrCapAbsent    = errors.New("capability absent")
	ErrCapEscapeSeq = errors.New("illegal escape sequence")
)
