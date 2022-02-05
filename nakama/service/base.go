package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
	"google.golang.org/protobuf/encoding/protojson"
)

type MatchService struct {
	Marshaler   *protojson.MarshalOptions
	Unmarshaler *protojson.UnmarshalOptions

	// LastAction int
	// Time       int
}

func (serv *MatchService) Next(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	//对应的actionID
	switch state.Actions.Last().ActionID {
	case int64(api.OpCode_OPCODE_ASSASSIN):
		serv.CompleteAssassin(ctx, dispatcher, state, message)
	case int64(api.OpCode_OPCODE_CHANGE_CARD):
		serv.CompleteChangeCard(ctx, dispatcher, state, message)
	case int64(api.OpCode_OPCODE_DENY_STEAL):
		serv.ContinueAssassin(ctx, dispatcher, state, message)
	case int64(api.OpCode_OPCODE_DENY_KILL):
		serv.ContinueAssassin(ctx, dispatcher, state, message)
	case int64(api.OpCode_OPCODE_DRAW_COINS):
		serv.ContinueAssassin(ctx, dispatcher, state, message)
	case int64(api.OpCode_OPCODE_DENY_MONEY):
		serv.ContinueAssassin(ctx, dispatcher, state, message)
	default:
	}
	if state.IsActionComplete {
		serv.NextTurn(state)
	}

	//下一个回合
}

func New() (s *MatchService) {
	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}
	return &MatchService{
		Marshaler:   marshaler,
		Unmarshaler: unmarshaler,
	}
}

func (serv *MatchService) NextTurn(state *model.MatchState) {
	var nextPlayer string
	for i, p := range state.PlayerSequence {
		if state.CurrentPlayerID == p {
			nextPlayer = state.PlayerSequence[(i+1)%len(state.PlayerSequence)]
		}
	}
	for _, p := range state.PlayerInfos {
		if nextPlayer == p.Id {
			p.State = api.State_START
		} else {
			p.State = api.State_IDLE
		}
	}
}

func (serv *MatchService) GetNextPlayer(playerID string, state *model.MatchState) (nextPlayer string) {
	for i, p := range state.PlayerSequence {
		if playerID == p {
			nextPlayer = state.PlayerSequence[(i+1)%len(state.PlayerSequence)]
		}
	}
	return
}
