package service

import (
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

// 玩家选择完牌后，发送choose，然后服务器换牌
func (serv *MatchService) CompleteChangeCard(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	//更换牌
	myTurn := message.GetUserId() == state.CurrentPlayerID

	msg := &api.ChangeCardResponse{}

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}

	// 手牌替换成目标id的两张牌
	handcards := state.PlayerInfos[state.CurrentPlayerID].Cards
	state.Deck = append(state.Deck, handcards...)
	var reservedCards []*api.Card
	for _, m := range msg.Cards {
		for _, c := range state.Deck {
			if c.Id == m {
				reservedCards = append(reservedCards, c)
			}
		}
	}
	state.PlayerInfos[state.CurrentPlayerID].Cards = reservedCards
	state.SufferDeck()
	//下一个回合
	state.NextTurn()
}
