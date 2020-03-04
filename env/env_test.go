/*
Copyright © 2019, 2020 Red Hat, Inc.

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

package env_test

import (
	"os"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/env"
)

const envVariableName = "MY_ENV_VARIABLE"
const envVariableDefaultValue = "fallback"

// TestGetEnvExistingVariable check whether the reading from existing environment variable is correct.
func TestGetEnvExistingVariable(t *testing.T) {
	const envVariableValue = "foobar"

	// set the environment variable
	err := os.Setenv(envVariableName, envVariableValue)
	if err != nil {
		t.Fatal(err)
	}

	// check the environment variable value
	value := env.GetEnv(envVariableName, envVariableDefaultValue)
	if value != envVariableValue {
		t.Fatal("Environment variable has no proper value:", value)
	}
}

// TestGetEnvNoVariable check how the non-existent environment variable is handled.
func TestGetEnvNoVariable(t *testing.T) {
	// make sure no environment variables are set up
	os.Clearenv()

	// check the environment variable value
	value := env.GetEnv(envVariableName, envVariableDefaultValue)
	if value != envVariableDefaultValue {
		t.Fatal("Environment variable has no proper value:", value)
	}
}

// TestGetEnvEmptyVariable check how the existing but empty environment variable is handled.
func TestGetEnvEmptyVariable(t *testing.T) {
	const envVariableValue = ""

	// set the environment variable
	err := os.Setenv(envVariableName, envVariableValue)
	if err != nil {
		t.Fatal(err)
	}

	// check the environment variable value
	value := env.GetEnv(envVariableName, envVariableDefaultValue)
	if value != envVariableValue {
		t.Fatal("Environment variable has no proper value:", value)
	}
}
