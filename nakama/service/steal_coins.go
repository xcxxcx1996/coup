package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type Steal struct {
	Victim string
	Thief  string
}

func (a Steal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.StealCoins{}
	valid := ValidAction(state, message, api.State_START, msg)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}

	a.Victim = msg.PlayerId
	a.Thief = message.GetUserId()
	state.Actions.Push(a)
	info := fmt.Sprintf("%v claims the Captain,want to steal %v 的金币", message.GetUsername(), state.GetPlayerNameByID(a.Victim))
	SendNotification(info, dispatcher)

	state.EnterQuestion()
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Steal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("question end, enter deny ")
	SendNotification(info, dispatcher)
	state.EnterDenySteal(a.Victim)
}

func (a Steal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("deny end, steal")
	SendNotification(info, dispatcher)

	coins := state.PlayerInfos[a.Victim].Coins
	if coins <= 2 {
		state.PlayerInfos[a.Victim].Coins = 0
		state.PlayerInfos[state.CurrentPlayerID].Coins += coins
	} else {
		state.PlayerInfos[a.Victim].Coins -= 2
		state.PlayerInfos[state.CurrentPlayerID].Coins += 2
	}
	state.Actions.Pop()
	state.NextTurn()
}

func (a Steal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("steal was stoped")
	SendNotification(info, dispatcher)

	state.Actions.Pop()
	state.NextTurn()
}

func (a Steal) GetRole() api.Role {
	return api.Role_CAPTAIN
}
