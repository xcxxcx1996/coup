package service

import (
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type DenySteal struct {
	Victim string
	Thief  string
}

func (d DenySteal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.Deny{}
	// 推进行动
	if err = ValidAction(state, message, api.State_DENY_STEAL, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	// 如果不阻止
	if !msg.IsDeny {
		ass, _ := state.Actions.Pop()
		ass.AfterDeny(dispatcher, state)
		return
	}
	// 阻止
	action, err := state.Actions.Last()
	if err != nil {
		return
	}
	ass, ok := action.(Steal)
	if !ok {
		return errors.New("wrong action")
	}

	d.Victim = ass.Victim
	d.Thief = ass.Thief

	state.Actions.Push(d)
	info := fmt.Sprintf("%v claim the Queen, want to stop the steal action", message.GetUsername())
	SendNotification(info, dispatcher)
	state.EnterQuestion()
	return nil
}

func (d DenySteal) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {

	// 不质疑删除IAction， 然后assain改为 isdeny
	info := fmt.Sprintln("question end, stop steal")
	SendNotification(info, dispatcher)
	_, err = state.Actions.Pop()
	if err != nil {
		return
	}
	steal, err := state.Actions.Last()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	err = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DENY_KILL), nil, nil, nil, true)
	if err != nil {
		return
	}
	err = steal.Stop(dispatcher, state)
	if err != nil {
		return
	}
	state.NextTurn()
	return
}

func (d DenySteal) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	return
}

func (d DenySteal) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {

	info := fmt.Sprintln("deny failed")
	SendNotification(info, dispatcher)

	_, err = state.Actions.Pop()
	if err != nil {
		return
	}
	action, err := state.Actions.Last()
	if err != nil {
		return
	}
	err = action.AfterDeny(dispatcher, state)
	return
}

func (d DenySteal) GetRole() api.Role {
	return api.Role_QUEEN
}

func (c DenySteal) GetActor() string {
	return c.Victim
}
