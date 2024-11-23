package consistenthash

import "sort"

type SortSlice struct {
	s []int
}

func (s *SortSlice) Add(v int) {
	s.s = append(s.s, v)
	for i := len(s.s) - 1; i > 0 && s.s[i] < s.s[i-1]; i-- {
		temp := s.s[i]
		s.s[i] = s.s[i-1]
		s.s[i-1] = temp
	}
}

func (s *SortSlice) Index(v int) int {
	idx := sort.Search(len(s.s), func(i int) bool {
		return s.s[i] >= v
	})
	return idx
}

func (s *SortSlice) Len() int {
	return len(s.s)
}

func (s *SortSlice) Get(index int) int {
	return s.s[index]
}

func (s *SortSlice) Delete(index int) {
	s.s = append(s.s[:index], s.s[index+1:]...)
}
