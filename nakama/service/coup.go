package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

func (serv *MatchService) Coup(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	state.ActionComplete = true

	msg := &api.Coup{}
	if err = ValidAction(state, message, api.State_START, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	err = state.LoseCoins(msg.PlayerId, 7)
	if err != nil {
		return
	}
	state.EnterDicardState(msg.PlayerId)
	info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> launch a coup to <span style={{ color: "red" }}>%v</span></p>`, message.GetUsername(), state.GetPlayerNameByID(msg.PlayerId))
	SendNotification(info, dispatcher)
	return nil

}
