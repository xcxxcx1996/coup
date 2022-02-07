package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type Steal struct {
	Victim string
	Thief  string
}

func (a Steal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.StealCoins{}
	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.Victim = msg.PlayerId
	a.Thief = message.GetUserId()
	state.Actions.Push(a)
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: fmt.Sprintf("%v 想偷 %v 的金币", message.GetUsername(), state.GetPlayerNameByID(a.Victim))})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true) // question状态
	state.EnterQuestion()
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Steal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("质疑结束，开始阻止阶段")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterDenySteal(a.Victim)
}

func (a Steal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("阻止结束，开始偷钱")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	coins := state.PlayerInfos[a.Victim].Coins
	if coins <= 2 {
		state.PlayerInfos[a.Victim].Coins = 0
		state.PlayerInfos[state.CurrentPlayerID].Coins += coins
	} else {
		state.PlayerInfos[a.Victim].Coins -= 2
		state.PlayerInfos[state.CurrentPlayerID].Coins += 2
	}
}

func (a Steal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("偷钱行为被停止")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	state.NextTurn()
}

func (a Steal) GetRole() api.Role {
	return api.Role_ASSASSIN
}
