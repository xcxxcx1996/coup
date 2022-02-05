package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type AssassinService struct {
}

// 玩家开始准备刺杀，进入question状态
func (serv *MatchService) BeforeAssassin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	msg := &api.Assassin{}

	if err := serv.Unmarshaler.Unmarshal(message.GetData(), msg); err != nil {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
	}
	var newAction = model.Action{ActionID: int64(api.OpCode_OPCODE_ASSASSIN), ActionArg: msg}
	state.Actions.Push(&newAction)
	// question状态
	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
	serv.EnterQuestionState(nextPlayer, state)
}

// 继续刺杀过程，进入是否阻止刺杀步骤
func (serv *MatchService) ContinueAssassin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {

	nextPlayer := serv.GetNextPlayer(message.GetUserId(), state)
	serv.EnterDenyAssassin(nextPlayer, state)

}

// 开始刺杀动作，刺杀完成，下一个玩家
func (serv *MatchService) CompleteAssassin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	action, ok := state.Actions.Pop()
	if !ok {
		//丢失
		return
	}
	arg := action.ActionArg.(*api.Assassin)
	cards := state.PlayerInfos[arg.PlayerId].Cards
	for i, c := range cards {
		if c.Id == arg.CardId {
			cards = append(cards[:i], cards[i+1:]...)
		}
	}
	// todo  next turn
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)
}
