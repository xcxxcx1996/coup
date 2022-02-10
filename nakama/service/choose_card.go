package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

// 玩家选择完牌后，发送choose，然后服务器换牌
func (serv *MatchService) CompleteChangeCard(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	//更换牌
	state.ActionComplete = true
	msg := &api.ChangeCardResponse{}

	err = ValidAction(state, message, api.State_CHOOSE_CARD, msg)
	// 推进行动
	if err != nil {
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
	info := fmt.Sprintf("<p>%v complete the card change.</p>", message.GetUsername())
	SendNotification(info, dispatcher)
	state.PlayerInfos[state.CurrentPlayerID].Cards = reservedCards
	state.SufferDeck()
	//下一个回合
	state.NextTurn()
	return nil
}
