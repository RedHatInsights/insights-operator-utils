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

	"go/scanner"
	"go/token"
)

// Operator type for all functions that implements any dyadic operator
type Operator func(int, int) int

// toint function convert Boolean value to integer
func toint(x bool) int {
	if x {
		return 1
	}

	return 0
}

// tobool function convert integer value to Boolean
func tobool(x int) bool {
	return x != 0
}

// evaluate function evaluates expression represented as token sequence in postfix notation
func evaluateRPN(expr []TokenWithValue, values map[string]int) (Stack, error) {
	// all implemented dyadic (arithmetic) operators
	operators := map[token.Token]Operator{
		token.ADD:  func(x int, y int) int { return x + y },
		token.SUB:  func(x int, y int) int { return x - y },
		token.MUL:  func(x int, y int) int { return x * y },
		token.QUO:  func(x int, y int) int { return x / y },
		token.REM:  func(x int, y int) int { return x % y },
		token.EQL:  func(x int, y int) int { return toint(x == y) },
		token.LSS:  func(x int, y int) int { return toint(x < y) },
		token.GTR:  func(x int, y int) int { return toint(x > y) },
		token.NEQ:  func(x int, y int) int { return toint(x != y) },
		token.LEQ:  func(x int, y int) int { return toint(x <= y) },
		token.GEQ:  func(x int, y int) int { return toint(x >= y) },
		token.LAND: func(x int, y int) int { return toint(tobool(x) && tobool(y)) },
		token.LOR:  func(x int, y int) int { return toint(tobool(x) || tobool(y)) },
	}

	// operand stack is empty at beginning
	var stack Stack

	// token sequence processing
	for _, tok := range expr {
		// does the token represents dyadic operator?
		operator, isOperator := operators[tok.Token]
		if isOperator {
			// we found and operator
			// -> evaluate the dyadic operation (with two operands)
			// -> store result onto the stack
			err := performArithmeticOperation(&stack, operator)
			if err != nil {
				return stack, err
			}
		} else {
			// token does not represent operator:
			// -> is it integer value or identifier then?
			switch tok.Token {
			case token.INT:
				// we found integer value
				// so store it to the operand stack
				stack.Push(tok.Value)
			case token.IDENT:
				// we found identifier name
				// so try to find the value + store the value onto the operand stack
				value, found := values[tok.Identifier]
				if !found {
					return stack, fmt.Errorf("Unknown identifier: %s", tok.Identifier)
				}
				stack.Push(value)
			default:
				// token of different type (shall it happen?)
				return stack, fmt.Errorf("Incorrect input token: %v", tok)
			}
		}
	}

	// now the operand stack should contain just one value -> the result
	return stack, nil
}

// performArithmeticOperation function perform selected arithmetic operator
// against two values taken from stack
func performArithmeticOperation(stack *Stack, operator Operator) error {
	// read the second operand from the stack + check for empty stack
	y, err := stack.Pop()
	if err != nil {
		return err
	}

	// read the first operand from the stack + check for empty stack
	x, err := stack.Pop()
	if err != nil {
		return err
	}

	// perform the selected arithmeric operation
	result := operator(x, y)

	// store result (one value) back onto the stack
	stack.Push(result)

	// no error
	return nil
}

// Evaluate function evaluates given algebraic expression and return its result
func Evaluate(expression string, values map[string]int) (int, error) {
	// scanner object (lexer)
	var s scanner.Scanner

	// structure that represents set of source file(s)
	fset := token.NewFileSet()

	// info about source file
	file := fset.AddFile("", fset.Base(), len(expression))

	// initialize the scanner
	s.Init(file, []byte(expression), nil, scanner.ScanComments)

	// transform input expression into postfix notation
	postfixExpression := toRPN(s)

	// evaluate the expression represented in postfix notation
	stack, err := evaluateRPN(postfixExpression, values)
	if err != nil {
		return -1, err
	}

	if stack.Empty() {
		return -1, fmt.Errorf("empty stack")
	}

	// stack is not empty, so its TOP is the result
	value, _ := stack.Pop()
	return value, nil
}
