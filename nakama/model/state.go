package model

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
	JoinsInProgress  int
	IsQuestion       bool
	IsActionComplete bool
	// True if there's a game currently in progress.
	Playing bool
	// CurrentAction   Action
	Actions         Actions
	CurrentPlayerID string
	// Ticks until they must submit their move.
	NextGameRemainingTicks int64
	DeadlineRemainingTicks int64
	Message                string
	// The winner of the current game.
	// winner runtime.Presence
}
