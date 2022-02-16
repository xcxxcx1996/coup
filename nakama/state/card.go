package state

import (
	"errors"
	"math/rand"
	"time"

	"github.com/xcxcx1996/coup/api"
)

func (s *MatchState) ChangeSingleCardByRole(role api.Role, playerID string) (err error) {
	player, ok := s.PlayerInfos[playerID]
	if !ok {
		return errors.New("no that player")
	}
	for i, c := range player.Cards {
		if c.Role == role {
			s.Deck = append(s.Deck, c)
			player.Cards[i] = s.Deck[0]
			s.Deck = append([]*api.Card{}, s.Deck[1:]...)
			s.SufferDeck()
			return
		}
	}
	return errors.New("no card")
}

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

func (s *MatchState) SufferDeck() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s.Deck), func(i, j int) { s.Deck[i], s.Deck[j] = s.Deck[j], s.Deck[i] })
}

func (s *MatchState) SufferPlayer() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(s.PlayerSequence), func(i, j int) { s.PlayerSequence[i], s.PlayerSequence[j] = s.PlayerSequence[j], s.PlayerSequence[i] })
}

func (s *MatchState) GetPlayerNameByID(playerID string) string {
	for _, p := range s.PlayerInfos {
		if playerID == p.Id {
			return p.Name
		}
	}
	return ""
}

func (s *MatchState) DeleteCard(card_id, player_id string) (*api.Card, error) {
	cards := s.PlayerInfos[player_id].Cards
	for i, c := range cards {
		if c.Id == card_id {
			cards = append(cards[:i], cards[i+1:]...)
			s.PlayerInfos[player_id].Cards = cards
			return c, nil
		}
	}
	return nil, errors.New("no record")
}

// 某人有这张牌
func hasCard(role api.Role, cards []*api.Card) (cardID string, ok bool) {
	for _, c := range cards {
		if c.Role == role {
			return c.Id, true
		}
	}
	return "", false
}
