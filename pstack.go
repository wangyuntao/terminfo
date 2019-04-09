package terminfo

import (
	"errors"
)

type stack struct {
	items []interface{}
	len   int
}

func newStack() *stack {
	return &stack{
		items: make([]interface{}, 4),
		len:   0,
	}
}

func (s *stack) push(i interface{}) {
	if s.len >= len(s.items) {
		s.items = append(s.items, 0, 0, 0, 0)
	}
	s.items[s.len] = i
	s.len++
}

func (s *stack) pop() (interface{}, error) {
	if s.len <= 0 {
		return nil, errors.New("empty stack")
	}
	s.len--
	return s.items[s.len], nil
}

func (s *stack) popN() (int, error) {
	i, err := s.pop()
	if err != nil {
		return 0, err
	}
	if n, ok := i.(int); ok {
		return n, nil
	}
	return 0, errors.New("not number")
}

func (s *stack) popN2() (int, int, error) {
	n1, err := s.popN()
	if err != nil {
		return 0, 0, err
	}
	n2, err := s.popN()
	if err != nil {
		return 0, 0, err
	}
	return n1, n2, nil
}

func (s *stack) popS() (string, error) {
	i, err := s.pop()
	if err != nil {
		return "", err
	}
	if r, ok := i.(string); ok {
		return r, nil
	}
	return "", errors.New("not string")
}
