package buckets

import (
	"fmt"
	"math"
	"strconv"
)

type floorBucketer struct {
	Width       float64
}

func FloorBucketer(width float64) (BucketingStrategy, error) {
	if width <= 0 {
		return nil, fmt.Errorf("invalid width %g", width)
	}
	return &floorBucketer{Width: width}, nil
}

func (b *floorBucketer) IndexOf(value float64) (int32, error) {
	if b.Width <= 0 {
		return 0, fmt.Errorf("invalid width %g", b.Width)
	}
	index := int32(math.Floor(value / b.Width))
	return index, nil
}

func (b *floorBucketer) Range(index int32) (Range, error) {
	min := float64(index)*b.Width
	max := min + b.Width
	return Range{From: min, To: max, FromBound: Closed, ToBound: Open}, nil
}

func (b *floorBucketer) String() string {
	if b.Width == 1 {
		return "floor"
	}
	return fmt.Sprintf("floor:width=%g", b.Width)
}

func init() {
	RegisterParser("floor", func(args map[string]string) (BucketingStrategy, error) {
		width := 1.0
		if arg, ok := args["width"]; ok {
			var err error
			width, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid width %q", arg)
			}
		}
		return FloorBucketer(width)
	})
}
