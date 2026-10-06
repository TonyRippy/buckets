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

type linearBucketer struct {
	M, B float64
}

// LinearBucketer returns a linear bucketer with the given slope and intercept.
func LinearBucketer(m, b float64) (BucketingStrategy, error) {
	if math.IsNaN(m) || math.IsInf(m, 0) || m <= 0 {
		return nil, fmt.Errorf("invalid slope %g", m)
	}
	if math.IsNaN(b) || math.IsInf(b, 0) {
		return nil, fmt.Errorf("invalid intercept %g", b)
	}
	return &linearBucketer{M: m, B: b}, nil
}

func (b *linearBucketer) IndexOf(value float64) (int32, error) {
	if math.IsNaN(value) {
		return 0, fmt.Errorf("invalid value %g", value)
	}
	if value > b.boundaryAt(math.MaxInt32-1) {
		return math.MaxInt32, nil
	}
	if value <= b.boundaryAt(math.MinInt32) {
		return math.MinInt32, nil
	}
	shifted := (value - b.B) / b.M
	return int32(math.Ceil(shifted)), nil
}

// boundaryAt returns the boundary value at the given index.
func (b *linearBucketer) boundaryAt(index int32) float64 {
	return float64(index)*b.M + b.B
}

func (b *linearBucketer) Range(index int32) (Range, error) {
	switch index {
	case math.MinInt32:
		return Range{From: math.Inf(-1), To: b.boundaryAt(index), FromBound: Open, ToBound: Closed}, nil
	case math.MaxInt32:
		return Range{From: b.boundaryAt(index - 1), To: math.Inf(1), FromBound: Open, ToBound: Open}, nil
	default:
		return Range{From: b.boundaryAt(index - 1), To: b.boundaryAt(index), FromBound: Open, ToBound: Closed}, nil
	}
}

func (b *linearBucketer) String() string {
	parts := []string{}
	if b.M != 1 {
		parts = append(parts, fmt.Sprintf("m=%g", b.M))
	}
	if b.B != 0 {
		parts = append(parts, fmt.Sprintf("b=%g", b.B))
	}
	if len(parts) == 0 {
		return "linear"
	}
	return fmt.Sprintf("linear:%s", strings.Join(parts, ","))
}

func parseLinearBucketer(args map[string]string) (BucketingStrategy, error) {
	m := 1.0
	if arg, ok := args["m"]; ok {
		var err error
		m, err = strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid slope %q", arg)
		}
	}
	b := 0.0
	if arg, ok := args["b"]; ok {
		var err error
		b, err = strconv.ParseFloat(arg, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid intercept %q", arg)
		}
	}
	return LinearBucketer(m, b)
}

func init() {
	RegisterParser("linear", parseLinearBucketer)
}
