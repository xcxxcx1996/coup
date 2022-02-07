package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type Assassin struct {
	Assassinor   string
	Assassinated string
	IsDeny       bool
}

func (a Assassin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	// 获得信息、核验
	msg := &api.Assassin{}
	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	ok := message.GetUserId() == state.CurrentPlayerID
	if err != nil || !ok {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 推进行动
	a.Assassinated = msg.PlayerId
	a.Assassinor = message.GetUserId()
	state.Actions.Push(a)

	//后处理

	info := fmt.Sprintf("%v 对 %v 刺杀", message.GetUsername(), state.GetPlayerNameByID(a.Assassinated))
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)

	state.EnterQuestion()
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Assassin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("质疑结束，开始阻止阶段")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterDenyAssassin(a.Assassinated)
}

func (a Assassin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("阻止结束，开始刺杀")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterDicardState(a.Assassinated)
}

// 行动被停止，删除，然后跳过
func (a Assassin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("刺杀被停止")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	state.NextTurn()
}

func (a Assassin) GetRole() api.Role {
	return api.Role_ASSASSIN
}
