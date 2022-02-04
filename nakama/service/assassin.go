package service

import (
	"context"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) BeforeAssassin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	//todo

	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)
}
func (serv *MatchService) AfterAssassin(ctx context.Context, dispatcher runtime.MatchDispatcher, state *model.MatchState, message runtime.MatchData) {
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_ASSASSIN), nil, []runtime.Presence{message}, nil, true)

}
