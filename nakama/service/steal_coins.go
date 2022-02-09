package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type Steal struct {
	Victim string
	Thief  string
}

const StealCoinsNum = 2

func (a Steal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.StealCoins{}

	if err = ValidAction(state, message, api.State_START, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}

	a.Victim = msg.PlayerId
	a.Thief = message.GetUserId()

	state.Actions.Push(a)
	info := fmt.Sprintf("%v claims the Captain, want to steal %v coins", message.GetUsername(), state.GetPlayerNameByID(a.Victim))
	SendNotification(info, dispatcher)
	state.EnterQuestion()
	return
}

// 1. 没人质疑，2. 有人质疑，但失败
func (a Steal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	info := fmt.Sprintln("question end, enter deny.")
	SendNotification(info, dispatcher)
	state.EnterDenySteal(a.Victim)
	return nil
}

func (a Steal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintln("deny end, steal")
	SendNotification(info, dispatcher)
	resCoins, err := state.GetCoins(a.Victim)
	if err != nil {
		return
	}
	if resCoins <= StealCoinsNum {
		err = state.SetCoins(a.Victim, 0)
		if err != nil {
			return
		}
		err = state.GainCoins(a.Thief, resCoins)
		if err != nil {
			return
		}
	} else {
		state.LoseCoins(a.Victim, StealCoinsNum)
		state.GainCoins(a.Thief, resCoins)
	}
	return
}

func (a Steal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintln("steal was stoped")
	SendNotification(info, dispatcher)
	return
}

func (a Steal) GetRole() api.Role {
	return api.Role_CAPTAIN
}

func (c Steal) GetActor() string {
	return c.Thief
}
