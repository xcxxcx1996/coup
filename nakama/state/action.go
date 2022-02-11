package state

import (
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
)

type ActionState int


const (
	ASSASSIN = iota
	NO_ONE_DENY
	BEFORE_CHANGE
	COMPLETE_CHANGE
)


type IAction interface {
	// 开始动作
	Start(runtime.MatchDispatcher, runtime.MatchData, *MatchState) error
	// 下一步动作
	AfterQuestion(runtime.MatchDispatcher, *MatchState) error
	AfterDeny(runtime.MatchDispatcher, *MatchState) error
	// 动作中止
	Stop(runtime.MatchDispatcher, *MatchState) error
	// getRole
	GetRole() api.Role
	GetActor() string
}


type Actions struct {
	data []IAction
}


func (s *Actions) Clear() {
	s.data = []IAction{}
}


func (s *Actions) Push(item IAction) {
	s.data = append(s.data, item)
}


func (s *Actions) Pop() (IAction, error) {
	length := s.Length()
	var item IAction
	if length == 0 {
		return item, errors.New("no record")
	} else if length == 1 {
		item = s.data[0]
		s.data = []IAction{}
		return item, nil
	} else {
		item = s.data[length-1]
		s.data = s.data[:length-1]
		return item, nil
	}
}


func (s *Actions) ElementAt(index int) (IAction, error) {
	length := len(s.data)
	if index >= length {
		return nil, errors.New("no record")
	}
	return s.data[index], nil
}

// 原则上讲，不应该存在的操作
func (s *Actions) UpdateElementAt(data IAction, index int) error {
	length := len(s.data)
	if index >= length {
		return errors.New("no record")
	}
	s.data[index] = data
	return nil
}

// 原则上讲，不应该存在的操作
func (s *Actions) RemoveElementAt(index int) (IAction, error) {
	length := len(s.data)
	if index >= length {
		return nil, errors.New("no record")
	}
	s.data = append(
		s.data[:index],
		s.data[index+1:]...,
	)
	return s.data[index], nil
}

func (s *Actions) Length() int {
	return len(s.data)
}

func (s *Actions) Last() (action IAction, err error) {
	if s.Length() >= 1 {
		return s.data[len(s.data)-1], nil
	}
	return nil, errors.New("no record")
}
