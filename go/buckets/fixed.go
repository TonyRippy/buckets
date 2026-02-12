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
	"math"
	"strconv"
	"strings"
)

type fixedBucketer struct {
	Width      float64
	Origin     float64
	ClosedSide ClosedSide
}

// FixedBucketer returns a fixed-width bucketer with the given origin and closed side.
func FixedBucketer(width, origin float64, closedSide ClosedSide) (BucketingStrategy, error) {
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %g", width)
	}
	if closedSide != Left && closedSide != Right {
		return nil, fmt.Errorf("invalid closed side %d", closedSide)
	}
	return &fixedBucketer{Width: width, Origin: origin, ClosedSide: closedSide}, nil
}

func (b *fixedBucketer) IndexOf(value float64) (int32, error) {
	if b.Width <= 0 {
		return 0, fmt.Errorf("invalid width %g", b.Width)
	}
	shifted := (value - b.Origin) / b.Width
	switch b.ClosedSide {
	case Left:
		return int32(math.Floor(shifted)), nil
	case Right:
		return int32(math.Ceil(shifted)), nil
	default:
		return 0, fmt.Errorf("invalid closed side %d", b.ClosedSide)
	}
}

func (b *fixedBucketer) Range(index int32) (Range, error) {
	switch b.ClosedSide {
	case Left:
		min := b.Origin + float64(index)*b.Width
		max := min + b.Width
		return Range{From: min, To: max, FromBound: Closed, ToBound: Open}, nil
	case Right:
		max := b.Origin + float64(index)*b.Width
		min := max - b.Width
		return Range{From: min, To: max, FromBound: Open, ToBound: Closed}, nil
	default:
		return Range{}, fmt.Errorf("invalid closed side %d", b.ClosedSide)
	}
}

func (b *fixedBucketer) String() string {
	parts := []string{}
	if b.Width != 1 {
		parts = append(parts, fmt.Sprintf("width=%g", b.Width))
	}
	if b.Origin != 0 {
		parts = append(parts, fmt.Sprintf("origin=%g", b.Origin))
	}
	if b.ClosedSide == Right {
		parts = append(parts, "closed=right")
	}
	if len(parts) == 0 {
		return "fixed"
	}
	return fmt.Sprintf("fixed:%s", strings.Join(parts, ","))
}

func init() {
	RegisterParser("fixed", func(args map[string]string) (BucketingStrategy, error) {
		width := 1.0
		if arg, ok := args["width"]; ok {
			var err error
			width, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid width %q", arg)
			}
		}

		origin := 0.0
		if arg, ok := args["origin"]; ok {
			var err error
			origin, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid origin %q", arg)
			}
		}

		closedSide := Left
		if arg, ok := args["closed"]; ok {
			var err error
			closedSide, err = ParseClosedSide(arg)
			if err != nil {
				return nil, err
			}
		}

		return FixedBucketer(width, origin, closedSide)
	})
}
