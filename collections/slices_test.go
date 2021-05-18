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

package collections_test

// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-utils/packages/collections/slices_test.html

import (
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/collections"
)

// TestStringInSliceEmptySlice tests the behaviour of StringInSlice for regular string and empty slice
func TestStringInSliceEmptySlice(t *testing.T) {
	str := "foo"
	slice := []string{}
	if collections.StringInSlice(str, slice) {
		t.Fatal("False expected for empty slice")
	}
}

// TestStringInSliceEmptySliceEmptyString tests the behaviour of StringInSlice for empty string and empty slice
func TestStringInSliceEmptySliceEmptyString(t *testing.T) {
	str := ""
	slice := []string{}
	if collections.StringInSlice(str, slice) {
		t.Fatal("False expected for empty slice")
	}
}

// TestStringInSliceRegularStringNotFound tests the behaviour of StringInSlice for regular string that is not contained in slice
func TestStringInSliceRegularStringNotFound(t *testing.T) {
	str := "foo"
	slice := []string{"aaa", "bbb", "ccc"}
	if collections.StringInSlice(str, slice) {
		t.Fatal("String should not be found in the slice")
	}
}

// TestStringInSliceEmptyStringNotFound tests the behaviour of StringInSlice for empty string that is not contained in slice
func TestStringInSliceEmptyStringNotFound(t *testing.T) {
	str := ""
	slice := []string{"aaa", "bbb", "ccc"}
	if collections.StringInSlice(str, slice) {
		t.Fatal("Empty string should not be found in the slice")
	}
}

// TestStringInSliceRegularStringFound tests the behaviour of StringInSlice for regules strings contained in slice
func TestStringInSliceRegularStringFound(t *testing.T) {
	slice := []string{"foo", "bar", "baz"}

	// try to find the first item
	if !collections.StringInSlice("foo", slice) {
		t.Fatal("String should be found in the slice")
	}

	// try to find middle item
	if !collections.StringInSlice("bar", slice) {
		t.Fatal("String should be found in the slice")
	}

	// try to find the last item
	if !collections.StringInSlice("baz", slice) {
		t.Fatal("String should be found in the slice")
	}
}

// TestStringInSliceEmptyStringFound tests the behaviour of StringInSlice for empty string contained in slice
func TestStringInSliceEmptyStringFound(t *testing.T) {
	slice := []string{"foo", "", "baz"}

	if !collections.StringInSlice("", slice) {
		t.Fatal("Empty string should be found in the slice")
	}
}

// TestStringInSliceUnicodeStringFound tests the behaviour of StringInSlice for Unicode strings
func TestStringInSliceUnicodeStringFound(t *testing.T) {
	slice := []string{"žluťoučká", "привет", "γεια"}

	// try to find the first item
	if !collections.StringInSlice("žluťoučká", slice) {
		t.Fatal("String should be found in the slice")
	}

	// try to find middle item
	if !collections.StringInSlice("привет", slice) {
		t.Fatal("String should be found in the slice")
	}

	// try to find the last item
	if !collections.StringInSlice("γεια", slice) {
		t.Fatal("String should be found in the slice")
	}
}
