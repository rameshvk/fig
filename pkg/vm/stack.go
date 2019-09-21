package vm

// Stack manages the data stack of arbitrary values
type Stack struct {
	values []interface{}
}

// Push pushes a value onto the stack
func (s *Stack) Push(v interface{}) {
	s.values = append(s.values, v)
}

// Pop pops the value out of the stack
func (s *Stack) Pop() interface{} {
	top := s.values[len(s.values)-1]
	s.values = s.values[:len(s.values)-1]
	return top
}
