package terminfo

import (
	"errors"
)

var (
	errStackEmpty          = errors.New("stack: empty")
	errStackNotEmpty       = errors.New("stack: not empty")
	errStackIllegalPopType = errors.New("stack: illegal pop type")
)

type stack struct {
	e []interface{}
	l int
}

func newStack() *stack {
	return &stack{
		e: make([]interface{}, 2),
		l: 0,
	}
}

func (s *stack) isEmpty() bool {
	return s.l == 0
}

func (s *stack) push(e interface{}) {
	if s.l < len(s.e) {
		s.e[s.l] = e
	} else {
		s.e = append(s.e, e)
	}
	s.l++
}

func (s *stack) pop() (interface{}, error) {
	if s.l <= 0 {
		return nil, errStackEmpty
	}
	s.l--
	return s.e[s.l], nil
}

func (s *stack) popI() (int, error) {
	e, err := s.pop()
	if err != nil {
		return 0, err
	}

	i, ok := e.(int)
	if !ok {
		s.push(e)
		return 0, errStackIllegalPopType
	}
	return i, nil
}

func (s *stack) popII() (int, int, error) {
	i1, err := s.popI()
	if err != nil {
		return 0, 0, err
	}
	i2, err := s.popI()
	if err != nil {
		s.push(i1)
		return 0, 0, err
	}
	return i1, i2, nil
}

func (s *stack) popS() (string, error) {
	e, err := s.pop()
	if err != nil {
		return "", err
	}

	es, ok := e.(string)
	if !ok {
		s.push(e)
		return "", errStackIllegalPopType
	}
	return es, nil
}
