package slice

type Int64 []int64

func (s Int64) Avg() float64 {
	if s.Count() == 0 {
		return 0
	}

	return s.Sum() / float64(s.Count())
}

func (s Int64) Count() int64 {

	return int64(len(s))
}

func (s Int64) Sum() float64 {
	total := 0.0

	for _, val := range s {
		total += float64(val)
	}

	return total
}
