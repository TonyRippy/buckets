// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buckets

import (
	"fmt"
	"testing"
)

func TestCeilBucketerIndex(t *testing.T) {
	for _, test := range []struct {
		width float64
		value float64
		index int32
	}{
		{10, 0, 0},
		{10, 0.001, 1},
		{10, 0.1, 1},
		{10, 9.9, 1},
		{10, 10, 1},
		{10, 11, 2},
		{10, -0.1, 0},
		{10, -9.9, 0},
		{10, -10, -1},
		{10, 111.1, 12},
		{10, -111.1, -11},
	} {
		t.Run(fmt.Sprintf("width=%g,value=%g", test.width, test.value), func(t *testing.T) {
			bucketer, err := CeilBucketer(test.width)
			if err != nil {
				t.Fatalf("CeilBucketer: %v", err)
			}

			// Test that the value is in the expected index.
			index, err := bucketer.IndexOf(test.value)
			if err != nil {
				t.Fatalf("IndexOf: %v", err)
			}
			if index != test.index {
				t.Fatalf("expected index %d, got %d", test.index, index)
			}

			// Test that the value is in the range for the expected index.
			r, err := bucketer.Range(test.index)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			if !r.Contains(test.value) {
				t.Fatalf("expected %v to be in range %v", test.value, r)
			}
		})
	}
}

func TestCeilBucketerRange(t *testing.T) {
	for _, test := range []struct {
		width float64
		index int32
		want  Range
	}{
		{10, 0, Range{From: -10, To: 0, FromBound: Open, ToBound: Closed}},
		{10, 1, Range{From: 0, To: 10, FromBound: Open, ToBound: Closed}},
		{10, -1, Range{From: -20, To: -10, FromBound: Open, ToBound: Closed}},
	} {
		t.Run(fmt.Sprintf("width=%g,index=%d", test.width, test.index), func(t *testing.T) {
			bucketer, err := CeilBucketer(test.width)
			if err != nil {
				t.Fatalf("CeilBucketer: %v", err)
			}
			got, err := bucketer.Range(test.index)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			assertRangeEquals(t, test.want, got)
		})
	}
}

func assertCeilBucketerEquals(t *testing.T, want *ceilBucketer, got BucketingStrategy) {
	t.Helper()
	ceil, ok := got.(*ceilBucketer)
	if !ok {
		t.Fatalf("expected ceilBucketer, got %T", got)
	}
	if ceil.Width != want.Width {
		t.Errorf("expected Width %v, got %v", want.Width, ceil.Width)
	}
}

func TestCeilBucketerParse(t *testing.T) {
	for _, test := range []struct {
		spec string
		want *ceilBucketer
		err  bool
	}{
		{"ceil", &ceilBucketer{Width: 1}, false},
		{"ceil:width=0.5", &ceilBucketer{Width: 0.5}, false},
		{"ceil:width=10", &ceilBucketer{Width: 10}, false},
		{"ceil:width=0", nil, true},
		{"ceil:width=oops", nil, true},
	} {
		t.Run(test.spec, func(t *testing.T) {
			// Verify that parsing works as expected.
			got, err := Parse(test.spec)
			if test.err {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check that the bucketer is correct.
			assertCeilBucketerEquals(t, test.want, got)

			// Check that the string representation is correct.
			if got != nil && got.String() != test.spec {
				t.Errorf("expected %q, got %q", test.spec, got.String())
			}
		})
	}
}
