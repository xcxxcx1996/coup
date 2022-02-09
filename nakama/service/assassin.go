package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type Assassin struct {
	Assassinor   string
	Assassinated string
}

func (a Assassin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) error {
	// 获得信息、核验

	msg := &api.Assassin{}
	err := ValidAction(state, message, api.State_START, msg)
	// 推进行动
	if err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return err
	}
	a.Assassinated = msg.PlayerId
	a.Assassinor = message.GetUserId()
	state.Actions.Push(a)

	//后处理

	info := fmt.Sprintf("%v want to assassin %v ", message.GetUsername(), state.GetPlayerNameByID(a.Assassinated))
	SendNotification(info, dispatcher)
	state.EnterQuestion()
	return nil
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Assassin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	info := fmt.Sprintln("question end, deny start")
	SendNotification(info, dispatcher)
	state.EnterDenyAssassin(a.Assassinated)
	return nil
}

// 玩家进入刺杀阶段
func (a Assassin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	state.ActionComplete = true
	info := fmt.Sprintln("deny end, action execute")
	SendNotification(info, dispatcher)
	state.EnterDicardState(a.Assassinated)
	return nil
}

// 行动被停止 1.被成功质疑，2.被阻止
func (a Assassin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	state.ActionComplete = true
	info := fmt.Sprintln("assassin was stopped")
	SendNotification(info, dispatcher)
	state.Actions.Pop()
	return nil
}

func (a Assassin) GetActor() string {
	return a.Assassinor
}

func (a Assassin) GetRole() api.Role {
	return api.Role_ASSASSIN
}
