package state

import "github.com/xcxcx1996/coup/api"

// 指定玩家可以阻止自我刺杀
func (s *MatchState) EnterDenyAssassin(playerID string) {
	s.State = api.State_DENY_ASSASSIN
	s.CurrentDenyer = playerID
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
	s.CurrentDenyer = playerID
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
	s.CurrentDenyer = nextPlayer
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
