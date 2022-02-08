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

func (a Assassin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	// 获得信息、核验

	msg := &api.Assassin{}
	valid := ValidAction(state, message, api.State_START, msg)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.Assassinated = msg.PlayerId
	a.Assassinor = message.GetUserId()
	state.Actions.Push(a)

	//后处理

	info := fmt.Sprintf("%v want to assassin %v ", message.GetUsername(), state.GetPlayerNameByID(a.Assassinated))
	SendNotification(info, dispatcher)
	state.EnterQuestion()
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Assassin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("question end, deny start")
	SendNotification(info, dispatcher)

	state.EnterDenyAssassin(a.Assassinated)
}

// 玩家进入刺杀阶段
func (a Assassin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("deny end, action execute")
	SendNotification(info, dispatcher)

	state.EnterDicardState(a.Assassinated)
}

// 行动被停止，删除，然后跳过
func (a Assassin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("assassin was stopped")
	SendNotification(info, dispatcher)

	state.Actions.Pop()
	state.NextTurn()
}

func (a Assassin) GetRole() api.Role {
	return api.Role_ASSASSIN
}
