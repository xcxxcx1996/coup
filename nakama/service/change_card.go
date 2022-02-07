package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	"github.com/xcxcx1996/coup/model"
)

const CHANGE_NUM = 2

type ChangeCard struct {
	message runtime.MatchData
}

// 玩家提出换牌的主张，然后进入质疑
func (c ChangeCard) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) {
	// state.InitAction()

	myTurn := message.GetUserId() == state.CurrentPlayerID
	if !myTurn {
		// Client sent bad data.
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	c.message = message
	state.Actions.Push(c)

	info := fmt.Sprintf("%v宣称大使，想要换牌", message.GetUsername())
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)

	// question状态
	state.EnterQuestion()
}

// 给玩家发送换牌的资源
func (c ChangeCard) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	newcards := &api.ChangeCardRequest{
		Cards: state.Deck[:CHANGE_NUM],
	}
	buf, err := global.Marshaler.Marshal(newcards)
	if err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, []runtime.Presence{c.message}, nil, true)
	} else {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_CHOOSE_CARD), buf, []runtime.Presence{c.message}, nil, true)
	}
	info := fmt.Sprintln("质疑结束，玩家开始换牌")
	buf, _ = global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.EnterChooseCard()
	state.Actions.Pop()
}

// 下一个是抽牌，此处不执行
func (c ChangeCard) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) {

}

// 被质疑成功，停止
func (c ChangeCard) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) {
	info := fmt.Sprintln("换牌被阻止")
	buf, _ := global.Marshaler.Marshal(&api.Info{Info: info})
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_INFO), buf, nil, nil, true)
	state.Actions.Pop()
	state.NextTurn()
}

func (c ChangeCard) GetRole() api.Role {
	return api.Role_DIPLOMAT
}
