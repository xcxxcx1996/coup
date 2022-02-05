package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// 开始偷钱，进入质询阶段
func (serv *MatchService) BeforeStealCoins(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.StealCoins{}
	err := serv.Unmarshaler.Unmarshal(message.GetData(), msg)

	if err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}

	var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_STEAL_COINS), ActionArg: msg}
	state.Actions.Push(&newAction)

	// question状态
	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
	serv.EnterQuestionState(nextPlayer, state)

}
func (serv *MatchService) AfterStealCard(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	// 去除某些人的卡
	action, ok := state.Actions.Pop()
	if !ok {
		//丢失
		return
	}
	arg := action.ActionArg.(*api.StealCoins)

	coins := state.PlayerInfos[arg.PlayerId].Coins
	if coins <= 2 {
		state.PlayerInfos[arg.PlayerId].Coins = 0
		state.PlayerInfos[state.CurrentPlayerID].Coins += coins
	} else {
		state.PlayerInfos[arg.PlayerId].Coins -= 2
		state.PlayerInfos[state.CurrentPlayerID].Coins += 2
	}
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

}
