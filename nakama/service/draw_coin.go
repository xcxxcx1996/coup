package service

import (
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/xcxcx1996/coup/api"
	model "github.com/xcxcx1996/coup/state"
)

type DrawCoin struct {
	message runtime.MatchData
	coins   int32
}

func (a DrawCoin) Start(dispatcher runtime.MatchDispatcher, message runtime.MatchData, state *model.MatchState) (err error) {
	msg := &api.GetCoin{}

	// 推进行动
	if err = ValidAction(state, message, api.State_START, msg); err != nil {
		_ = dispatcher.BroadcastMessage(int64(api.OpCode_OPCODE_REJECTED), nil, nil, nil, true)
		return
	}
	a.coins = msg.Coins
	a.message = message

	info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> want to gain <span style={{ color: "red" }}>%v</span> coins.</p >`, a.message.GetUsername(), a.coins)
	SendNotification(info, dispatcher)
	if a.coins == 2 {
		state.Actions.Push(a)
		state.EnterDenyMoney()
		return
	}
	a.AfterDeny(dispatcher, state)
	state.NextTurn()
	return
}

// 正常拿钱没有人会质疑，只有公爵阻止,所以不用写
func (a DrawCoin) AfterQuestion(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	return
}

// 1. 未阻止
// 2. 阻止失败，弃牌，阻止失败，开始拿钱
func (a DrawCoin) AfterDeny(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	err = state.GainCoins(state.CurrentPlayerID, a.coins)
	if err != nil {
		return
	}
	info := fmt.Sprintf(`<p><span style={{ color: "red" }}>%v</span> success get <span style={{ color: "red" }}>%v</span> coins.</p >`, a.message.GetUsername(), a.coins)
	SendNotification(info, dispatcher)
	return
}

// 行动被阻止
// 被公爵阻止 deny_money=>after question

func (a DrawCoin) Stop(dispatcher runtime.MatchDispatcher, state *model.MatchState) (err error) {
	state.ActionComplete = true
	return
}

func (a DrawCoin) GetRole() api.Role {
	return api.Role_UNROLE
}

func (c DrawCoin) GetActor() string {
	return c.message.GetUserId()
}
