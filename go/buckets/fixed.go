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

// FixedBucketer returns a fixed-width bucketer with the given origin and closed side.
func FixedBucketer(width, origin float64, align Alignment) (BucketingStrategy, error) {
	if math.IsNaN(width) || math.IsInf(width, 0) || width <= 0 {
		return nil, fmt.Errorf("invalid width %g", width)
	}
	if math.IsNaN(origin) || math.IsInf(origin, 0) {
		return nil, fmt.Errorf("invalid origin %g", origin)
	}
	if align != Left && align != Right {
		return nil, fmt.Errorf("invalid alignment %d", align)
	}
	return &linearBucketer{M: width, B: origin, Alignment: align}, nil
}

func parseFixedBucketer(args map[string]string) (BucketingStrategy, error) {
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

	alignment := Right
	if arg, ok := args["align"]; ok {
		var err error
		alignment, err = ParseAlignment(arg)
		if err != nil {
			return nil, err
		}
	}

	return FixedBucketer(width, origin, alignment)
}

func init() {
	RegisterParser("fixed", parseFixedBucketer)
}
