package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

type DenyAssassian struct {
	Assassinated string
	Assassinor   string
	IsDeny       bool
}

func (d DenyAssassian) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	msg := &api.Deny{}
	myTurn := message.GetUserId() == state.CurrentPlayerID

	err := global.Unmarshaler.Unmarshal(message.GetData(), msg)
	if err != nil || !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 如果不阻止，没有下一个人，继续刺杀
	if !msg.IsDeny {
		ass, _ := state.Actions.Pop()
		ass.AfterDeny(dispatcher, state)
	}
	// 阻止
	action, _ := state.Actions.Last()
	ass := action.(Assassin)

	d.Assassinor = ass.Assassinated
	d.Assassinated = ass.Assassinor

	state.Actions.Push(d)

	info := fmt.Sprintln("我有女王")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	// question状态
	state.EnterQuestion()
}

// 1.我有女王，你们都不质疑我,刺杀行为停止 2. 有人质疑我，我有女王，刺杀行为停止
func (d DenyAssassian) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	// 不质疑删除IAction， 然后assain改为 isdeny
	state.Actions.Pop()
	action, _ := state.Actions.Last()
	action.(Assassin).Stop(dispatcher, state)
	// ass.IsDeny = true
	info := fmt.Sprintln("质疑结束，阻止刺杀")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
}

//
func (d DenyAssassian) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
}

// 被质疑成功
func (d DenyAssassian) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("阻止刺杀行动被停止，继续刺杀")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	action, _ := state.Actions.Last()
	action.AfterDeny(dispatcher, state)
}

func (d DenyAssassian) GetRole() api.Role {
	return api.Role_QUEEN
}
