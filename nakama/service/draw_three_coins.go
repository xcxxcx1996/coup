package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) BeforeDrawThreeCoins(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

	var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_DRAW_THREE_COINS)}
	state.Actions.Push(&newAction)

	// question状态
	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
	serv.EnterQuestionState(nextPlayer, state)

}

// 完成拿钱，进入拿钱
func (serv *MatchService) CompleteDrawThreeCoins(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

	state.PlayerInfos[state.CurrentPlayerID].Coins += 3
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

}
