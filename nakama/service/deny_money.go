package service

import (
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/model"
)

type DenyMoney struct {
	message runtime.MatchData
}

// 公爵阻止别人拿钱
func (d DenyMoney) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.Deny{}
	if err = ValidAction(state, message, api.State_DENY_MONEY, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	d.message = message
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

	info := fmt.Sprintf("%v cliam the barron, want to stop get money", message.GetUsername())
	SendNotification(info, dispatcher)

	// question状态
	state.EnterQuestion()
	return nil
}

//所有人都不质疑,那么阻止别人拿两块钱
func (d DenyMoney) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	info := fmt.Sprintln("question end, action was stop")
	SendNotification(info, dispatcher)
	_, err = state.Actions.Pop()
	if err != nil {
		return
	}
	action, _ := state.Actions.Last()
	gainCoins, ok := action.(Assassin)
	if !ok {
		return errors.New("wrong action")
	}
	err = gainCoins.Stop(dispatcher, state)
	if err != nil {
		return err
	}
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_DENY_MONEY), nil, nil, nil, true)

	// 下一个回合
	defer state.NextTurn()
	return nil
}

func (d DenyMoney) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	return nil
}

func (d DenyMoney) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	//
	_, err = state.Actions.Pop()
	if err != nil {
		return
	}
	action, err := state.Actions.Last()
	if err != nil {
		return
	}

	gainCoins, ok := action.(Assassin)
	if !ok {
		return errors.New("wrong action")
	}

	info := fmt.Sprintln("deny end, assassin excute")
	SendNotification(info, dispatcher)
	gainCoins.AfterDeny(dispatcher, state)
	return nil
}

func (d DenyMoney) GetRole() api.Role {
	return api.Role_QUEEN
}

func (d DenyMoney) GetActor() string {
	return d.message.GetUserId()
}
