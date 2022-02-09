package model

import (
	"errors"

	"github.com/xcxcx1996/coup/api"
)

func (s *MatchState) ChangeSingleCard(cardID, playerID string) (err error) {
	player, ok := s.PlayerInfos[playerID]
	if !ok {
		return errors.New("no that player")
	}
	for i, c := range player.Cards {
		if c.Id == cardID {
			s.Deck = append(s.Deck, c)
			player.Cards[i] = s.Deck[0]
			s.Deck = append([]*api.Card{}, s.Deck[1:]...)
		}
	}
	s.SufferDeck()
	return nil
}
