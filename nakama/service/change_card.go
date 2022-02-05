package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) BeforeChangeCard(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_CHANGE_CARD)}
	state.Actions.Push(&newAction)
	// question状态
	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
	serv.EnterQuestionState(nextPlayer, state)
}

func (serv *MatchService) CompleteChangeCard(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	//更换牌
	_, ok := state.Actions.Pop()
	if !ok {
		//丢失
		return
	}

	handcards := state.PlayerInfos[state.CurrentPlayerID].Cards
	life := len(handcards)
	newcards := state.Deck[:life]
	state.Deck = append(state.Deck[life:], handcards...)
	state.PlayerInfos[state.CurrentPlayerID].Cards = newcards
	//下一个回合
}
