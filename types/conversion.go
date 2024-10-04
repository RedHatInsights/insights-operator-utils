// Copyright 2024 Red Hat, Inc
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

package types

// Uint64ToUint32 safely converts a uint64 into a uint32 without overflow
func Uint64ToUint32(v uint64) (uint32, error) {
	// Check for overflow without casting
	if v > uint64(^uint32(0)) { // ^uint32(0) is the maximum uint32 value (4,294,967,295)
		return 0, &OutOfRangeError{
			Value: v,
			Type:  "uint32",
		}
	}

	return uint32(v), nil // #nosec G103: safe after bounds check
}
