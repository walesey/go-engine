package renderer

import (
	"github.com/go-gl/mathgl/mgl32"
	"fmt"
)

type MatStack struct {
	top *Element
	bottom *Element
	size int
} 

type Element struct {
	value mgl32.Mat4
	next *Element
	previous *Element
}

func CreateMatStack() *MatStack{
	return &MatStack{ size: 0 }
}

// Return the stack's length
func (s *MatStack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *MatStack) Push(value mgl32.Mat4) {
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
// If the stack is empty, return error
func (s *MatStack) Pop() (mgl32.Mat4, error) {
	if s.size > 0 {
		value := s.top.value
		s.top = s.top.next
		s.size--
		return value, nil
	}
	return mgl32.Mat4{}, fmt.Errorf("MatStack Empty")
}

//multiply every matrix together
func (s *MatStack) MultiplyAll() mgl32.Mat4 {
	if s.size <= 0 {
		return mgl32.Ident4()
	}
	elem := s.bottom
	result := elem.value
	for i:=1 ; i<s.size ; i++ {
		elem = elem.previous
		result = result.Mul4(elem.value)
	}

	return result
}