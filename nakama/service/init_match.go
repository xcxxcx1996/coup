package service

import (
	"context"
	"math/rand"

	// "github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) InitMatch(ctx context.Context, dispatcher runtime.MatchDispatcher, logger runtime.Logger, s *model.MatchState, tickRate int64) *model.MatchState {
	// Notify the players a new game has started.
	logger.Info("开始初始化")
	//
	initDeck(s)
	logger.Info("初始化卡牌池")
	//
	initPlayer(s)
	logger.Info("初始化角色: %v", s.PlayerInfos)
	//
	s.DeadlineRemainingTicks = 50
	s.State = api.State_START
	buf, err := global.Marshaler.Marshal(&api.Start{
		PlayerInfos:     s.PlayerInfos,
		CurrentPlayerId: s.CurrentPlayerID,
		Message:         s.Message,
		Deadline:        s.DeadlineRemainingTicks / tickRate,
	})
	if err != nil {
		logger.Error("error encoding message: %v", err)
	} else {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_START), buf, nil, nil, true)
	}
	return s
}

func initPlayer(state *model.MatchState) {
	state.PlayerInfos = make(map[string]*api.PlayerInfo, 4)
	for userID := range state.Presences {
		var playerinfo = &api.PlayerInfo{
			State: api.State_IDLE,
			Id:    userID,
			Coins: 3,
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

func initDeck(state *model.MatchState) {
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
