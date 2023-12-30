package stack

type Stack struct {
	stack []int
}

func (s *Stack) Push(v int) {
	s.stack = append(s.stack, v)
}

func (s *Stack) Pop() (top int) {
	top = s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]
	return top
}

func (s *Stack) Len() int {
	return len(s.stack)
}
