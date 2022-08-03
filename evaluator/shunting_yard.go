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
	"strconv"

	"go/scanner"
	"go/token"
)

// toRPN function transforms sequence of tokens with expression into PRN code
func toRPN(s scanner.Scanner) []TokenWithValue {
	const stringFmt = "%s "

	// operators with precedence
	var operators = map[token.Token]int{
		// arithmetic operators
		token.MUL: 5,
		token.QUO: 5,
		token.REM: 5,
		token.ADD: 4,
		token.SUB: 4,

		// relational operators
		token.EQL: 3,
		token.LSS: 3,
		token.GTR: 3,
		token.NEQ: 3,
		token.LEQ: 3,
		token.GEQ: 3,

		// logic operators
		token.LAND: 2,
		token.LOR:  1,
	}

	var stack []token.Token

	var output []TokenWithValue

	// tokenization implementation and token processing
loop:
	for {
		_, tok, value := s.Scan()

		switch tok {
		case token.INT:
			// integer value can be added directly into output
			intValue, _ := strconv.Atoi(value)
			output = append(output, ValueToken(tok, intValue))
			fmt.Printf("%d ", intValue)
		case token.IDENT:
			// identifier can be added directly into output
			output = append(output, IdentifierToken(tok, value))
			fmt.Printf(stringFmt, value)
		case token.LPAREN:
			// left paren is pushed into stack
			stack = append(stack, tok)
		case token.RPAREN:
			// right paren start processing operands on stack (until first left paren is found)
			var tok token.Token
			for {
				// read value from stack (POP)
				tok, stack = stack[len(stack)-1], stack[:len(stack)-1]
				if tok == token.LPAREN {
					// remove left paren if found + stop operand processing
					break
				}
				// other tokens poped from stack can be added to output
				output = append(output, OperatorToken(tok))
				fmt.Printf("%v ", tok)
			}
		case token.EOF:
			// special token marking end of tokenization
			break loop
		default:
			priority1, isOperator := operators[tok]
			if isOperator {
				// traverse through values on stack
				for len(stack) > 0 {
					// TOP operation
					tok := stack[len(stack)-1]

					// read priority for operator read from stack
					priority2 := operators[tok]

					// compare operator priorities
					if priority1 > priority2 {
						// priority of read operator is greater than:
						// -> end of processing
						break
					}

					// priority of read operator is less than or equal:
					// -> process read operator and POP it from stack
					stack = stack[:len(stack)-1] // POP
					output = append(output, OperatorToken(tok))
					fmt.Printf(stringFmt, tok)
				}

				// newly read operator needs to be pushed onto stack
				stack = append(stack, tok)
			}
		}
	}
	// clean out the stack at end of processing
	for len(stack) > 0 {
		fmt.Printf(stringFmt, stack[len(stack)-1])
		output = append(output, OperatorToken(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	fmt.Println()

	return output
}
