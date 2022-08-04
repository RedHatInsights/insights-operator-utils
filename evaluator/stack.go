/*
Copyright © 2022 Pavel Tisnovsky

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package evaluator

import (
	"fmt"
)

// Stack struct is naive but fully functional implementation of operand stack
type Stack struct {
	stack []int
}

// Push method pushes value onto the stack
func (stack *Stack) Push(value int) {
	stack.stack = append(stack.stack, value)
}

// Pop method pops value from stack with check if stack is empty
func (stack *Stack) Pop() (int, error) {
	if stack.Empty() {
		return -1, fmt.Errorf("Empty stack")
	}

	// index of top of stack (TOS)
	tos := len(stack.stack) - 1

	// read element from the stack
	element := stack.stack[tos]

	// remove element from the stack
	stack.stack = stack.stack[:tos]
	return element, nil
}

// Empty method checks if the stack is empty
func (stack *Stack) Empty() bool {
	return len(stack.stack) == 0
}

// Size method returns number of items on stack
func (stack *Stack) Size() int {
	return len(stack.stack)
}
