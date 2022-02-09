package state

import "github.com/xcxcx1996/coup/api"

// ==========状态改变==========================
//
func (s *MatchState) EnterQuestion() {
	//下一个用户进入question状态
	s.State = api.State_QUESTION
	action, _ := s.Actions.Last()
	nextPlayer := s.GetNextPlayer(action.GetActor())
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
	action, _ := s.Actions.Last()

	if nextPlayer == action.GetActor() {
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

// 质疑是否成功
func (s *MatchState) ValidQuestion() (ok bool) {
	// getRole
	action, err := s.Actions.Last()
	if err != nil {
		return false
	}
	role := action.GetRole()
	// get player
	cards := s.PlayerInfos[action.GetActor()].Cards
	cardId, has := hasCard(role, cards)
	if !ok {
		return !has
	}
	s.ChangeSingleCard(cardId, action.GetActor())
	return has
}
