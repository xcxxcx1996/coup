package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type DrawThreeCoins struct {
	message runtime.MatchData
}

func (a DrawThreeCoins) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	valid := ValidAction(state, message, api.State_START, nil)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.message = message
	state.Actions.Push(a)
	info := fmt.Sprintf("%v claims the barron, want to gain 3 coins", message.GetUsername())
	SendNotification(info, dispatcher)
	// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: fmt.Sprintf()})

	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterQuestion()

}

// 只有质疑没有阻止
func (a DrawThreeCoins) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.PlayerInfos[state.CurrentPlayerID].Coins += 3
	info := fmt.Sprintf("%v successful gain 3 coins", a.message.GetUsername())
	SendNotification(info, dispatcher)

	// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: })
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	state.NextTurn()
}

// 阻止失败，开始拿钱
func (a DrawThreeCoins) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

func (a DrawThreeCoins) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.Actions.Pop()
	state.NextTurn()
}

func (a DrawThreeCoins) GetRole() api.Role {
	return api.Role_BARON
}
