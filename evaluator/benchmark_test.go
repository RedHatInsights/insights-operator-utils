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

	"github.com/RedHatInsights/insights-operator-utils/evaluator"
)

func BenchmarkSingleToken(b *testing.B) {
	var values = make(map[string]int)
	expression := "42"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkConstantAddExpression(b *testing.B) {
	var values = make(map[string]int)
	expression := "1+2"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkConstantAddMulExpression(b *testing.B) {
	var values = make(map[string]int)
	expression := "1+2*3"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkConstantLongLogicalExpression(b *testing.B) {
	var values = make(map[string]int)
	expression := "1 < 2 && 1 > 2 && 1 <= 2 && 1 >= 2 && 1==2 && 1 != 2"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkConstantExpressionWithParenthesis(b *testing.B) {
	var values = make(map[string]int)
	expression := "(1+2)*3"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkOneVariableInExpression(b *testing.B) {
	var values = make(map[string]int)
	values["a"] = 1

	expression := "a + 1"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkTwoVariablesInExpression(b *testing.B) {
	var values = make(map[string]int)
	values["a"] = 1
	values["b"] = 2

	expression := "a + b"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}

func BenchmarkThreeVariablesInExpression(b *testing.B) {
	var values = make(map[string]int)
	values["a"] = 1
	values["b"] = 2
	values["c"] = 3

	expression := "a + b * c"

	for i := 0; i < b.N; i++ {
		_, _ = evaluator.Evaluate(expression, values)
	}
}
