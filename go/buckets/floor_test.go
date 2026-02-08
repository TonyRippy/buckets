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

func TestFloorBucketerIndex(t *testing.T) {
	for _, test := range []struct {
		width float64
		value float64
		index int32
	}{
		{10, 0, 0},
		{10, 0.5, 0},
		{10, 9.999, 0},
		{10, 10, 1},
		{10, 11, 1},
		{10, -10, -1},
		{10, -11, -2},
		{10, 111, 11},
	} {
		t.Run(fmt.Sprintf("width=%g,value=%g", test.width, test.value), func(t *testing.T) {
			bucketer, err := FloorBucketer(test.width)
			if err != nil {
				t.Fatalf("FloorBucketer: %v", err)
			}

			// Test that the value is in the expected index
			index, err := bucketer.IndexOf(test.value)
			if err != nil {
				t.Fatalf("IndexOf: %v", err)
			}
			if index != test.index {
				t.Fatalf("expected index %d, got %d", test.index, index)
			}

			// Test that the value is the range for the expected index
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

func TestFloorBucketerRange(t *testing.T) {
	for _, test := range []struct {
		width float64
		index int32
		want  Range
	}{
		{10, 0, Range{From: 0, To: 10, FromBound: Closed, ToBound: Open}},
		{10, 1, Range{From: 10, To: 20, FromBound: Closed, ToBound: Open}},
		{10, -1, Range{From: -10, To: 0, FromBound: Closed, ToBound: Open}},
	} {
		t.Run(fmt.Sprintf("width=%g,index=%d", test.width, test.index), func(t *testing.T) {
			bucketer, err := FloorBucketer(test.width)
			if err != nil {
				t.Fatalf("FloorBucketer: %v", err)
			}
			got, err := bucketer.Range(test.index)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			assertRangeEquals(t, test.want, got)
		})
	}
}

func assertFloorBucketerEquals(t *testing.T, want *floorBucketer, got BucketingStrategy) {
	t.Helper()
	floor, ok := got.(*floorBucketer)
	if !ok {
		t.Fatalf("expected floorBucketer, got %T", got)
	}
	if floor.Width != want.Width {
		t.Errorf("expected Width %v, got %v", want.Width, floor.Width)
	}
}

func TestFloorBucketerParse(t *testing.T) {
	for _, test := range []struct {
		spec string
		want *floorBucketer
		err  bool
	}{
		{"floor", &floorBucketer{Width: 1}, false},
		{"floor:width=0.5", &floorBucketer{Width: 0.5}, false},
		{"floor:width=10", &floorBucketer{Width: 10}, false},
		{"floor:width=0", nil, true},
		{"floor:width=oops", nil, true},
	} {
		t.Run(test.spec, func(t *testing.T) {
			// Verify that the parsing works as expected
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

			// Check that the bucketer is correct
			assertFloorBucketerEquals(t, test.want, got)

			// Check that the string representation is correct
			if got != nil {
				if got.String() != test.spec {
					t.Errorf("expected %q, got %q", test.spec, got.String())
				}
			}
		})
	}
}
