package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	"github.com/xcxcx1996/coup/global"
	model "github.com/xcxcx1996/coup/state"
)

const CHANGE_NUM = 2

type ChangeCard struct {
	message runtime.MatchData
}

// 玩家提出换牌的主张，然后进入质疑
func (c ChangeCard) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	// state.InitAction()

	err = ValidAction(state, message, api.State_START, nil)
	// 推进行动
	if err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	c.message = message
	state.Actions.Push(c)

	info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> claims the <span style={{ color: "red" }}>DIPLOMAT</span> and want to change the card.</p >`, message.GetUsername())
	SendNotification(info, dispatcher)

	// question状态
	state.EnterQuestion()
	return nil
}

// 给玩家发送换牌的资源
func (c ChangeCard) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) error {
	newcards := &api.ChangeCardRequest{
		Cards: state.Deck[:CHANGE_NUM],
	}
	buf, err := global.Marshaler.Marshal(newcards)
	if err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, []runtime.Presence{c.message}, nil, true)
		return err
	}
	_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_CHOOSE_CARD), buf, []runtime.Presence{c.message}, nil, true)
	info := fmt.Sprintln("<p>Questioning ends, choosing card action begins.</p>")
	SendNotification(info, dispatcher)
	state.EnterChooseCard()
	return err
}

// 下一个是抽牌，此处不执行
func (c ChangeCard) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	return
}

// 被质疑成功，停止
func (c ChangeCard) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	info := fmt.Sprintln("<p>Changing card was stoped.</p>")
	SendNotification(info, dispatcher)
	return nil
}

func (c ChangeCard) GetRole() api.Role {
	return api.Role_DIPLOMAT
}

func (c ChangeCard) GetActor() string {
	return c.message.GetUserId()
}
