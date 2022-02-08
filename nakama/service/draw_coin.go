package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type DrawCoin struct {
	message runtime.MatchData
	coins   int32
}

func (a DrawCoin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.GetCoin{}
	valid := ValidAction(state, message, api.State_START, msg)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.coins = msg.Coins
	a.message = message
	state.Actions.Push(a)

	if a.coins == 2 {
		info := fmt.Sprintf("%v want to gain 2 coins", a.message.GetUsername())
		SendNotification(info, dispatcher)
		state.EnterDenyMoney()
	} else {
		info := "%v want to gain 1 coins"
		SendNotification(info, dispatcher)
		a.AfterDeny(dispatcher, state)
	}
}

// 正常拿钱没有人会质疑，只有公爵阻止,所以不用写
func (a DrawCoin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

// 阻止失败，开始拿钱
func (a DrawCoin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.Actions.Pop()
	state.PlayerInfos[state.CurrentPlayerID].Coins += int64(a.coins)
	info := fmt.Sprintf("%v success get the coins", a.message.GetUsername())
	SendNotification(info, dispatcher)

	// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: info})
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.NextTurn()
}

// 行动被阻止
func (a DrawCoin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.Actions.Pop()
	state.NextTurn()
}

func (a DrawCoin) GetRole() api.Role {
	return api.Role_UNROLE
}
