package service

import (
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

// 被别人看出来了，然后把自己的牌放回并重新抽一张
func (serv *MatchService) ChangeSingleCard(cardID, playerID string, state *model.MatchState) {
	cards := state.PlayerInfos[playerID].Cards
	for i, c := range cards {
		if c.Id == cardID {
			state.Deck = append(state.Deck, c)
			cards[i] = state.Deck[0]
			state.Deck = append([]*api.Card{}, state.Deck[1:]...)
		}
	}
}
