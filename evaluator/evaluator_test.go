// Copyright 2022 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evaluator_test

import (
	"fmt"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/evaluator"
)

type TestCase struct {
	name          string
	expression    string
	expectedValue int
	expectedError bool
}

// TestEvaluatorEmptyInput function checks the evaluator.Evaluate function for
// empty input
func TestEvaluatorEmptyInput(t *testing.T) {
	var values = make(map[string]int)
	expression := ""

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEvaluatorSingleToken function checks the evaluator.Evaluate function for
// single token input
func TestEvaluatorSingleToken(t *testing.T) {
	var values = make(map[string]int)
	expression := "42"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 42, result)
}

// TestEvaluatorArithmetic checks the evaluator.Evaluate function for simple
// arithmetic expression
func TestEvaluatorArithmetic(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "short expression",
			expression:    "1+2*3",
			expectedValue: 7,
		},
		{
			name:          "long expression",
			expression:    "4/2-1+5%2",
			expectedValue: 2,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorParenthesis checks the evaluator.Evaluate function for simple
// arithmetic expression with parenthesis
func TestEvaluatorParenthesis(t *testing.T) {
	var values = make(map[string]int)
	expression := "(1+2)*3"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 9, result)
}

// TestEvaluatorRelational checks the evaluator.Evaluate function for simple
// relational expression
func TestEvaluatorRelational(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "less than",
			expression:    "1 < 2",
			expectedValue: 1,
		},
		{
			name:          "greater or equal",
			expression:    "1 >= 2",
			expectedValue: 0,
		},
		{
			name:          "long expression",
			expression:    "1 < 2 && 1 > 2 && 1 <= 2 && 1 >= 2 && 1==2 && 1 != 2",
			expectedValue: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorBoolean checks the evaluator.Evaluate function for simple
// boolean expressions
func TestEvaluatorBoolean(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "and",
			expression:    "1 && 0",
			expectedValue: 0,
		},
		{
			name:          "or",
			expression:    "1 || 0",
			expectedValue: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestEvaluatorValues checks the evaluator.Evaluate function for expression
// with named values
func TestEvaluatorValues(t *testing.T) {
	var values = make(map[string]int)
	values["x"] = 1
	values["y"] = 2
	expression := "x+y*2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 5, result)
}

// TestEvaluatorWrongInput checks the evaluator.Evaluate function for
// expression that is not correct
func TestEvaluatorWrongInput(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "mul instead of right operand",
			expression:    "1**",
			expectedError: true,
		},
		{
			name:          "forgot closing parenthesis",
			expression:    "(1+2*",
			expectedError: true,
		},
		{
			name:          "no operands",
			expression:    "+",
			expectedError: true,
		},
		{
			name:          "no right operand",
			expression:    "2+",
			expectedError: true,
		},
		{
			name:          "no left operand",
			expression:    "+2",
			expectedError: true,
		},
		{
			name:          "no left operand (minus)",
			expression:    "-2",
			expectedError: true,
		},
		{
			name:          "== typo",
			expression:    "0=0",
			expectedError: true,
		},
		{
			name:          "zero division",
			expression:    "1/0",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedError {
				result, err := evaluator.Evaluate(tc.expression, values)
				assert.Error(t, err, "error is expected")
				assert.Equal(t, -1, result)
			}
		})
	}
}

// TestEvaluatorMissingValue checks the evaluator.Evaluate function for
// expression that use value not provided
func TestEvaluatorMissingValue(t *testing.T) {
	var values = make(map[string]int)
	expression := "value"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEdgeCases tests expressions that rarely happen
// in the real world
func TestEdgeCases(t *testing.T) {
	var values = make(map[string]int)
	testCases := []TestCase{
		{
			name:          "useless parenthesis",
			expression:    "(2)*(2)",
			expectedValue: 4,
		},
		{
			name:          "multiple useless parenthesis",
			expression:    "(((0))) >= 0",
			expectedValue: 1,
		},
		{
			name:          "scrambled useless parenthesis",
			expression:    "((((0==0)))+1)",
			expectedValue: 2,
		},
		{
			name:          "0 addition idempotence",
			expression:    "1+0+0+0+0",
			expectedValue: 1,
		},
		{
			name:          "1 division idempotence",
			expression:    "5/1/1/1/1/1",
			expectedValue: 5,
		},
		{
			name:          "transitivity",
			expression:    "(3 > 2) && (2 > 1) == (3 > 1)",
			expectedValue: 1,
		},
		{
			name:          "big integer",
			expression:    "9223372036854775807+100-100",
			expectedValue: 9223372036854775807,
		},
		{
			name:          "overflow",
			expression:    "9223372036854775807+1",
			expectedValue: -9223372036854775808,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evaluator.Evaluate(tc.expression, values)
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expectedValue, result)
		})
	}
}

// TestToInt tests the function toint
func TestToInt(t *testing.T) {
	// conversion from false to 0
	result := evaluator.ToInt(false)
	assert.Equal(t, 0, result)

	// conversion from true to 1
	result = evaluator.ToInt(true)
	assert.Equal(t, 1, result)
}

// TestToBool tests the function tobool
func TestToBool(t *testing.T) {
	// conversion 0 to false
	result := evaluator.ToBool(0)
	assert.False(t, result)

	// conversion 1 to true
	result = evaluator.ToBool(1)
	assert.True(t, result)
}

// TestEvaluateRPNNoTokens tests the function evaluateRPN when no tokens are
// provided at input
func TestEvaluateRPNNoTokens(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	stack, err := evaluator.EvaluateRPN(tokens, values)

	// check the output
	assert.NoError(t, err)
	assert.True(t, stack.Empty())
}

// TestEvaluateRPNInvalidToken tests the function evaluateRPN when invalid
// token is provided at input
func TestEvaluateRPNInvalidToken(t *testing.T) {
	// these tokens are not supported
	invalidTokens := []token.Token{
		token.ILLEGAL,
		token.EOF,
		token.COMMENT,
		token.BREAK,
		token.CASE,
		token.CHAN,
		token.CONST,
		token.CONTINUE,
		token.DEFAULT,
		token.DEFER,
		token.ELSE,
		token.FALLTHROUGH,
		token.FOR,
		token.FUNC,
		token.GO,
		token.GOTO,
		token.IF,
		token.IMPORT,
		token.INTERFACE,
		token.MAP,
		token.PACKAGE,
		token.RANGE,
		token.RETURN,
		token.SELECT,
		token.STRUCT,
		token.SWITCH,
		token.TYPE,
		token.VAR,
	}

	// check all invalid tokens
	for _, invalidToken := range invalidTokens {
		name := fmt.Sprintf("EvaluatingInvalidToken %v", invalidToken)
		t.Run(name, func(t *testing.T) {
			// tokens to be tokenized
			tokens := []evaluator.TokenWithValue{
				evaluator.TokenWithValue{invalidToken, -1, ""},
			}

			// value map used during evaluation
			var values = make(map[string]int)

			// evaluate expression represented as sequence of tokens in RPN order
			_, err := evaluator.EvaluateRPN(tokens, values)

			// check the output -> error needs to be detected
			assert.Error(t, err)
		})
	}
}

// TestEvaluateRPNIntValue tests the function evaluateRPN when just one token
// with integer values is provided at input
func TestEvaluateRPNIntValue(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{
		evaluator.ValueToken(token.INT, 42),
	}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	stack, err := evaluator.EvaluateRPN(tokens, values)

	// check the output
	assert.NoError(t, err)
	assert.False(t, stack.Empty())
	assert.Equal(t, stack.Size(), 1)

	value, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, value, 42)
}

// TestEvaluateRPNTwoIntValues tests the function evaluateRPN when two tokens
// with integer values are provided at input
func TestEvaluateRPNTwoIntValues(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{
		evaluator.ValueToken(token.INT, 1),
		evaluator.ValueToken(token.INT, 2),
	}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	stack, err := evaluator.EvaluateRPN(tokens, values)

	// check the output
	assert.NoError(t, err)
	assert.False(t, stack.Empty())
	assert.Equal(t, stack.Size(), 2)

	// values to be poped from stack in reverse order
	value, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, value, 2)

	// values to be poped from stack in reverse order
	value, err = stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, value, 1)
}

// TestEvaluateRPNArithmeticOperation tests the function evaluateRPN when three tokens
// representing arithmetic expression is evaluated
func TestEvaluateRPNArithmeticOperation(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{
		// RPN order (postfix)
		evaluator.ValueToken(token.INT, 1),
		evaluator.ValueToken(token.INT, 2),
		evaluator.OperatorToken(token.ADD),
	}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	stack, err := evaluator.EvaluateRPN(tokens, values)

	// check the output
	assert.NoError(t, err)
	assert.False(t, stack.Empty())
	assert.Equal(t, stack.Size(), 1)

	value, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, value, 3)
}

// TestEvaluateRPNJustArithmeticOperator tests the function evaluateRPN when just
// arithmetic operator is provided
func TestEvaluateRPNJustArithmeticOperator(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{
		evaluator.OperatorToken(token.ADD),
	}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	_, err := evaluator.EvaluateRPN(tokens, values)

	// check the output -> error needs to be detected
	assert.Error(t, err)
}

// TestEvaluateRPNInsuficientOperand the function evaluateRPN when just
// arithmetic operator and one operand are provided
func TestEvaluateRPNInsuficientOperand(t *testing.T) {
	// tokens to be tokenized
	tokens := []evaluator.TokenWithValue{
		evaluator.ValueToken(token.INT, 1),
		evaluator.OperatorToken(token.ADD),
	}

	// value map used during evaluation
	var values = make(map[string]int)

	// evaluate expression represented as sequence of tokens in RPN order
	_, err := evaluator.EvaluateRPN(tokens, values)

	// check the output -> error needs to be detected
	assert.Error(t, err)
}

// TestPerformArithmeticOperation check the behaviour of function
// performArithmeticOperation for correct operators and tokens
func TestPerformArithmeticOperation(t *testing.T) {
	// operand stack (also known as data stack)
	stack := evaluator.Stack{}

	// push two values onto the stack
	stack.Push(1)
	stack.Push(2)

	// any token that is not token.QUO or token.REM
	tok := token.ADD

	// perform the selected arithmetic operation
	addOperation := func(x int, y int) int { return x + y }

	// perform the selected arithmetic operation
	err := evaluator.PerformArithmeticOperation(&stack, addOperation, tok)
	assert.NoError(t, err)

	// check stack
	assert.False(t, stack.Empty())
	assert.Equal(t, stack.Size(), 1)

	// stack should contain one value
	value, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, value, 3)
}

// TestPerformArithmeticOperationMissingOperand check the behaviour of function
// performArithmeticOperation for incorrect number of operands
func TestPerformArithmeticOperationMissingOperand(t *testing.T) {
	// operand stack (also known as data stack)
	stack := evaluator.Stack{}

	// push just one value onto the stack
	stack.Push(1)

	// any token that is not token.QUO or token.REM
	tok := token.ADD

	// perform the selected arithmetic operation
	addOperation := func(x int, y int) int { return x + y }

	// perform the selected arithmetic operation
	err := evaluator.PerformArithmeticOperation(&stack, addOperation, tok)
	assert.Error(t, err)
}

// TestPerformArithmeticOperationMissingBothOperands check the behaviour of
// function performArithmeticOperation for incorrect number of operands
func TestPerformArithmeticOperationMissingBothOperands(t *testing.T) {
	// operand stack (also known as data stack)
	stack := evaluator.Stack{}

	// stack is empty!

	// any token that is not token.QUO or token.REM
	tok := token.ADD

	// perform the selected arithmetic operation
	addOperation := func(x int, y int) int { return x + y }

	// perform the selected arithmetic operation
	err := evaluator.PerformArithmeticOperation(&stack, addOperation, tok)
	assert.Error(t, err)
}

// TestPerformArithmeticOperationDivideByNotZero check the behaviour of function
// performArithmeticOperation for divide by any value different from zero
func TestPerformArithmeticOperationDivideByNotZero(t *testing.T) {
	// operand stack (also known as data stack)
	stack := evaluator.Stack{}

	// push two values onto the stack
	stack.Push(4)
	stack.Push(2)

	// any token that is not token.QUO or token.REM
	tok := token.QUO

	// perform the selected arithmetic operation
	addOperation := func(x int, y int) int { return x + y }

	// perform the selected arithmetic operation
	err := evaluator.PerformArithmeticOperation(&stack, addOperation, tok)
	assert.NoError(t, err)
}

// TestPerformArithmeticOperationDivideByZero check the behaviour of function
// performArithmeticOperation for divide by zero
func TestPerformArithmeticOperationDivideByZero(t *testing.T) {
	// operand stack (also known as data stack)
	stack := evaluator.Stack{}

	// push two values onto the stack
	stack.Push(1)
	stack.Push(0)

	// any token that is not token.QUO or token.REM
	tok := token.QUO

	// perform the selected arithmetic operation
	addOperation := func(x int, y int) int { return x + y }

	// perform the selected arithmetic operation
	err := evaluator.PerformArithmeticOperation(&stack, addOperation, tok)
	assert.Error(t, err)
}
