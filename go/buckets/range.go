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

func (b BoundType) String() string {
	switch b {
	case Open:
		return "open"
	case Closed:
		return "closed"
	default:
		return fmt.Sprintf("unknown(%d)", b)
	}
}

func parseBoundType(value string) (BoundType, error) {
	switch strings.ToLower(value) {
	case "open":
		return Open, nil
	case "closed":
		return Closed, nil
	default:
		return Open, fmt.Errorf(`expected "open" or "closed", got %q`, value)
	}
}

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
