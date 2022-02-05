package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// 被别人看出来了，然后把自己的牌放回并重新抽一张
func (serv *MatchService) ChangeSingleCard(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	// handcards := state.PlayerInfos[state.CurrentPlayerID].Cards
	msg := &api.ChangeSingleCard{}
	if err := serv.Unmarshaler.Unmarshal(message.GetData(), msg); err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	cards := state.PlayerInfos[msg.PlayerId].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			state.Deck = append(state.Deck, c)
			cards[i] = state.Deck[0]
			state.Deck = append([]*api.Card{}, state.Deck[1:]...)
		}
	}
	
}
