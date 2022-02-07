package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type DenySteal struct {
	Victim string
	Thief  string
}

func (d DenySteal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	msg := &api.Deny{}
	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 如果不阻止
	if !msg.IsDeny {
		return
	}
	// 阻止
	action, _ := state.Actions.Last()
	ass := action.(Steal)

	d.Victim = ass.Victim
	d.Thief = ass.Thief

	state.Actions.Push(d)
	info := fmt.Sprintln("我有大使或者女王，你不准偷钱")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	// question状态
	state.EnterQuestion()
}

func (d DenySteal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

	// 不质疑删除IAction， 然后assain改为 isdeny
	info := fmt.Sprintln("质疑结束，阻止成功")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	action, _ := state.Actions.Last()
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DENY_KILL), nil, nil, nil, true)
	action.Stop(dispatcher, state)
}
func (d DenySteal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

func (d DenySteal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("阻止失败")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	action, _ := state.Actions.Last()
	action.AfterDeny(dispatcher, state)
}

func (d DenySteal) GetRole() api.Role {
	return api.Role_QUEEN
}
