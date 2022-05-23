package brainfuck

import "github.com/pkg/errors"

type SliceBasedStack struct {
	data []CommandStackItem
}

func (s SliceBasedStack) Len() int {
	return len(s.data)
}

func (s SliceBasedStack) Top() (CommandStackItem, bool) {
	if len(s.data) > 0 {
		return s.data[len(s.data)-1], true
	}
	return nil, false
}

func (s *SliceBasedStack) Push(item CommandStackItem) error {
	s.data = append(s.data, item)
	return nil
}

func (s *SliceBasedStack) Pop() (CommandStackItem, error) {
	if item, ok := s.Top(); ok {
		s.data = s.data[0 : len(s.data)-1]
		return item, nil
	}
	return nil, errors.New("cannot pop from empty stack")
}
