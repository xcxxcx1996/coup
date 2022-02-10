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
	info := fmt.Sprintf(`<p><span style="color:red;">%v</span> claims the <span style="color:red;">CAPTAIN</span> and want to steal <span style="color:red;">%v</span> coins.</p >`, message.GetUsername(), state.GetPlayerNameByID(a.Victim))
	SendNotification(info, dispatcher)
	state.EnterQuestion()
	return
}

// 1. 没人质疑 不会删，2. 有人质疑，但是失败，这个时候会把action提出来
func (a Steal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.Actions.Push(a)
	info := fmt.Sprintln(`<p>Questioning ends, Denying begins.</p >`)
	SendNotification(info, dispatcher)
	state.EnterDenySteal(a.Victim)
	return nil
}

//1. 没人阻止，偷钱，下一个回合，
//2 .有人阻止，但是阻止失败，deny steal
func (a Steal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintln(`<p><span style="color:green;">Successful</span> theft and refusal to end.</p >`)
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
		state.GainCoins(a.Thief, StealCoinsNum)
	}
	return
}

func (a Steal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintln(`<p>Steal was <span style="color:red;">stoped</span>.</p >`)
	SendNotification(info, dispatcher)
	return
}

func (a Steal) GetRole() api.Role {
	return api.Role_CAPTAIN
}

func (c Steal) GetActor() string {
	return c.Thief
}
