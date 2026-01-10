package coll

// Stack is a generic stack data structure that stores elements of type `T`.
// It provides basic stack operations like push, pop, and peek.
// The type `T` must be comparable, allowing the stack to handle any type that supports equality comparison.
type Stack[T comparable] struct {
	items []T
}

// NewStack creates and returns a new, empty stack.
// Returns:
//   - A pointer to an empty `Stack`.
//
// Example usage:
//
//	stack := NewStack[int]() // Creates a new empty stack of integers.
func NewStack[T comparable]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

// Push adds a new element to the top of the stack.
// Parameters:
//   - `element`: The element to be added to the stack.
//
// Example usage:
//
//	stack.Push(5) // Pushes the value 5 onto the stack.
func (stack *Stack[T]) Push(element T) {
	stack.items = append(stack.items, element)
}

// Peek returns the top element of the stack without removing it.
// If the stack is empty, it returns the zero value for type `T`.
// Returns:
//   - The top element of the stack, or the zero value of `T` if the stack is empty.
//
// Example usage:
//
//	top := stack.Peek() // Gets the top element of the stack without removing it.
func (stack *Stack[T]) Peek() T {
	var lastElement T
	if len(stack.items) > 0 {
		lastElement = stack.items[len(stack.items)-1]
	}
	return lastElement
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns the zero value for type `T`.
// Returns:
//   - The top element of the stack, or the zero value of `T` if the stack is empty.
//
// Example usage:
//
//	top := stack.Pop() // Removes and returns the top element of the stack.
func (stack *Stack[T]) Pop() T {
	var lastElement T
	if len(stack.items) > 0 {
		lastElement = stack.items[len(stack.items)-1]
		stack.items = stack.items[:len(stack.items)-1]
	}
	return lastElement
}

// IsEmpty checks whether the stack contains any elements.
// Returns:
//   - `true` if the stack is empty, `false` otherwise.
//
// Example usage:
//
//	isEmpty := stack.IsEmpty() // Returns true if the stack is empty.
func (stack *Stack[T]) IsEmpty() bool {
	return len(stack.items) == 0
}

// Size returns the number of elements currently in the stack.
// Returns:
//   - The number of elements in the stack.
//
// Example usage:
//
//	size := stack.Size() // Returns the size of the stack.
func (stack *Stack[T]) Size() int {
	return len(stack.items)
}

// Clear removes all elements from the stack, leaving it empty.
//
// Example usage:
//
//	stack.Clear() // Clears all elements from the stack.
func (stack *Stack[T]) Clear() {
	stack.items = make([]T, 0)
}
