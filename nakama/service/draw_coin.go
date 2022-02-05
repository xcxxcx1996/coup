package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// 开始拿钱，进入是否阻止拿钱阶段
func (serv *MatchService) BeforeDrawCoin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.GetCoin{}
	err := serv.Unmarshaler.Unmarshal(message.GetData(), msg)

	if err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}

	var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_DRAW_COINS), ActionArg: msg}
	state.Actions.Push(&newAction)

	// 改变用户状态
	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)

	serv.EnterDenyCoins(nextPlayer, state)
}

// 完成拿钱动作，下一个人
func (serv *MatchService) CompleteDrawCoin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

	state.PlayerInfos[state.CurrentPlayerID].Coins += 3

	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

}
