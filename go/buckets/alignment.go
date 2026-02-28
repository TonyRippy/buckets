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

// Alignment describes where an index's value should be on the bucket for that index.
// An index's value is computed using a formula $f(i)$ for index $i$.
// A left-aligned bucketing strategy would use buckets of $[f(i), f(i+1))$ for index `i`.
// A right-aligned bucketing strategy would use buckets of $(f(i-1), f(i)]$ for index `i`.
// This does not apply to all bucketing strategies.
type Alignment uint8

const (
	// Left indicates that the index's value should be on the left side of the bucket for that index.
	// For example, a range that is closed on the left side is [0, 10).
	Left Alignment = iota

	// Right indicates that the range is aligned on the right side.
	// For example, a range that is closed on the right side is (0, 10].
	Right
)

// String returns a string representation of the closed side.
func (s Alignment) String() string {
	switch s {
	case Left:
		return "left"
	case Right:
		return "right"
	default:
		return fmt.Sprintf("unknown(%d)", s)
	}
}

// ParseAlignment parses a string representation of a closed side.
func ParseAlignment(s string) (Alignment, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "left":
		return Left, nil
	case "right":
		return Right, nil
	default:
		return 0, fmt.Errorf("invalid alignment %q", s)
	}
}
