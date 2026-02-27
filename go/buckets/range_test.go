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
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

type rangeTestCase struct {
	file  string
	line  int
	spec  string
	index int32
	want  Range
}

func (tc rangeTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func TestRanges(t *testing.T) {
	for _, tc := range loadRangeTestCases(t) {
		tc := tc
		t.Run(tc.Name(), func(t *testing.T) {
			bucketer, err := Parse(tc.spec)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.spec, err)
			}

			got, err := bucketer.Range(tc.index)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			assertRangeEquals(t, tc.want, got)
		})
	}
}

func loadRangeTestCases(t *testing.T) []rangeTestCase {
	t.Helper()

	pattern := filepath.Join(testCaseDirectory(t), "*_range.csv")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob range fixtures %q: %v", pattern, err)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		t.Fatalf("no range fixtures found matching %q", pattern)
	}

	testCases := make([]rangeTestCase, 0)
	for _, path := range paths {
		testCases = append(testCases, loadRangeTestFile(t, path)...)
	}
	return testCases
}

func loadRangeTestFile(t *testing.T, path string) []rangeTestCase {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open range fixture %q: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read range fixture %q: %v", path, err)
	}
	if len(records) == 0 {
		t.Fatalf("range fixture %q has no header row", path)
	}

	headers := make([]string, len(records[0]))
	columns := make(map[string]int, len(records[0]))
	for i, field := range records[0] {
		header := strings.TrimSpace(field)
		if header == "" {
			t.Fatalf("range fixture %q line 1: empty header", path)
		}
		headers[i] = header
		columns[header] = i
	}

	specCol := requiredColumn(t, path, columns, "spec")
	indexCol := requiredColumn(t, path, columns, "index")
	fromCol := requiredColumn(t, path, columns, "from")
	toCol := requiredColumn(t, path, columns, "to")
	fromBoundCol := requiredColumn(t, path, columns, "from_bound")
	toBoundCol := requiredColumn(t, path, columns, "to_bound")

	testCases := make([]rangeTestCase, 0, len(records)-1)
	filename := filepath.Base(path)
	for i, fields := range records[1:] {
		line := i + 2
		if len(fields) != len(headers) {
			t.Fatalf(
				"range fixture %q line %d: expected %d fields, got %d",
				path,
				line,
				len(headers),
				len(fields),
			)
		}

		spec := strings.TrimSpace(fields[specCol])
		if spec == "" {
			t.Fatalf(`%q:%d: empty spec`, path, line)
		}

		index, err := strconv.ParseInt(strings.TrimSpace(fields[indexCol]), 10, 32)
		if err != nil {
			t.Fatalf(`%q:%d: parsing index: %v`, path, line, err)
		}

		from, err := strconv.ParseFloat(strings.TrimSpace(fields[fromCol]), 64)
		if err != nil {
			t.Fatalf(`%q:%d: parsing from: %v`, path, line, err)
		}
		fromBound, err := parseBoundType(strings.TrimSpace(fields[fromBoundCol]))
		if err != nil {
			t.Fatalf(`%q:%d: parsing from bound: %v`, path, line, err)
		}

		to, err := strconv.ParseFloat(strings.TrimSpace(fields[toCol]), 64)
		if err != nil {
			t.Fatalf(`%q:%d: parsing to: %v`, path, line, err)
		}
		toBound, err := parseBoundType(strings.TrimSpace(fields[toBoundCol]))
		if err != nil {
			t.Fatalf(`%q:%d: parsing to bound: %v`, path, line, err)
		}

		testCases = append(testCases, rangeTestCase{
			file:  filename,
			line:  line,
			spec:  spec,
			index: int32(index),
			want: Range{
				From:      from,
				To:        to,
				FromBound: fromBound,
				ToBound:   toBound,
			},
		})
	}
	return testCases
}
