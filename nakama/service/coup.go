package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Coup(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.Coup{}
	err := serv.Unmarshaler.Unmarshal(message.GetData(), msg)

	if err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}

	cards := state.PlayerInfos[msg.PlayerId].Cards
	for i, c := range cards {
		if c.Id == msg.CardId {
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)
}
