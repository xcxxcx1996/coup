package state

import (
	"log"
	"math/rand"

	"github.com/xcxcx1996/coup/api"
)

func (s *MatchState) Init(deadlineticks int64) {
	//设置时间
	s.DeadlineRemainingTicks = deadlineticks
	s.MaxDeadlineTicks = deadlineticks
	s.initDeck()
	s.initPlayer()
	log.Printf("state:%v", s.PlayerInfos)
	log.Printf("state:%v", s.PlayerInfos)
	log.Printf("state:%v", s.PlayerInfos)
}

func (state *MatchState) initPlayer() {
	state.PlayerInfos = make(map[string]*api.PlayerInfo, 4)
	state.Players = state.Presences
	for userID := range state.Presences {
		var playerinfo = &api.PlayerInfo{
			State: api.State_IDLE,
			Id:    userID,
			Coins: 2,
			Cards: state.Deck[:2],
			Name:  state.Presences[userID].GetUsername(),
		}
		state.PlayerInfos[userID] = playerinfo
		state.PlayerSequence = append(state.PlayerSequence, userID)
		state.Deck = append([]*api.Card{}, state.Deck[2:]...)
	}
	state.SufferPlayer()
	if len(state.PlayerSequence) > 0 {
		state.CurrentPlayerID = state.PlayerSequence[0]
		state.PlayerInfos[state.CurrentPlayerID].State = api.State_START
	}
}

func (state *MatchState) initDeck() {
	var deck []*api.Card
	for i := 0; i < 3; i++ {
		for i := 1; i < 6; i++ {
			card := &api.Card{
				Id:   getRandomString(),
				Role: int2Role(i),
			}
			deck = append(deck, card)
		}
	}
	state.Deck = deck
	state.SufferDeck()
}

func int2Role(i int) api.Role {
	switch i - 1 {
	case 0:
		return api.Role_DIPLOMAT
	case 1:
		return api.Role_QUEEN
	case 2:
		return api.Role_CAPTAIN
	case 3:
		return api.Role_ASSASSIN
	case 4:
		return api.Role_BARON
	}
	return api.Role_DIPLOMAT
}

func getRandomString() string {
	n := 32
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
