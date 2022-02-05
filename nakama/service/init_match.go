package service

import (
	"context"
	"math/rand"
	"time"

	// "github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
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
	logger.Info("初始化角色")
	//
	s.DeadlineRemainingTicks = 50

	t := time.Now().UTC()
	buf, err := serv.Marshaler.Marshal(&api.Start{
		PlayerInfos:     s.PlayerInfos,
		CurrentPlayerId: s.CurrentPlayerID,
		Message:         s.Message,
		Deadline:        t.Add(time.Duration(s.DeadlineRemainingTicks/tickRate) * time.Second).Unix(),
	})
	if err != nil {
	} else {
		logger.Error("error encoding message: %v", err)
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
		}
		state.PlayerInfos[userID] = playerinfo
		state.PlayerSequence = append(state.PlayerSequence, userID)
		state.Deck = append([]*api.Card{}, state.Deck[2:]...)
	}
	sufferPlayer(state.PlayerSequence)
	if len(state.PlayerSequence) > 0 {
		state.CurrentPlayerID = state.PlayerSequence[0]

	}
}

func initDeck(state *model.MatchState) {
	var deck []*api.Card
	for i := 0; i < 3; i++ {
		for i := 0; i < 5; i++ {
			card := &api.Card{
				Id:   getRandomString(),
				Role: int2Role(i),
			}
			deck = append(deck, card)
		}
	}
	sufferCard(deck)
	state.Deck = deck
}

func int2Role(i int) api.Role {
	switch i {
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

func sufferCard(deck []*api.Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}
func sufferPlayer(deck []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
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
