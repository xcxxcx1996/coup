package state

import (
	"math/rand"

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
	ActionComplete    bool
	CurrentPlayerID   string
	Currentquestioner string
	CurrentDiscarder  string
	CurrentDenyer     string
	// Ticks until they must submit their move.
	NextGameRemainingTicks int64
	DeadlineRemainingTicks int64
	MaxDeadlineTicks       int64
	Message                string
	// The winner of the current game.
	// winner runtime.Presence
}

func (s *MatchState) NextTurn() {
	s.ActionComplete = false
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
	s.CurrentPlayerID = nextPlayer
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

//next
