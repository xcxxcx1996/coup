package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type DrawThreeCoins struct {
	message runtime.MatchData
}

func (a DrawThreeCoins) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	myTurn := message.GetUserId() == state.CurrentPlayerID

	if !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.message = message
	state.Actions.Push(a)

	buf, _ := global.Marshaler.Marshal(&api.Info{Info: fmt.Sprintf("%v 发动公爵技能，想拿三块钱", message.GetUsername())})

	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterQuestion()

}

// 只有质疑没有阻止
func (a DrawThreeCoins) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.PlayerInfos[state.CurrentPlayerID].Coins += 3

	buf, _ := global.Marshaler.Marshal(&api.Info{Info: fmt.Sprintf("%v 成功拿了三块钱", a.message.GetUsername())})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
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
