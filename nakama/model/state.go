package model

import (
	"math/rand"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
)

type MatchState struct {
	// ** 场面状态
	//剩余卡牌
	Deck           []*api.Card
	PlayerSequence []string
	Presences      map[string]runtime.Presence
	PlayerInfos    map[string]*api.PlayerInfo

	//玩家id对应的信息（卡牌，钱）

	//功能
	Random     *rand.Rand
	Label      *MatchLabel
	EmptyTicks int

	// Currently connected users, or reserved spaces.
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress int
	// // 质疑环节
	// IsQuestion bool
	// // 阻止环节
	// IsDeny    bool
	// IsDiscard bool
	State api.State
	// True if there's a game currently in progress.
	Playing bool
	// CurrentAction   Action
	Actions           Actions
	CurrentPlayerID   string
	Currentquestioner string
	CurrentDiscarder  string
	CurrentDenyer     string
	// Ticks until they must submit their move.
	NextGameRemainingTicks int64
	DeadlineRemainingTicks int64
	Message                string
	// The winner of the current game.
	// winner runtime.Presence
}

// ==========状态改变==========================
//
func (s *MatchState) EnterQuestion() {
	//下一个用户进入question状态
	s.State = api.State_QUESTION
	nextPlayer := s.GetNextPlayer(s.CurrentPlayerID)
	s.Currentquestioner = nextPlayer
	for _, p := range s.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_QUESTION
		} else {
			p.State = api.State_IDLE
		}
	}
}

func (s *MatchState) NextQuestionor() (end bool) {
	nextPlayer := s.GetNextPlayer(s.Currentquestioner)
	if nextPlayer == s.CurrentPlayerID {
		return true
	}
	for _, p := range s.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_QUESTION
		} else {
			p.State = api.State_IDLE
		}
	}
	return false
}

func (s *MatchState) NextTurn() {
	s.State = api.State_START
	var nextPlayer string
	for i, p := range s.PlayerSequence {
		if s.CurrentPlayerID == p {
			nextPlayer = s.PlayerSequence[(i+1)%len(s.PlayerSequence)]
		}
	}
	for _, p := range s.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_START
		} else {
			p.State = api.State_IDLE
		}
	}
}

func (s *MatchState) GetNextPlayer(playerID string) (nextPlayer string) {
	for i, p := range s.PlayerSequence {
		if playerID == p {
			nextPlayer = s.PlayerSequence[(i+1)%len(s.PlayerSequence)]
		}
	}
	return
}

// 指定玩家进入弃牌阶段
func (s *MatchState) EnterDicardState(playerID string) {
	s.State = api.State_DISCARD
	//下一个用户进入question状态
	for _, p := range s.PlayerInfos {
		if playerID == p.Id {
			s.CurrentDiscarder = playerID
			p.State = api.State_DISCARD
		} else {
			p.State = api.State_IDLE
		}
	}
}

func (s *MatchState) EnterChooseCard() {
	s.State = api.State_CHOOSE_CARD
	for _, p := range s.PlayerInfos {
		if s.CurrentPlayerID == p.Id {
			p.State = api.State_CHOOSE_CARD
		} else {
			p.State = api.State_IDLE
		}
	}
}

// 指定玩家可以阻止自我刺杀
func (s *MatchState) EnterDenyAssassin(playerID string) {
	s.State = api.State_DENY_ASSASSIN
	for _, p := range s.PlayerInfos {
		if playerID == p.Id {
			p.State = api.State_DENY_ASSASSIN
		} else {
			p.State = api.State_IDLE
		}
	}
}

// 指定玩家阻止偷钱
func (s *MatchState) EnterDenySteal(playerID string) {
	s.State = api.State_DENY_STEAL
	for _, p := range s.PlayerInfos {
		if playerID == p.Id {
			p.State = api.State_DENY_STEAL
		} else {
			p.State = api.State_IDLE
		}
	}
}

// 指定玩家阻止玩家偷钱，我有公爵，你不准拿钱！
func (s *MatchState) EnterDenyMoney() {
	s.State = api.State_DENY_MONEY
	nextPlayer := s.GetNextPlayer(s.CurrentPlayerID)
	s.CurrentDenyer = nextPlayer
	for _, p := range s.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_DENY_MONEY
		} else {
			p.State = api.State_IDLE
		}
	}
}
func (s *MatchState) NextDenyer() (end bool) {
	nextPlayer := s.GetNextPlayer(s.CurrentDenyer)
	if nextPlayer == s.CurrentPlayerID {
		return true
	}
	for _, p := range s.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_DENY_MONEY
		} else {
			p.State = api.State_IDLE
		}
	}
	return false
}

//next

// 质疑是否成功
func (s *MatchState) ValidQuestion() (cardID string, ok bool) {
	action, ok := s.Actions.Last()
	if !ok {
		return "", false
	}

	role := action.GetRole()
	cards := s.PlayerInfos[s.CurrentPlayerID].Cards
	return hasCard(role, cards)
}

// 某人有这张牌
func hasCard(role api.Role, cards []*api.Card) (cardID string, ok bool) {
	for _, c := range cards {
		if c.Role == role {
			return c.Id, false
		}
	}
	return "", true
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
