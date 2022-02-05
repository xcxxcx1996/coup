package model

type Action struct {
	ActionID  int64
	ActionArg interface{}
}
type Actions struct {
	data []*Action
}

func (s *Actions) Push(item *Action) (string, bool) {
	s.data = append(s.data, item)
	return "ok", true
}

func (s *Actions) Pop() (*Action, bool) {
	length := s.Length()
	var item *Action
	if length == 0 {
		return item, false
	} else if length == 1 {
		item = s.data[0]
		s.data = []*Action{}
		return item, true
	} else {
		item = s.data[0]
		s.data = s.data[1:]
		return item, true
	}
}

func (s *Actions) ElementAt(index int) (*Action, bool) {
	length := len(s.data)
	if index >= length {
		return &Action{}, false
	}
	return s.data[index], true
}

// 原则上讲，不应该存在的操作
func (s *Actions) UpdateElementAt(data *Action, index int) bool {
	length := len(s.data)
	if index >= length {
		return false
	}
	s.data[index] = data
	return true
}

// 原则上讲，不应该存在的操作
func (s *Actions) RemoveElementAt(index int) (*Action, bool) {
	length := len(s.data)
	if index >= length {
		return &Action{}, false
	}
	s.data = append(
		s.data[:index],
		s.data[index+1:]...,
	)
	return s.data[index], true
}

func (s *Actions) Length() int {
	return len(s.data)
}

func (s *Actions) Last() *Action {
	if s.Length() >= 1 {
		return s.data[len(s.data)-1]
	}
	return nil
}
