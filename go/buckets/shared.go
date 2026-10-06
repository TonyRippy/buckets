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

import "math"

// bucketRange constructs a range from the same shared boundaries used by
// correctBucketIndex. Ordinary buckets have an open lower boundary and a closed
// upper boundary. The underflow and overflow buckets extend to negative and
// positive infinity respectively, with those infinite endpoints excluded.
// boundaryAt supplies the upper boundary at each index; adjacent buckets reuse
// that value so rounding cannot introduce gaps or overlaps between them.
func bucketRange(index int32, boundaryAt func(int32) float64) Range {
	switch index {
	case UnderflowBucketIndex:
		return Range{From: math.Inf(-1), To: boundaryAt(index), FromBound: Open, ToBound: Closed}
	case OverflowBucketIndex:
		return Range{From: boundaryAt(index - 1), To: math.Inf(1), FromBound: Open, ToBound: Open}
	default:
		return Range{From: boundaryAt(index - 1), To: boundaryAt(index), FromBound: Open, ToBound: Closed}
	}
}

// correctBucketIndex adjusts an estimated index to agree with the boundaries
// used by Range. Inverse calculations such as division or logarithms can round
// differently from the multiplication or exponentiation used to build boundaries,
// assigning a value to a bucket whose range excludes it.
//
// Buckets have open lower and closed upper boundaries. The result is the first
// index whose upper boundary is at least value, or OverflowBucketIndex if none
// qualifies. The function walks up to four steps toward the correct bucket, then
// binary-searches within the bounds established by the walk. This keeps nearby
// corrections cheap while bounding the work when many boundaries round to the
// same value and the intervening buckets are empty.
//
// The caller must supply a finite value and a clamped candidate index. boundaryAt
// must return nondecreasing, non-NaN upper boundaries for indices from
// UnderflowBucketIndex through OverflowBucketIndex-1; boundaries may be infinite.
func correctBucketIndex(value float64, index int32, boundaryAt func(int32) float64) int32 {
	// Try nearby buckets first, retaining bounds for the binary-search fallback.
	const maxWalkSteps = 4
	low, high := int64(UnderflowBucketIndex), int64(OverflowBucketIndex)
	for steps := 0; ; steps++ {
		var direction int32
		switch {
		case index > UnderflowBucketIndex && value <= boundaryAt(index-1):
			high = int64(index) - 1
			direction = -1
		case index < OverflowBucketIndex && value > boundaryAt(index):
			low = int64(index) + 1
			direction = 1
		default:
			return index
		}
		if steps == maxWalkSteps {
			break
		}
		index += direction
	}

	// Find the first closed upper boundary that contains value. Binary search
	// also skips arbitrarily many empty buckets caused by rounded boundaries.
	// The overflow bucket is the fallback if no finite-index boundary qualifies.
	for low < high {
		mid := low + (high-low)/2
		if value <= boundaryAt(int32(mid)) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return int32(low)
}
