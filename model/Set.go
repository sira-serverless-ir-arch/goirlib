package model

type Set struct {
	data map[string]bool
}

func NewSet() *Set {
	s := &Set{}
	s.data = make(map[string]bool)
	return s
}

func (s *Set) GetData() map[string]bool {
	return s.data
}

func (s *Set) Add(element string) {
	s.data[element] = true
}

func (s *Set) Remove(element string) {
	delete(s.data, element)
}

func (s *Set) Contains(element string) bool {
	_, exists := s.data[element]
	return exists
}

func (s *Set) Size() int {
	return len(s.data)
}
