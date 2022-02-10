package service

import (
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type DenySteal struct {
	Victim string
	Thief  string
	Role   api.Role
}

func (d DenySteal) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.DenySteal{}
	// 推进行动
	if err = ValidAction(state, message, api.State_DENY_STEAL, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	d.Role = msg.Role
	// 如果不阻止，偷窃进行，下一个回合
	if !msg.IsDeny {
		ass, _ := state.Actions.Pop()
		ass.AfterDeny(dispatcher, state)
		state.NextTurn()
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
	action, err := state.Actions.Pop()
	if err != nil {
		return
	}
	steal, ok := action.(Steal)
	if !ok {
		return errors.New("wrong steal")
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

	info := fmt.Sprintln("Deny failed.")
	SendNotification(info, dispatcher)

	action, err := state.Actions.Pop()
	if err != nil {
		return
	}

	steal, ok := action.(Steal)
	if !ok {
		return errors.New("wrong action")
	}
	err = steal.AfterDeny(dispatcher, state)
	return
}

func (d DenySteal) GetRole() api.Role {
	return d.Role
}

func (c DenySteal) GetActor() string {
	return c.Victim
}
