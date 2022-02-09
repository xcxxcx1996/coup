package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Coup(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	state.ActionComplete = true

	msg := &api.Coup{}
	if err = ValidAction(state, message, api.State_START, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}

	state.EnterDicardState(msg.PlayerId)
	info := fmt.Sprintf("%v launching a coup to %v", message.GetUsername(), state.GetPlayerNameByID(msg.PlayerId))
	SendNotification(info, dispatcher)
	return nil
}
