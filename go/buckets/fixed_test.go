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

func assertFixedBucketerEquals(t *testing.T, want *fixedBucketer, got BucketingStrategy) {
	t.Helper()
	fixed, ok := got.(*fixedBucketer)
	if !ok {
		t.Fatalf("expected fixedBucketer, got %T", got)
	}
	if fixed.Width != want.Width {
		t.Errorf("expected width %v, got %v", want.Width, fixed.Width)
	}
	if fixed.Origin != want.Origin {
		t.Errorf("expected origin %v, got %v", want.Origin, fixed.Origin)
	}
	if fixed.Alignment != want.Alignment {
		t.Errorf("expected alignment %v, got %v", want.Alignment, fixed.Alignment)
	}
}

func TestFixedBucketerParse(t *testing.T) {
	for _, tc := range loadFixedParseTestCases(t) {
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
			assertFixedBucketerEquals(t, tc.want, got)
			if got.String() != tc.canonical {
				t.Errorf("expected canonical string %q, got %q", tc.canonical, got.String())
			}
		})
	}
}

func TestFixedBucketerInvalidAlignment(t *testing.T) {
	_, err := FixedBucketer(1, 0, Alignment(255))
	if err == nil {
		t.Fatalf("expected error")
	}
}

type fixedParseTestCase struct {
	file          string
	line          int
	spec          string
	wantErr       bool
	errorContains string
	want          *fixedBucketer
	canonical     string
}

func (tc fixedParseTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func loadFixedParseTestCases(t *testing.T) []fixedParseTestCase {
	t.Helper()

	path := filepath.Join(testCaseDirectory(t), "fixed_parse.csv")
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

	specCol := requiredColumn(t, path, columns, "spec")
	wantErrCol := requiredColumn(t, path, columns, "want_error")
	errorContainsCol := requiredColumn(t, path, columns, "error_contains")
	widthCol := requiredColumn(t, path, columns, "width")
	originCol := requiredColumn(t, path, columns, "origin")
	alignmentCol := requiredColumn(t, path, columns, "alignment")
	canonicalCol := requiredColumn(t, path, columns, "canonical")

	testCases := make([]fixedParseTestCase, 0, len(records)-1)
	fileName := filepath.Base(path)
	for i, fields := range records[1:] {
		lineNo := i + 2
		if len(fields) != len(headers) {
			t.Fatalf("fixture %q line %d: expected %d fields, got %d", path, lineNo, len(headers), len(fields))
		}

		spec := fields[specCol]
		wantErr, err := strconv.ParseBool(strings.TrimSpace(fields[wantErrCol]))
		if err != nil {
			t.Fatalf("%s:%d: parse want_error: %v", fileName, lineNo, err)
		}
		errorContains := strings.TrimSpace(fields[errorContainsCol])

		tc := fixedParseTestCase{
			file:          fileName,
			line:          lineNo,
			spec:          spec,
			wantErr:       wantErr,
			errorContains: errorContains,
		}

		if !wantErr {
			width, err := strconv.ParseFloat(strings.TrimSpace(fields[widthCol]), 64)
			if err != nil {
				t.Fatalf("%s:%d: parse width: %v", fileName, lineNo, err)
			}
			origin, err := strconv.ParseFloat(strings.TrimSpace(fields[originCol]), 64)
			if err != nil {
				t.Fatalf("%s:%d: parse origin: %v", fileName, lineNo, err)
			}
			alignment, err := ParseAlignment(strings.TrimSpace(fields[alignmentCol]))
			if err != nil {
				t.Fatalf("%s:%d: parse alignment: %v", fileName, lineNo, err)
			}
			tc.want = &fixedBucketer{
				Width:     width,
				Origin:    origin,
				Alignment: alignment,
			}
			tc.canonical = fields[canonicalCol]
		}

		testCases = append(testCases, tc)
	}
	return testCases
}
