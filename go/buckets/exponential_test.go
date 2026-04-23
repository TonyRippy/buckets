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

func TestExponentialBucketerParse(t *testing.T) {
	for _, tc := range loadExponentialParseTestCases(t) {
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
			assertExponentialBucketerEquals(t, tc.want, got)
			if got.String() != tc.canonical {
				t.Errorf("expected canonical string %q, got %q", tc.canonical, got.String())
			}
		})
	}
}

func TestExponentialOutOfCoverage(t *testing.T) {
	eb, err := ExponentialBucketer(2, 0)
	if err != nil {
		t.Fatalf("ExponentialBucketer positive: %v", err)
	}
	if _, err := eb.IndexOf(-0.001); err == nil {
		t.Fatalf("expected error for value below positive exponential coverage")
	}
}

func assertExponentialBucketerEquals(t *testing.T, want exponentialParseExpectation, got BucketingStrategy) {
	t.Helper()
	switch want.kind {
	case "exponential":
		exponential, ok := got.(*exponentialBucketer)
		if !ok {
			t.Fatalf("expected exponentialBucketer, got %T", got)
		}
		if exponential.Base != want.base {
			t.Errorf("expected base %v, got %v", want.base, exponential.Base)
		}
		if exponential.Origin != want.origin {
			t.Errorf("expected origin %v, got %v", want.origin, exponential.Origin)
		}
	default:
		t.Fatalf("unknown kind %q", want.kind)
	}
}

type exponentialParseTestCase struct {
	file          string
	line          int
	spec          string
	wantErr       bool
	errorContains string
	want          exponentialParseExpectation
	canonical     string
}

type exponentialParseExpectation struct {
	kind   string
	base   float64
	origin float64
}

func (tc exponentialParseTestCase) Name() string {
	return fmt.Sprintf("%s:%d", tc.file, tc.line)
}

func loadExponentialParseTestCases(t *testing.T) []exponentialParseTestCase {
	t.Helper()

	path := filepath.Join(testCaseDirectory(t), "exponential_parse.csv")
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
	kindCol := requiredColumn(t, path, columns, "kind")
	baseCol := requiredColumn(t, path, columns, "base")
	originCol := requiredColumn(t, path, columns, "origin")
	canonicalCol := requiredColumn(t, path, columns, "canonical")

	testCases := make([]exponentialParseTestCase, 0, len(records)-1)
	fileName := filepath.Base(path)
	for i, fields := range records[1:] {
		lineNo := i + 2
		if len(fields) != len(headers) {
			t.Fatalf("fixture %q line %d: expected %d fields, got %d", path, lineNo, len(headers), len(fields))
		}

		spec := fields[specCol]
		wantErr, err := strconv.ParseBool(strings.TrimSpace(fields[wantErrCol]))
		if err != nil {
			t.Fatalf("%s:%d: parse want_error: %v", path, lineNo, err)
		}

		tc := exponentialParseTestCase{
			file:          fileName,
			line:          lineNo,
			spec:          spec,
			wantErr:       wantErr,
			errorContains: strings.TrimSpace(fields[errorContainsCol]),
		}

		if !wantErr {
			kind := strings.TrimSpace(fields[kindCol])
			if kind == "" {
				t.Fatalf("%s:%d: empty kind", path, lineNo)
			}

			base, err := strconv.ParseFloat(strings.TrimSpace(fields[baseCol]), 64)
			if err != nil {
				t.Fatalf("%s:%d: parse base: %v", path, lineNo, err)
			}

			origin := 0.0
			if strings.TrimSpace(fields[originCol]) != "" {
				origin, err = strconv.ParseFloat(strings.TrimSpace(fields[originCol]), 64)
				if err != nil {
					t.Fatalf("%s:%d: parse origin: %v", path, lineNo, err)
				}
			}

			tc.want = exponentialParseExpectation{
				kind:   kind,
				base:   base,
				origin: origin,
			}
			tc.canonical = fields[canonicalCol]
		}

		testCases = append(testCases, tc)
	}

	return testCases
}
