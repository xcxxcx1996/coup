package model

import (
	"math/rand"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
)

type MatchState struct {
	//剩余卡牌
	Deck []api.Role
	//玩家id对应的信息（卡牌，钱）
	Random     *rand.Rand
	Label      *MatchLabel
	EmptyTicks int
	// Currently connected users, or reserved spaces.
	Presences map[string]runtime.Presence
	// Number of users currently in the process of connecting to the match.
	JoinsInProgress int

	// True if there's a game currently in progress.
	Playing bool

	// Ticks until they must submit their move.
	DeadlineRemainingTicks int64
	// The winner of the current game.
	// winner runtime.Presence
}
