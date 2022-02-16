package state

import (
	"log"

	"github.com/xcxcx1996/coup/api"
)

// ==========状态改变==========================
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
	action, err := s.Actions.Last()
	s.Currentquestioner = nextPlayer
	if err != nil {
		log.Println("error:", err)
	}
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

// 质疑是否成功 如果没有质疑成功，如果有质疑失败
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
	// 如果玩家没有卡牌，那么
	if has {
		s.ChangeSingleCard(cardId, action.GetActor())
	}
	return !has
}
