package model

import (
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

// // 动作
// type Action struct {
// 	ActionID  int64
// 	State     ActionState
// 	Role      api.Role
// 	ActionArg interface{}
// }

type IAction interface {
	// 开始动作
	Start(runtime.MatchDispatcher, runtime.MatchData, *MatchState)
	// 下一步动作
	AfterQuestion(runtime.MatchDispatcher, *MatchState)
	AfterDeny(runtime.MatchDispatcher, *MatchState)
	// 动作中止
	Stop(runtime.MatchDispatcher, *MatchState)
	// getRole
	GetRole() api.Role
}

type Actions struct {
	data []IAction
}

func (s *Actions) Push(item IAction) (string, bool) {
	s.data = append(s.data, item)
	return "ok", true
}

func (s *Actions) Pop() (IAction, bool) {
	length := s.Length()
	var item IAction
	if length == 0 {
		return item, false
	} else if length == 1 {
		item = s.data[0]
		s.data = []IAction{}
		return item, true
	} else {
		item = s.data[0]
		s.data = s.data[1:]
		return item, true
	}
}

func (s *Actions) ElementAt(index int) (IAction, bool) {
	length := len(s.data)
	if index >= length {
		return nil, false
	}
	return s.data[index], true
}

// 原则上讲，不应该存在的操作
func (s *Actions) UpdateElementAt(data IAction, index int) bool {
	length := len(s.data)
	if index >= length {
		return false
	}
	s.data[index] = data
	return true
}

// 原则上讲，不应该存在的操作
func (s *Actions) RemoveElementAt(index int) (IAction, bool) {
	length := len(s.data)
	if index >= length {
		return nil, false
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

func (s *Actions) Last() (action IAction, ok bool) {
	if s.Length() >= 1 {
		return s.data[len(s.data)-1], true
	}
	return nil, false
}
