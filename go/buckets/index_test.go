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

type indexTestCase struct {
	file      string
	line      int
	spec      string
	value     float64
	wantIndex int32
}

func (tc indexTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func TestIndexes(t *testing.T) {
	for _, tc := range loadIndexTestCases(t) {
		tc := tc
		t.Run(tc.Name(), func(t *testing.T) {
			bucketer, err := Parse(tc.spec)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tc.spec, err)
			}

			gotIndex, err := bucketer.IndexOf(tc.value)
			if err != nil {
				t.Fatalf("IndexOf: %v", err)
			}
			if gotIndex != tc.wantIndex {
				t.Fatalf("expected index %d, got %d", tc.wantIndex, gotIndex)
			}

			r, err := bucketer.Range(tc.wantIndex)
			if err != nil {
				t.Fatalf("Range: %v", err)
			}
			if !r.Contains(tc.value) {
				t.Fatalf("expected %v to be in range %v", tc.value, r)
			}
		})
	}
}

func loadIndexTestCases(t *testing.T) []indexTestCase {
	t.Helper()

	pattern := filepath.Join(testCaseDirectory(t), "*_index.csv")
	paths, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("glob index fixtures %q: %v", pattern, err)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		t.Fatalf("no index fixtures found matching %q", pattern)
	}

	testCases := make([]indexTestCase, 0)
	for _, path := range paths {
		testCases = append(testCases, loadIndexTestFile(t, path)...)
	}
	return testCases
}

func loadIndexTestFile(t *testing.T, path string) []indexTestCase {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open index fixture %q: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read index fixture %q: %v", path, err)
	}
	if len(records) == 0 {
		t.Fatalf("index fixture %q has no header row", path)
	}

	headers := make([]string, len(records[0]))
	columns := make(map[string]int, len(records[0]))
	for i, field := range records[0] {
		header := strings.TrimSpace(field)
		if header == "" {
			t.Fatalf("index fixture %q line 1: empty header", path)
		}
		headers[i] = header
		columns[header] = i
	}

	specCol := requiredColumn(t, path, columns, "spec")
	valueCol := requiredColumn(t, path, columns, "value")
	indexCol := requiredColumn(t, path, columns, "index")

	testCases := make([]indexTestCase, 0, len(records)-1)
	fileName := filepath.Base(path)
	for i, fields := range records[1:] {
		lineNo := i + 2
		if len(fields) != len(headers) {
			t.Fatalf(
				"index fixture %q line %d: expected %d fields, got %d",
				path,
				lineNo,
				len(headers),
				len(fields),
			)
		}

		spec := strings.TrimSpace(fields[specCol])
		if spec == "" {
			t.Fatalf("index fixture %q line %d: empty spec", path, lineNo)
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(fields[valueCol]), 64)
		if err != nil {
			t.Fatalf("index fixture %q line %d: parse value: %v", path, lineNo, err)
		}

		index, err := strconv.ParseInt(strings.TrimSpace(fields[indexCol]), 10, 32)
		if err != nil {
			t.Fatalf("index fixture %q line %d: parse index: %v", path, lineNo, err)
		}

		testCases = append(testCases, indexTestCase{
			file:      fileName,
			line:      lineNo,
			spec:      spec,
			value:     value,
			wantIndex: int32(index),
		})
	}
	return testCases
}
