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
)

type floorBucketer struct {
	Width       float64
}

// FloorBucketer returns a bucketer that rounds down to the nearest multiple of width.
func FloorBucketer(width float64) (BucketingStrategy, error) {
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %g", width)
	}
	return &floorBucketer{Width: width}, nil
}

func (b *floorBucketer) IndexOf(value float64) (int32, error) {
	if b.Width <= 0 {
		return 0, fmt.Errorf("invalid width %g", b.Width)
	}
	index := int32(math.Floor(value / b.Width))
	return index, nil
}

func (b *floorBucketer) Range(index int32) (Range, error) {
	min := float64(index)*b.Width
	max := min + b.Width
	return Range{From: min, To: max, FromBound: Closed, ToBound: Open}, nil
}

func (b *floorBucketer) String() string {
	if b.Width == 1 {
		return "floor"
	}
	return fmt.Sprintf("floor:width=%g", b.Width)
}

func init() {
	RegisterParser("floor", func(args map[string]string) (BucketingStrategy, error) {
		width := 1.0
		if arg, ok := args["width"]; ok {
			var err error
			width, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid width %q", arg)
			}
		}
		return FloorBucketer(width)
	})
}
