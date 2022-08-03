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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-operator-utils/evaluator"
)

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

// TestEvaluatorArithmetic1 checks the evaluator.Evaluate function for simple
// arithmetic expression
func TestEvaluatorArithmetic1(t *testing.T) {
	var values = make(map[string]int)
	expression := "1+2*3"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 7, result)
}

// TestEvaluatorArithmetic2 checks the evaluator.Evaluate function for simple
// arithmetic expression
func TestEvaluatorArithmetic2(t *testing.T) {
	var values = make(map[string]int)
	expression := "4/2-1+5%2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 2, result)
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

// TestEvaluatorRelational1 checks the evaluator.Evaluate function for simple
// relational expression
func TestEvaluatorRelational1(t *testing.T) {
	var values = make(map[string]int)
	expression := "1 < 2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 1, result)
}

// TestEvaluatorRelational2 checks the evaluator.Evaluate function for simple
// relational expression
func TestEvaluatorRelational2(t *testing.T) {
	var values = make(map[string]int)
	expression := "1 >= 2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 0, result)
}

// TestEvaluatorRelational3 checks the evaluator.Evaluate function for simple
// relational expression
func TestEvaluatorRelational3(t *testing.T) {
	var values = make(map[string]int)
	expression := "1 < 2 && 1 > 2 && 1 <= 2 && 1 >= 2 && 1==2 && 1 != 2"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 0, result)
}

// TestEvaluatorBoolean1 checks the evaluator.Evaluate function for simple
// boolean expression
func TestEvaluatorBoolean1(t *testing.T) {
	var values = make(map[string]int)
	expression := "1 && 0"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 0, result)
}

// TestEvaluatorBoolean2 checks the evaluator.Evaluate function for simple
// boolean expression
func TestEvaluatorBoolean2(t *testing.T) {
	var values = make(map[string]int)
	expression := "1 || 0"

	result, err := evaluator.Evaluate(expression, values)

	assert.Nil(t, err, "unexpected error")
	assert.Equal(t, 1, result)
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

// TestEvaluatorWrongInput1 checks the evaluator.Evaluate function for
// expression that is not correct
func TestEvaluatorWrongInput1(t *testing.T) {
	var values = make(map[string]int)
	expression := "1**"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEvaluatorWrongInput2 checks the evaluator.Evaluate function for
// expression that is not correct
func TestEvaluatorWrongInput2(t *testing.T) {
	var values = make(map[string]int)
	expression := "(1+2*"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEvaluatorWrongInput3 checks the evaluator.Evaluate function for
// expression that is not correct
func TestEvaluatorWrongInput3(t *testing.T) {
	var values = make(map[string]int)
	expression := "+"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}

// TestEvaluatorMissingValue checks the evaluator.Evaluate function for
// expression that use value not provided
func TestEvaluatorMissingValue(t *testing.T) {
	var values = make(map[string]int)
	expression := "value"

	_, err := evaluator.Evaluate(expression, values)

	assert.Error(t, err, "error is expected")
}
