package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Coup(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.Coup{}
	state.ActionComplete = true

	valid := ValidAction(state, message, api.State_START, msg)
	// 推进行动
	if !valid {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_COUP), nil, nil, nil, true)
	state.EnterDicardState(msg.PlayerId)
	info := fmt.Sprintf("%v launching a coup to %v", message.GetUsername(), state.GetPlayerNameByID(msg.PlayerId))
	SendNotification(info, dispatcher)

	// buf, _ := global.Marshaler.Marshal(&api.ActionInfo{Message: info})
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}
