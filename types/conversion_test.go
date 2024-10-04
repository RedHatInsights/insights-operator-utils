// Copyright 2024 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package types

import (
	"testing"
)

func TestUint64ToUint32(t *testing.T) {
	tests := []struct {
		name    string
		input   uint64
		want    uint32
		errType error
	}{
		{"Zero", 0, 0, nil},
		{"Max uint32", 4294967295, 4294967295, nil},
		{"Overflow", 4294967296, 0, &OutOfRangeError{}},
		{"Large overflow", 18446744073709551615, 0, &OutOfRangeError{}},
		{"Mid-range value", 2147483648, 2147483648, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uint64ToUint32(tt.input)
			if tt.errType != nil {
				// Check if the error is of the expected type
				if err == nil {
					t.Errorf("Expected error type %T, got no error", tt.errType)
				} else if _, ok := err.(*OutOfRangeError); !ok {
					t.Errorf("Expected error type *OutOfRangeError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Uint64ToUint32() error = %v, want no error", err)
					return
				}
				if got != tt.want {
					t.Errorf("Uint64ToUint32() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
