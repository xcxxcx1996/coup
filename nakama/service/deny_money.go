package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// 是否阻止，如果是进入质询阶段，如果否下一个人
func (serv *MatchService) BeforeDenyMoney(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.DenyMoney{}
	// get msg
	if err := serv.Unmarshaler.Unmarshal(message.GetData(), msg); err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}
	if msg.IsDeny {
		var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_DENY_MONEY), ActionArg: msg}
		state.Actions.Push(&newAction)
		// question状态
		nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
		serv.EnterQuestionState(nextPlayer, state)
		return
	}
	// 不阻止的话下一个人

}

// 完成某人阻止
func (serv *MatchService) CompletedDenyMoney(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

}

// 进入
func (serv *MatchService) EnterDenyCoins(playerID string, state *model.MatchState) {

}
