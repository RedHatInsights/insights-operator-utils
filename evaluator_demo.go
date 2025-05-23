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

package main

import (
	"fmt"

	"github.com/RedHatInsights/insights-operator-utils/evaluator"
)

// expression that we need to evaluate
const source = `
1 + 2*3 > (severity)
`

func main() {
	// values that can be used in expression
	values := make(map[string]int)
	values["confidence"] = 2
	values["severity"] = 3
	values["priority"] = 1

	value, err := evaluator.Evaluate(source, values)
	fmt.Println(value, err)
}
