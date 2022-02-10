package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type DrawThreeCoins struct {
	message runtime.MatchData
}

func (a DrawThreeCoins) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {

	if err = ValidAction(state, message, api.State_START, nil); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.message = message
	state.Actions.Push(a)
	info := fmt.Sprintf("%v claims the barron and want to gain 3 coins.", message.GetUsername())
	SendNotification(info, dispatcher)
	state.EnterQuestion()
	return
}

// 只有质疑没有阻止
func (a DrawThreeCoins) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	err = state.GainCoins(state.CurrentPlayerID, 3)
	if err != nil {
		return
	}
	info := fmt.Sprintf("%v successful gain 3 coins.", a.message.GetUsername())
	SendNotification(info, dispatcher)
	state.NextTurn()
	return
}

// 没有阻止函数
func (a DrawThreeCoins) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	return
}

// 冒充公爵被质疑成功了
func (a DrawThreeCoins) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintf("%v was denied to get 3 coins.", a.message.GetUsername())
	SendNotification(info, dispatcher)
	return
}

func (a DrawThreeCoins) GetRole() api.Role {
	return api.Role_BARON
}

func (c DrawThreeCoins) GetActor() string {
	return c.message.GetUserId()
}
