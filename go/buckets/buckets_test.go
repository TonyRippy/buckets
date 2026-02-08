package buckets

import "testing"

func assertRangeEquals(t *testing.T, want, got Range) {
	t.Helper()
	if got.From != want.From {
		t.Errorf("expected From %v, got %v", want.From, got.From)
	}
	if got.To != want.To {
		t.Errorf("expected To %v, got %v", want.To, got.To)
	}
	if got.FromBound != want.FromBound {
		t.Errorf("expected FromBound %v, got %v", want.FromBound, got.FromBound)
	}
	if got.ToBound != want.ToBound {
		t.Errorf("expected ToBound %v, got %v", want.ToBound, got.ToBound)
	}
}
