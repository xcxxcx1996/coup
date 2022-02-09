package state

import "errors"

func (s *MatchState) Alive(playerID string) (bool, error) {
	for _, p := range s.PlayerInfos {
		if playerID == p.Id {
			return len(p.Cards) >= 1, nil
		}
	}
	return false, errors.New("找不到该角色")

}

func (s *MatchState) EliminatePlayer(playerID string) error {
	for i, p := range s.PlayerSequence {
		if playerID == p {
			s.PlayerSequence = append(s.PlayerSequence[:i], s.PlayerSequence[i+1:]...)
			return nil
		}
	}
	return errors.New("找不到该角色")
}

func (s *MatchState) GetNextPlayer(playerID string) (nextPlayer string) {
	for i, p := range s.PlayerSequence {
		if playerID == p {
			nextPlayer = s.PlayerSequence[(i+1)%len(s.PlayerSequence)]
		}
	}
	return
}
