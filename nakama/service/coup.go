package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

func (serv *MatchService) Coup(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	msg := &api.Coup{}

	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// _ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_COUP), nil, nil, nil, true)
	state.EnterDicardState(msg.PlayerId)
	info := fmt.Sprintf("%v 对%v发动政变", message.GetUsername(), state.GetPlayerNameByID(msg.PlayerId))
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}
