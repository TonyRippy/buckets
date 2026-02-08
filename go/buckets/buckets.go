package buckets

import (
	"fmt"
)

type BoundType uint8

const (
	Open BoundType = iota
	Closed
)

type Range struct {
	From      float64
	To        float64
	FromBound BoundType
	ToBound   BoundType
}

func (r Range) Contains(x float64) bool {
	switch r.FromBound {
	case Open:
		if x <= r.From {
			return false
		}
	case Closed:
		if x < r.From {
			return false
		}
	}
	switch r.ToBound {
	case Open:
		if x >= r.To {
			return false
		}
	case Closed:
		if x > r.To {
			return false
		}
	}
	return true
}

type BucketingStrategy interface {
	fmt.Stringer

	// IndexOf returns the index of the bucket that contains the given value.
	IndexOf(value float64) (int32, error)

	// Range returns the range of values that are in the bucket with the given index.
	Range(index int32) (Range, error)
}
