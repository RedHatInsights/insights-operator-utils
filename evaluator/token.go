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
	"go/token"
)

// TokenWithValue structure represents token tied with its value (if any)
type TokenWithValue struct {
	Token      token.Token
	Value      int
	Identifier string
}

// ValueToken is constructor for TokenWithValue structure
func ValueToken(tok token.Token, value int) TokenWithValue {
	return TokenWithValue{
		Token: tok,
		Value: value,
	}
}

// OperatorToken is constructor for TokenWithValue structure
func OperatorToken(tok token.Token) TokenWithValue {
	return TokenWithValue{
		Token: tok,
	}
}

// IdentifierToken is constructor for TokenWithValue structure
func IdentifierToken(tok token.Token, identifier string) TokenWithValue {
	return TokenWithValue{
		Token:      tok,
		Identifier: identifier,
	}
}
