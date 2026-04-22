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

const (
	// DefaultZeroWidth is the default width of the zero bucket.
	// The value is borrowed from the Prometheus codebase, which says:
	// "2^-128 (or 0.5*2^-127 in the actual IEEE 754 representation), which is a bucket boundary at all possible resolutions."
	DefaultZeroWidth = 2.938735877055719e-39
)

type exponentialBucketer struct {
	Base   float64
	Origin float64
}

// ExponentialBucketer returns an exponential bucketing strategy.
func ExponentialBucketer(base, origin float64) (BucketingStrategy, error) {
	if math.IsNaN(base) || math.IsInf(base, 0) || base <= 1 {
		return nil, fmt.Errorf("invalid base %g", base)
	}
	if math.IsNaN(origin) || math.IsInf(origin, 0) {
		return nil, fmt.Errorf("invalid origin %g", origin)
	}
	return &exponentialBucketer{Base: base, Origin: origin}, nil
}

func (b *exponentialBucketer) IndexOf(value float64) (int32, error) {
	if math.IsNaN(value) {
		return 0, fmt.Errorf("invalid value %g", value)
	}
	if math.IsInf(value, 1) {
		return math.MaxInt32, nil
	}
	shifted := value - b.Origin
	if value < 0 {
		return math.MinInt32, ErrOutOfRange
	}
	bucket := math.Ceil(math.Log(shifted) / math.Log(b.Base))
	if bucket > float64(math.MaxInt32) {
		return math.MaxInt32, nil
	}
	if bucket < float64(math.MinInt32) {
		return math.MinInt32, nil
	}
	return int32(bucket), nil
}

func (b *exponentialBucketer) Range(index int32) (Range, error) {
	to := b.Origin + math.Pow(b.Base, float64(index))
	if index == UnderflowBucketIndex {
		return Range{From: math.Inf(-1), To: to, FromBound: Open, ToBound: Closed}, nil
	}
	if index == OverflowBucketIndex {
		return Range{From: to, To: math.Inf(1), FromBound: Open, ToBound: Open}, nil
	}
	from := b.Origin + math.Pow(b.Base, float64(index-1))
	return Range{From: from, To: to, FromBound: Open, ToBound: Closed}, nil
}

func (b *exponentialBucketer) String() string {
	var parts []string
	if b.Base != 2 {
		if b.Base == math.E {
			parts = append(parts, "base=e")
		} else {
			parts = append(parts, fmt.Sprintf("base=%g", b.Base))
		}
	}
	if b.Origin != 0 {
		parts = append(parts, fmt.Sprintf("origin=%g", b.Origin))
	}
	if len(parts) == 0 {
		return "exponential"
	}
	return fmt.Sprintf("exponential:%s", strings.Join(parts, ","))
}

func parseExponentialBucketer(args map[string]string) (BucketingStrategy, error) {
	base := 2.0
	if arg, ok := args["base"]; ok {
		var err error
		base, err = parseExponentialBase(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid base %q: %w", arg, err)
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
	return ExponentialBucketer(base, origin)
}

func parseExponentialBase(value string) (float64, error) {
	trimmed := strings.TrimSpace(value)
	if strings.EqualFold(trimmed, "e") {
		return math.E, nil
	}
	return strconv.ParseFloat(trimmed, 64)
}

func init() {
	RegisterParser("exponential", parseExponentialBucketer)
	RegisterParser("exp", parseExponentialBucketer)
}
