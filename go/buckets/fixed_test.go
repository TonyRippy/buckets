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

func TestFixedBucketerIndex(t *testing.T) {
	for _, test := range []struct {
		name   string
		width  float64
		origin float64
		closed ClosedSide
		value  float64
		index  int32
	}{
		{
			name:   "left-closed/default-origin",
			width:  10,
			origin: 0,
			closed: Left,
			value:  9.999,
			index:  0,
		},
		{
			name:   "left-closed/upper-boundary",
			width:  10,
			origin: 0,
			closed: Left,
			value:  10,
			index:  1,
		},
		{
			name:   "left-closed/negative",
			width:  10,
			origin: 0,
			closed: Left,
			value:  -11,
			index:  -2,
		},
		{
			name:   "right-closed/default-origin",
			width:  10,
			origin: 0,
			closed: Right,
			value:  0.001,
			index:  1,
		},
		{
			name:   "right-closed/lower-boundary",
			width:  10,
			origin: 0,
			closed: Right,
			value:  0,
			index:  0,
		},
		{
			name:   "right-closed/negative",
			width:  10,
			origin: 0,
			closed: Right,
			value:  -10,
			index:  -1,
		},
		{
			name:   "left-closed/shifted-origin",
			width:  10,
			origin: 5,
			closed: Left,
			value:  14.9,
			index:  0,
		},
		{
			name:   "left-closed/shifted-origin-boundary",
			width:  10,
			origin: 5,
			closed: Left,
			value:  15,
			index:  1,
		},
		{
			name:   "right-closed/shifted-origin-boundary",
			width:  10,
			origin: 5,
			closed: Right,
			value:  5,
			index:  0,
		},
		{
			name:   "right-closed/shifted-origin-open-lower-boundary",
			width:  10,
			origin: 5,
			closed: Right,
			value:  -5.001,
			index:  -1,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			bucketer, err := FixedBucketer(test.width, test.origin, test.closed)
			if err != nil {
				t.Fatalf("FixedBucketer: %v", err)
			}

			index, err := bucketer.IndexOf(test.value)
			if err != nil {
				t.Fatalf("IndexOf: %v", err)
			}
			if index != test.index {
				t.Fatalf("expected index %d, got %d", test.index, index)
			}

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

func TestFixedBucketerRange(t *testing.T) {
	for _, test := range []struct {
		name   string
		width  float64
		origin float64
		closed ClosedSide
		index  int32
		want   Range
	}{
		{
			name:   "left-closed/default-origin",
			width:  10,
			origin: 0,
			closed: Left,
			index:  0,
			want:   Range{From: 0, To: 10, FromBound: Closed, ToBound: Open},
		},
		{
			name:   "left-closed/negative-index",
			width:  10,
			origin: 0,
			closed: Left,
			index:  -1,
			want:   Range{From: -10, To: 0, FromBound: Closed, ToBound: Open},
		},
		{
			name:   "right-closed/default-origin",
			width:  10,
			origin: 0,
			closed: Right,
			index:  0,
			want:   Range{From: -10, To: 0, FromBound: Open, ToBound: Closed},
		},
		{
			name:   "right-closed/positive-index",
			width:  10,
			origin: 0,
			closed: Right,
			index:  1,
			want:   Range{From: 0, To: 10, FromBound: Open, ToBound: Closed},
		},
		{
			name:   "left-closed/shifted-origin",
			width:  10,
			origin: 5,
			closed: Left,
			index:  0,
			want:   Range{From: 5, To: 15, FromBound: Closed, ToBound: Open},
		},
		{
			name:   "right-closed/shifted-origin",
			width:  10,
			origin: 5,
			closed: Right,
			index:  0,
			want:   Range{From: -5, To: 5, FromBound: Open, ToBound: Closed},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			bucketer, err := FixedBucketer(test.width, test.origin, test.closed)
			if err != nil {
				t.Fatalf("FixedBucketer: %v", err)
			}
			got, err := bucketer.Range(test.index)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			assertRangeEquals(t, test.want, got)
		})
	}
}

func assertFixedBucketerEquals(t *testing.T, want *fixedBucketer, got BucketingStrategy) {
	t.Helper()
	fixed, ok := got.(*fixedBucketer)
	if !ok {
		t.Fatalf("expected fixedBucketer, got %T", got)
	}
	if fixed.Width != want.Width {
		t.Errorf("expected Width %v, got %v", want.Width, fixed.Width)
	}
	if fixed.Origin != want.Origin {
		t.Errorf("expected Origin %v, got %v", want.Origin, fixed.Origin)
	}
	if fixed.ClosedSide != want.ClosedSide {
		t.Errorf("expected ClosedSide %v, got %v", want.ClosedSide, fixed.ClosedSide)
	}
}

func TestFixedBucketerParse(t *testing.T) {
	for _, test := range []struct {
		spec   string
		want   *fixedBucketer
		canon  string
		hasErr bool
	}{
		{
			spec:  "fixed",
			want:  &fixedBucketer{Width: 1, Origin: 0, ClosedSide: Left},
			canon: "fixed",
		},
		{
			spec:  "fixed:width=0.5",
			want:  &fixedBucketer{Width: 0.5, Origin: 0, ClosedSide: Left},
			canon: "fixed:width=0.5",
		},
		{
			spec:  "fixed:width=10,origin=5",
			want:  &fixedBucketer{Width: 10, Origin: 5, ClosedSide: Left},
			canon: "fixed:width=10,origin=5",
		},
		{
			spec:  "fixed:closed=right",
			want:  &fixedBucketer{Width: 1, Origin: 0, ClosedSide: Right},
			canon: "fixed:closed=right",
		},
		{
			spec:  "fixed:width=10,origin=5,closed=right",
			want:  &fixedBucketer{Width: 10, Origin: 5, ClosedSide: Right},
			canon: "fixed:width=10,origin=5,closed=right",
		},
		{
			spec:  " fixed : WIDTH=10, ORIGIN=5, CLOSED=RIGHT ",
			want:  &fixedBucketer{Width: 10, Origin: 5, ClosedSide: Right},
			canon: "fixed:width=10,origin=5,closed=right",
		},
		{
			spec:   "fixed:width=0",
			hasErr: true,
		},
		{
			spec:   "fixed:width=oops",
			hasErr: true,
		},
		{
			spec:   "fixed:origin=oops",
			hasErr: true,
		},
		{
			spec:   "fixed:closed=oops",
			hasErr: true,
		},
	} {
		t.Run(test.spec, func(t *testing.T) {
			got, err := Parse(test.spec)
			if test.hasErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assertFixedBucketerEquals(t, test.want, got)

			if got.String() != test.canon {
				t.Errorf("expected canonical string %q, got %q", test.canon, got.String())
			}
		})
	}
}

func TestFixedBucketerInvalidClosedSide(t *testing.T) {
	_, err := FixedBucketer(1, 0, ClosedSide(255))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestClosedSideString(t *testing.T) {
	for _, test := range []struct {
		side ClosedSide
		want string
	}{
		{Left, "left"},
		{Right, "right"},
		{ClosedSide(42), "unknown(42)"},
	} {
		t.Run(fmt.Sprintf("side=%d", test.side), func(t *testing.T) {
			if got := test.side.String(); got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}
