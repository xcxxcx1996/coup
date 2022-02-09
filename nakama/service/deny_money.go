package service

import (
	"errors"
	"fmt"
	"log"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type DenyMoney struct {
	message runtime.MatchData
	// Role    api.Role
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
		// 下一个人
		end := state.NextDenyer()
		// 所有人都不阻止
		if end {
			action, _ := state.Actions.Pop()
			ass, ok := action.(DrawCoin)
			if !ok {
				return errors.New("wrong draw coins action")
			}
			ass.AfterDeny(dispatcher, state)
		}
		return
	}
	// 阻止
	log.Printf("Actions:%v,deny_action:%v", state.Actions, d)
	state.Actions.Push(d)
	log.Printf("Actions:%v", state.Actions)
	info := fmt.Sprintf("%v cliam the barron adn want to stop get money.", message.GetUsername())
	SendNotification(info, dispatcher)

	// question状态
	state.EnterQuestion()
	return nil
}

//所有人都不质疑,那么阻止别人拿两块钱
func (d DenyMoney) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	info := fmt.Sprintln("Question end and coins stop.")
	SendNotification(info, dispatcher)
	action, err := state.Actions.Pop()
	if err != nil {
		return
	}

	gainCoins, ok := action.(DrawCoin)
	if !ok {
		log.Println("error: wrong action")
		return errors.New("wrong action")
	}
	err = gainCoins.Stop(dispatcher, state)
	if err != nil {
		return err
	}
	// 下一个回合
	state.NextTurn()
	return nil
}

func (d DenyMoney) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	return nil
}

func (d DenyMoney) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	//
	action, err := state.Actions.Pop()
	if err != nil {
		return
	}
	gainCoins, ok := action.(DrawCoin)
	if !ok {
		return errors.New("wrong action")
	}
	info := fmt.Sprintln("deny end, get money excute")
	SendNotification(info, dispatcher)
	gainCoins.AfterDeny(dispatcher, state)
	return nil
}

func (d DenyMoney) GetRole() api.Role {
	return api.Role_BARON
}

func (d DenyMoney) GetActor() string {
	return d.message.GetUserId()
}
