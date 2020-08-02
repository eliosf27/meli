package slice

type SliceInt64 []int64

func (s SliceInt64) Avg() float64 {
	if s.Count() == 0 {
		return 0
	}

	return s.Sum() / float64(s.Count())
}

func (s SliceInt64) Count() int64 {

	return int64(len(s))
}

func (s SliceInt64) Sum() float64 {
	total := 0.0

	for _, val := range s {
		total += float64(val)
	}

	return total
}
