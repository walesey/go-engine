package renderer

type Stack struct {
	top *Element
	bottom *Element
	size int
} 

type Element struct {
	value interface{}
	next *Element
	previous *Element
}

func CreateStack() *Stack{
	return &Stack{ size: 0 }
}

// Return the stack's length
func (s *Stack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *Stack) Push( value interface{} ) {
	s.top = &Element{value: value, next: s.top}
	if s.size == 0 {
		s.bottom = s.top
	}
	if s.top.next != nil {
		s.top.next.previous = s.top
	}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, returns nil
func (s *Stack) Pop() interface{} {
	if s.size > 0 {
		value := s.top.value
		s.top = s.top.next
		s.size--
		return value
	}
	return nil
}

// Get value at index
func (s *Stack) Get( index int ) interface{} {
	if index >= s.size || index < 0 {
		return nil
	}
	result := s.bottom
	for i := 0 ; i < index ; i++ {
		result = result.previous
	}
	return result.value
}
