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
	"errors"
	"fmt"
	"math"
)

const (
	// UnderflowBucketIndex is the index of the underflow bucket.
	// It is reserved for all values less than the interval for bucket UnderflowBucketIndex + 1.
	UnderflowBucketIndex int32 = math.MinInt32

	// OverflowBucketIndex is the index of the overflow bucket.
	// It is reserved for all values greater than the interval for bucket OverflowBucketIndex - 1.
	OverflowBucketIndex  int32 = math.MaxInt32
)

var (
	ErrOutOfRange = errors.New("value is out of range")
)

// BucketingStrategy is a strategy for bucketing values into ranges.
type BucketingStrategy interface {
	fmt.Stringer

	// IndexOf returns the index of the bucket that contains the given value.
	IndexOf(value float64) (int32, error)

	// Range returns the range of values that are in the bucket with the given index.
	Range(index int32) (Range, error)
}
