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
	"strings"
)

type BoundType uint8

const (
	Open BoundType = iota
	Closed
)

type Range struct {
	From      float64
	To        float64
	FromBound BoundType
	ToBound   BoundType
}

func (r Range) Contains(x float64) bool {
	switch r.FromBound {
	case Open:
		if x <= r.From {
			return false
		}
	case Closed:
		if x < r.From {
			return false
		}
	}
	switch r.ToBound {
	case Open:
		if x >= r.To {
			return false
		}
	case Closed:
		if x > r.To {
			return false
		}
	}
	return true
}

func (r Range) String() string {
	var s strings.Builder
	if r.FromBound == Closed {
		s.WriteRune('[')
	} else {
		s.WriteRune('(')
	}
	s.WriteString(fmt.Sprintf("%g, %g", r.From, r.To))
	if r.ToBound == Closed {
		s.WriteRune(']')
	} else {
		s.WriteRune(')')
	}
	return s.String()
}

// BucketingStrategy is a strategy for bucketing values into ranges.
type BucketingStrategy interface {
	fmt.Stringer

	// IndexOf returns the index of the bucket that contains the given value.
	IndexOf(value float64) (int32, error)

	// Range returns the range of values that are in the bucket with the given index.
	Range(index int32) (Range, error)
}

// ClosedSide represents the side of a partially closed range that is closed.
// It is used to indicate whether a bucketer is closed on the left or right side.
type ClosedSide uint8

const (
	// Left indicates that the range is closed on the left side.
	// For example, a range that is closed on the left side is [0, 10).
	Left ClosedSide = iota

	// Right indicates that the range is closed on the right side.
	// For example, a range that is closed on the right side is (0, 10].
	Right
)

// String returns a string representation of the closed side.
func (s ClosedSide) String() string {
	switch s {
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// ParseClosedSide parses a string representation of a closed side.
func ParseClosedSide(s string) (ClosedSide, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	default:
		return 0, fmt.Errorf("invalid closed side %q", s)
	}
}
