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
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
)

func assertLinearBucketerEquals(t *testing.T, want linearParseExpectation, got BucketingStrategy) {
	t.Helper()
	linear, ok := got.(*linearBucketer)
	if !ok {
		t.Fatalf("expected linearBucketer, got %T", got)
	}
	if linear.M != want.m {
		t.Errorf("expected slope %v, got %v", want.m, linear.M)
	}
	if linear.B != want.b {
		t.Errorf("expected intercept %v, got %v", want.b, linear.B)
	}
}

func TestLinearBucketerParse(t *testing.T) {
	for _, tc := range loadLinearParseTestCases(t) {
		tc := tc
		t.Run(tc.Name(), func(t *testing.T) {
			got, err := Parse(tc.spec)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Fatalf("expected error containing %q, got %q", tc.errorContains, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertLinearBucketerEquals(t, tc.want, got)
			if got.String() != tc.canonical {
				t.Errorf("expected canonical string %q, got %q", tc.canonical, got.String())
			}
		})
	}
}

type linearParseTestCase struct {
	file          string
	line          int
	spec          string
	wantErr       bool
	errorContains string
	want          linearParseExpectation
	canonical     string
}

type linearParseExpectation struct {
	m         float64
	b         float64
}

func (tc linearParseTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func loadLinearParseTestCases(t *testing.T) []linearParseTestCase {
	t.Helper()
	return slices.Concat(
		loadLinearParseTestFile(t, "linear_parse.csv"),
		loadLinearParseTestFile(t, "fixed_parse.csv"))
}

func loadLinearParseTestFile(t *testing.T, filename string) []linearParseTestCase {
	t.Helper()

	path := filepath.Join(testCaseDirectory(t), filename)
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture %q: %v", filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read fixture %q: %v", filename, err)
	}
	if len(records) == 0 {
		t.Fatalf("fixture %q has no header row", filename)
	}

	headers := make([]string, len(records[0]))
	columns := make(map[string]int, len(records[0]))
	for i, field := range records[0] {
		header := strings.TrimSpace(field)
		if header == "" {
			t.Fatalf("fixture %q line 1: empty header", filename)
		}
		headers[i] = header
		columns[header] = i
	}

	specCol := requiredColumn(t, path, columns, "spec")
	wantErrCol := requiredColumn(t, path, columns, "want_error")
	errorContainsCol := requiredColumn(t, path, columns, "error_contains")
	mCol := requiredColumn(t, path, columns, "m")
	bCol := requiredColumn(t, path, columns, "b")
	canonicalCol := requiredColumn(t, path, columns, "canonical")

	testCases := make([]linearParseTestCase, 0, len(records)-1)
	for i, fields := range records[1:] {
		lineNo := i + 2
		if len(fields) != len(headers) {
			t.Fatalf("%s:%d: expected %d fields, got %d", filename, lineNo, len(headers), len(fields))
		}

		spec := fields[specCol]
		wantErr, err := strconv.ParseBool(strings.TrimSpace(fields[wantErrCol]))
		if err != nil {
			t.Fatalf("%s:%d: parse want_error: %v", filename, lineNo, err)
		}
		errorContains := strings.TrimSpace(fields[errorContainsCol])

		tc := linearParseTestCase{
			file:          filename,
			line:          lineNo,
			spec:          spec,
			wantErr:       wantErr,
			errorContains: errorContains,
		}

		if !wantErr {
			m, err := strconv.ParseFloat(strings.TrimSpace(fields[mCol]), 64)
			if err != nil {
				t.Fatalf("%s:%d: parse m: %v", filename, lineNo, err)
			}
			b, err := strconv.ParseFloat(strings.TrimSpace(fields[bCol]), 64)
			if err != nil {
				t.Fatalf("%s:%d: parse b: %v", filename, lineNo, err)
			}
			tc.want = linearParseExpectation{
				m:         m,
				b:         b,
			}
			tc.canonical = fields[canonicalCol]
		}

		testCases = append(testCases, tc)
	}
	return testCases
}

func TestLinearRangesShareIdenticalBoundaries(t *testing.T) {
	for _, spec := range []string{"linear:m=0.1", "linear:m=0.1,align=left", "linear:m=0.2,b=0.1", "linear:m=0.2,b=0.1,align=left", "fixed:width=0.01"} {
		t.Run(spec, func(t *testing.T) {
			strategy, err := Parse(spec)
			if err != nil {
				t.Fatal(err)
			}
			for index := int32(-1000); index < 1000; index++ {
				lower, err := strategy.Range(index)
				if err != nil {
					t.Fatal(err)
				}
				upper, err := strategy.Range(index + 1)
				if err != nil {
					t.Fatal(err)
				}
				if lower.To != upper.From {
					t.Fatalf("adjacent buckets %d and %d disagree: %.17g != %.17g", index, index+1, lower.To, upper.From)
				}
			}
		})
	}
}

func TestLinearRangesDoNotWrapEndpointIndexes(t *testing.T) {
	for _, alignment := range []Alignment{Left, Right} {
		strategy, err := LinearBucketer(1, 0, alignment)
		if err != nil {
			t.Fatal(err)
		}
		for _, index := range []int32{math.MinInt32, math.MaxInt32} {
			interval, err := strategy.Range(index)
			if err != nil {
				t.Fatal(err)
			}
			if interval.To-interval.From != 1 {
				t.Fatalf("index %d alignment %v has invalid range %v", index, alignment, interval)
			}
		}
	}
}
