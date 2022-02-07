package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type DenyMoney struct {
}

// 公爵阻止别人拿钱
func (d DenyMoney) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {

	msg := &api.Deny{}
	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 如果不阻止，下一个刺杀者
	if !msg.IsDeny {

		end := state.NextDenyer()
		// 所有人都不阻止
		if end {
			action, _ := state.Actions.Last()
			ass := action.(DrawCoin)
			ass.AfterDeny(dispatcher, state)
			_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DENY_MONEY), nil, nil, nil, true)
		}
		return
	}
	// 阻止
	state.Actions.Push(d)

	info := fmt.Sprintln("男爵")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	// question状态
	state.EnterQuestion()
}

//所有人都不质疑,那么阻止别人拿两块钱
func (d DenyMoney) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

	// 不质疑删除IAction， 然后assain改为 isdeny
	info := fmt.Sprintln("阻止成功")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)

	state.Actions.Pop()
	action, _ := state.Actions.Last()
	action.Stop(dispatcher, state)
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DENY_MONEY), nil, nil, nil, true)
}

func (d DenyMoney) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

func (d DenyMoney) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	state.Actions.Pop()
	action, _ := state.Actions.Last()
	action.AfterDeny(dispatcher, state)
	info := fmt.Sprintln("阻止失败")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}

func (d DenyMoney) GetRole() api.Role {
	return api.Role_QUEEN
}
