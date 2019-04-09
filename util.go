package terminfo

import (
	"errors"
)

func getStr(bs []byte, offset int) ([]byte, error) {
	s, _, err := getStrAndIdx(bs, offset)
	return s, err
}

func getStrAndIdx(bs []byte, offset int) ([]byte, int, error) {
	for i := offset; i < len(bs); i++ {
		if bs[i] == 0 {
			return bs[offset:i], i, nil
		}
	}
	return nil, 0, errors.New("Illegal string value")
}
