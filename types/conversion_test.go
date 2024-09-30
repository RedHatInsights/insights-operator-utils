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
		wantErr bool
	}{
		{"Zero", 0, 0, false},
		{"Max uint32", 4294967295, 4294967295, false},
		{"Overflow", 4294967296, 0, true},
		{"Large overflow", 18446744073709551615, 0, true},
		{"Mid-range value", 2147483648, 2147483648, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Uint64ToUint32(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Uint64ToUint32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Uint64ToUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}
