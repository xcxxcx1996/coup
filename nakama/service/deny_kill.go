package service

import (
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type DenyAssassian struct {
	Assassinated string
	Assassinor   string
}

// 阻止刺杀
func (d DenyAssassian) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.Deny{}

	// 验证
	if err = ValidAction(state, message, api.State_DENY_ASSASSIN, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}

	if !msg.IsDeny {
		ass, _ := state.Actions.Pop()
		ass.AfterDeny(dispatcher, state)
		return
	}

	// 阻止
	action, _ := state.Actions.Last()
	ass, ok := action.(Assassin)
	if !ok {
		return errors.New("wrong action")
	}

	d.Assassinor = ass.Assassinor
	d.Assassinated = ass.Assassinated

	state.Actions.Push(d)

	info := fmt.Sprintf("%v claim the queen, want to stop the kill", message.GetUsername())
	SendNotification(info, dispatcher)

	// question状态
	state.EnterQuestion()
	return nil
}

// 1.我有女王，你们都不质疑我,刺杀行为停止， 下一个回合
// 2. 有人质疑我，但是质疑失败，弃牌之后，我有女王，刺杀行为停止，下一个回合
func (d DenyAssassian) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	// 不质疑删除IAction， 然后assain改为 isdeny
	action, err := state.Actions.Pop()
	ass, ok := action.(Assassin)
	if !ok {
		return errors.New("wrong ass action")
	}
	ass.Stop(dispatcher, state)
	// ass.IsDeny = true
	info := fmt.Sprintln("question end, assassin was stopped")
	SendNotification(info, dispatcher)
	state.NextTurn()
	return
}

//
func (d DenyAssassian) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	return
}

// 阻止刺杀被质疑成功，刺杀进行
func (d DenyAssassian) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	action, err := state.Actions.Pop()
	ass, ok := action.(Assassin)
	if !ok {
		return errors.New("wrong action")
	}
	info := fmt.Sprintln("deny end, assassin excute")
	SendNotification(info, dispatcher)
	ass.AfterDeny(dispatcher, state)
	return
}

func (d DenyAssassian) GetRole() api.Role {
	return api.Role_QUEEN
}

func (c DenyAssassian) GetActor() string {
	return c.Assassinated
}
