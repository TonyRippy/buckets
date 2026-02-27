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
	"strconv"
	"strings"
	"testing"
)

func TestAlignmentString(t *testing.T) {
	for _, tc := range loadAlignmentStringTestCases(t) {
		tc := tc
		t.Run(tc.Name(), func(t *testing.T) {
			side := Alignment(tc.side)
			want := tc.want
			if got := side.String(); got != want {
				t.Fatalf("expected %q, got %q", want, got)
			}
		})
	}
}

type alignmentStringTestCase struct {
	file string
	line int
	side int
	want string
}

func (tc alignmentStringTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func loadAlignmentStringTestCases(t *testing.T) []alignmentStringTestCase {
	t.Helper()

	path := filepath.Join(testCaseDirectory(t), "alignment_string.csv")
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open fixture %q: %v", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read fixture %q: %v", path, err)
	}
	if len(records) == 0 {
		t.Fatalf("fixture %q has no header row", path)
	}

	headers := make([]string, len(records[0]))
	columns := make(map[string]int, len(records[0]))
	for i, field := range records[0] {
		header := strings.TrimSpace(field)
		if header == "" {
			t.Fatalf("fixture %q line 1: empty header", path)
		}
		headers[i] = header
		columns[header] = i
	}

	sideCol := requiredColumn(t, path, columns, "side")
	wantCol := requiredColumn(t, path, columns, "want")

	testCases := make([]alignmentStringTestCase, 0, len(records)-1)
	fileName := filepath.Base(path)
	for i, fields := range records[1:] {
		lineNo := i + 2
		if len(fields) != len(headers) {
			t.Fatalf("fixture %q line %d: expected %d fields, got %d", path, lineNo, len(headers), len(fields))
		}

		side, err := strconv.Atoi(strings.TrimSpace(fields[sideCol]))
		if err != nil {
			t.Fatalf("fixture %q line %d: parse side: %v", path, lineNo, err)
		}

		want := fields[wantCol]
		testCases = append(testCases, alignmentStringTestCase{
			file: fileName,
			line: lineNo,
			side: side,
			want: want,
		})
	}
	return testCases
}

