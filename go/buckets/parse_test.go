package buckets

import (
	"strings"
	"testing"
)

func TestParseErrors(t *testing.T) {
	for _, test := range []struct {
		spec     string
		contains string
	}{
		{"", "empty"},
		{"unknown", "unknown"},
	} {
		t.Run(test.spec, func(t *testing.T) {
			_, err := Parse(test.spec)
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), test.contains) {
				t.Fatalf("expected error containing %q, got %q", test.contains, err.Error())
			}
		})
	}
}
